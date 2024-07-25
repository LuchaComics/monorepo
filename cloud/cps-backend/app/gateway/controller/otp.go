package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/dchest/uniuri"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	gateway_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/gateway/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

type OTPGenerateResponseIDO struct {
	Base32      string `json:"base32"`
	OTPAuthURL  string `json:"otpauth_url"`
	AccountName string `json:"account_name"`
}

// GenerateOTP function generates the time-based one-time password (TOTP) secret for the user. The user must use these values to generate a QR to present to the user.
func (impl *GatewayControllerImpl) GenerateOTP(ctx context.Context) (*OTPGenerateResponseIDO, error) {
	// Extract from our session the following data.
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	sessionID, _ := ctx.Value(constants.SessionID).(string)

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByID(sessCtx, userID)
		if err != nil {
			impl.Logger.Error("failed getting session user", slog.Any("err", err))
			return nil, err
		}
		if u == nil {
			impl.Logger.Warn("user does not exist validation error")
			return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
		}

		res := &OTPGenerateResponseIDO{}

		// Only generate a new OTP if no previous secret was generated. If
		// previously generated then reuse existing opt secret and otp auth url.
		if u.OTPSecret == "" && u.OTPAuthURL == "" {
			// STEP 1: Generate the OPT.
			key, err := totp.Generate(totp.GenerateOpts{
				Issuer:      impl.TemplatedEmailer.GetFrontendDomainName(),
				AccountName: u.Email,
				SecretSize:  15,
			})
			if err != nil {
				impl.Logger.Error("failed generating otp", slog.Any("err", err))
				return nil, err
			}

			// STEP 2: Save the secret to the user's profile.

			// DEVELOPERS NOTE: This is not a mistake! We do not enable OTP in
			// this function, we will make `OTPEnabled=true` in `ValidateOTP`
			// function when the user successfully verifies their 2FA code.
			u.OTPEnabled = false

			// DEVELOPERS NOTE: The most important variables for this function
			// is generating the OTP secret.
			u.OTPVerified = false
			u.OTPValidated = false
			u.OTPSecret = key.Secret()
			u.OTPAuthURL = key.URL()
			u.ModifiedAt = time.Now()

			if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
				impl.Logger.Error("failed updating session user with opt secret", slog.Any("err", err))
				return nil, err
			}

			// STEP 3: Update the authenticated user session.
			uBin, err := json.Marshal(u)
			if err != nil {
				impl.Logger.Error("marshalling error", slog.Any("err", err))
				return nil, err
			}
			atExpiry := 14 * 24 * time.Hour
			err = impl.Cache.SetWithExpiry(sessCtx, sessionID, uBin, atExpiry)
			if err != nil {
				impl.Logger.Error("cache set with expiry error", slog.Any("err", err))
				return nil, err
			}

			// STEP 4: Share the secret with the user so they may give to their
			// third-party authenticator app.
			res.Base32 = key.Secret()
			res.OTPAuthURL = key.URL()
			res.AccountName = fmt.Sprintf("%s: %s", impl.TemplatedEmailer.GetFrontendDomainName(), u.Email)

			impl.Logger.Debug("successfully generated opt secret and auth url", slog.Any("base_32", res.Base32), slog.Any("opt_auth_url", res.OTPAuthURL), slog.Any("account_name", res.AccountName))
		} else {
			// Reuse the existing opt secret and auth url.
			res.Base32 = u.OTPSecret
			res.OTPAuthURL = u.OTPAuthURL
			res.AccountName = fmt.Sprintf("%s: %s", impl.TemplatedEmailer.GetFrontendDomainName(), u.Email)
			impl.Logger.Warn("reusing previously generated opt secret and auth url", slog.Any("base_32", res.Base32), slog.Any("opt_auth_url", res.OTPAuthURL), slog.Any("account_name", res.AccountName))
		}

		impl.Logger.Debug("generated otp", slog.Any("res", res))

		return res, nil
	}

	// Start a transaction
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return result.(*OTPGenerateResponseIDO), nil
}

func (impl *GatewayControllerImpl) GenerateOTPAndQRCodePNGImage(ctx context.Context) ([]byte, error) {
	otpResponse, err := impl.GenerateOTP(ctx)
	if err != nil {
		impl.Logger.Error("failed generating otp",
			slog.Any("error", err))
		return nil, err
	}

	// Generate the QR code for the specific URL and return the `png` binary
	// file bytes.
	var png []byte
	png, err = qrcode.Encode(otpResponse.OTPAuthURL, qrcode.Medium, 256)
	if err != nil {
		impl.Logger.Error("encode error", slog.Any("error", err))
		return nil, err
	}

	impl.Logger.Debug("qr code ready",
		slog.Any("payload", otpResponse.OTPAuthURL))

	return png, err
}

type VerificationTokenRequestIDO struct {
	VerificationToken string `json:"verification_token"`
}

type VerificationTokenResponseIDO struct {
	User          *u_d.User `json:"user"`
	OTPBackupCode string    `json:"otp_backup_code"`
}

// VerifyOTP function verifies provided token from the third-party authenticator app. The purpose of this function is to finish the otp setup.
func (impl *GatewayControllerImpl) VerifyOTP(ctx context.Context, req *VerificationTokenRequestIDO) (*VerificationTokenResponseIDO, error) {
	// Extract from our session the following data.
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	sessionID, _ := ctx.Value(constants.SessionID).(string)

	// Variable used to store 2FA backup code when the OTP was successfully
	// verified in this function.
	var otpBackupCode string

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByID(sessCtx, userID)
		if err != nil {
			impl.Logger.Error("failed getting session user", slog.Any("err", err))
			return nil, err
		}
		if u == nil {
			impl.Logger.Warn("user does not exist validation error")
			return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
		}
		if u.OTPSecret == "" {
			impl.Logger.Warn("user did not run generate otp")
			return nil, httperror.NewForBadRequestWithSingleField("message", "you did not setup two-factor authentication")
		}

		//
		// STEP 1: Validate the inputted totp code.
		//

		if valid := totp.Validate(req.VerificationToken, u.OTPSecret); valid == false {

			//
			// STEP 2: Invalid tokens for whatever reason must return with error.
			//

			impl.Logger.Warn("totp verification failed or expired",
				slog.String("verification_token", req.VerificationToken),
				slog.String("otp_secret", u.OTPSecret))
			return nil, httperror.NewForBadRequestWithSingleField("verification_token", "token expired or invalid")
		}

		//
		// STEP 3: Update the user's profile.
		//

		// Enable 2FA once verified for all future logins.
		u.OTPEnabled = true

		// Set this to be `true` to indicate that the user successfully setup
		// our 2FA system because backend received a validated token.
		u.OTPVerified = true

		// Set this `true` vecause we successfully validated the token to
		// indicate the 2FA was successful.
		u.OTPValidated = true

		// Keep track of when user's account changes.
		u.ModifiedAt = time.Now()

		//
		// STEP 4: Generate 2FA backup code then hash it.
		//

		// Generate a random string of 20 characters (UUIDLen) length and store
		// it in our variable `otpBackupCode` which we will return later.
		otpBackupCode = fmt.Sprintf("%v%v%v",
			uniuri.NewLen(uniuri.UUIDLen),
			uniuri.NewLen(uniuri.UUIDLen),
			uniuri.NewLen(uniuri.UUIDLen))

		// Hash our backup code so it's saved in our database protected.
		otpBackupCodeHash, err := impl.Password.GenerateHashFromPassword(otpBackupCode)
		if err != nil {
			impl.Logger.Error("hashing error", slog.Any("error", err))
			return nil, err
		}

		u.OTPBackupCodeHash = otpBackupCodeHash
		u.OTPBackupCodeHashAlgorithm = impl.Password.AlgorithmName()

		//
		// STEP 5: Save to database.
		//

		if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
			impl.Logger.Error("failed updating user", slog.Any("err", err))
			return nil, err
		}

		//
		// STEP 6: Update the authenticated user session.
		//

		uBin, err := json.Marshal(u)
		if err != nil {
			impl.Logger.Error("marshalling error", slog.Any("err", err))
			return nil, err
		}
		atExpiry := 14 * 24 * time.Hour
		err = impl.Cache.SetWithExpiry(sessCtx, sessionID, uBin, atExpiry)
		if err != nil {
			impl.Logger.Error("cache set with expiry error", slog.Any("err", err))
			return nil, err
		}

		return u, nil
	}

	// Start a transaction
	u, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	res := &VerificationTokenResponseIDO{
		User:          u.(*u_d.User),
		OTPBackupCode: otpBackupCode, // Return the non-hashed value.
	}

	return res, nil
}

type ValidateTokenRequestIDO struct {
	Token string `json:"token"`
}

type ValidateTokenResponseIDO struct {
	User *u_d.User `json:"user"`
}

// ValidateOTP function verifies provided token from the third-party authenticator app. The purpose of this function is enable the loggin for 2FA.
func (impl *GatewayControllerImpl) ValidateOTP(ctx context.Context, req *ValidateTokenRequestIDO) (*ValidateTokenResponseIDO, error) {
	// Extract from our session the following data.
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	sessionID, _ := ctx.Value(constants.SessionID).(string)

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByID(sessCtx, userID)
		if err != nil {
			impl.Logger.Error("failed getting session user", slog.Any("err", err))
			return nil, err
		}
		if u == nil {
			impl.Logger.Warn("user does not exist validation error")
			return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
		}
		if u.OTPSecret == "" {
			impl.Logger.Warn("user did not run generate otp")
			return nil, httperror.NewForBadRequestWithSingleField("message", "you did not setup two-factor authentication")
		}

		//
		// STEP 1: Validate the inputted totp code.
		//

		if valid := totp.Validate(req.Token, u.OTPSecret); valid == false {

			//
			// STEP 2: Invalid tokens for whatever reason must return with error.
			//

			impl.Logger.Warn("totp verification failed or expired",
				slog.String("token", req.Token),
				slog.String("otp_secret", u.OTPSecret))
			return nil, httperror.NewForBadRequestWithSingleField("token", "expired or invalid")
		}

		//
		// STEP 3: Update the user's profile.
		//

		// Set this `true` vecause we successfully validated the token to
		// indicate the 2FA was successful.
		u.OTPValidated = true

		// Keep track of when user's account changes.
		u.ModifiedAt = time.Now()
		if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
			impl.Logger.Error("failed updating user", slog.Any("err", err))
			return nil, err
		}

		//
		// STEP 4: Update the authenticated user session.
		//

		uBin, err := json.Marshal(u)
		if err != nil {
			impl.Logger.Error("marshalling error", slog.Any("err", err))
			return nil, err
		}
		atExpiry := 14 * 24 * time.Hour
		err = impl.Cache.SetWithExpiry(sessCtx, sessionID, uBin, atExpiry)
		if err != nil {
			impl.Logger.Error("cache set with expiry error", slog.Any("err", err))
			return nil, err
		}

		return u, nil
	}

	// Start a transaction
	u, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	res := &ValidateTokenResponseIDO{
		User: u.(*u_d.User),
	}

	return res, nil
}

// DisableOTP function disables 2FA.
func (impl *GatewayControllerImpl) DisableOTP(ctx context.Context) (*u_d.User, error) {
	// Extract from our session the following data.
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	sessionID, _ := ctx.Value(constants.SessionID).(string)

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByID(sessCtx, userID)
		if err != nil {
			impl.Logger.Error("failed getting session user", slog.Any("err", err))
			return nil, err
		}
		if u == nil {
			impl.Logger.Warn("user does not exist validation error")
			return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
		}

		//
		// STEP 3: Update the user's profile.
		//

		u.OTPEnabled = false
		u.OTPVerified = false
		u.OTPValidated = false
		u.OTPSecret = ""
		u.OTPAuthURL = ""
		u.ModifiedAt = time.Now()
		if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
			impl.Logger.Error("failed updating user", slog.Any("err", err))
			return nil, err
		}

		//
		// STEP 4: Update the authenticated user session.
		//

		uBin, err := json.Marshal(u)
		if err != nil {
			impl.Logger.Error("marshalling error", slog.Any("err", err))
			return nil, err
		}
		atExpiry := 14 * 24 * time.Hour
		err = impl.Cache.SetWithExpiry(sessCtx, sessionID, uBin, atExpiry)
		if err != nil {
			impl.Logger.Error("cache set with expiry error", slog.Any("err", err))
			return nil, err
		}

		return u, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*u_d.User), nil
}

type RecoveryRequestIDO struct {
	BackupCode string `json:"backup_code"`
}

// RecoveryOTP function verifies provided `backup code` for the logged in user and disable 2FA plus log user in if successfully verified.
func (impl *GatewayControllerImpl) RecoveryOTP(ctx context.Context, req *RecoveryRequestIDO) (*gateway_s.LoginResponseIDO, error) {
	// Extract from our session the following data.
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Validate the inputted totp code.
	if req.BackupCode == "" {
		impl.Logger.Warn("user did not include backup code",
			slog.String("user_id", userID.Hex()))
		return nil, httperror.NewForBadRequestWithSingleField("backup_code", "missing value")
	}

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByID(sessCtx, userID)
		if err != nil {
			impl.Logger.Error("failed getting session user",
				slog.String("user_id", userID.Hex()),
				slog.Any("err", err))
			return nil, err
		}
		if u == nil {
			impl.Logger.Warn("user does not exist validation error",
				slog.String("user_id", userID.Hex()))
			return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
		}
		if u.OTPBackupCodeHash == "" || u.OTPBackupCodeHashAlgorithm == "" {
			impl.Logger.Error("user did not run generate backup code",
				slog.String("user_id", userID.Hex()),
				slog.String("user_email", u.Email),
				slog.String("hash", u.OTPBackupCodeHash),
				slog.String("algorithm", u.OTPBackupCodeHashAlgorithm),
			)
			return nil, httperror.NewForBadRequestWithSingleField("message", "you did not setup a 2fa backup code")
		}

		//
		// Validate backup code.
		//

		// Make sure we are using the correct hashing algorithm.
		if impl.Password.AlgorithmName() != u.OTPBackupCodeHashAlgorithm {
			impl.Logger.Warn("backup code encoded in different hashing algorithm",
				slog.String("given_algorithm", u.OTPBackupCodeHashAlgorithm),
				slog.String("system_algorithm", impl.Password.AlgorithmName()),
				slog.String("user_email", u.Email),
				slog.String("user_id", u.ID.Hex()),
			)
			return nil, httperror.NewForBadRequestWithSingleField("message", "please contact administrator")
		}

		// Verify the inputted password and hashed password match.
		otpBackupCodedMatch, _ := impl.Password.ComparePasswordAndHash(req.BackupCode, u.OTPBackupCodeHash)
		if otpBackupCodedMatch == false {
			impl.Logger.Warn("backup code check validation error",
				slog.String("user_id", userID.Hex()),
				slog.String("user_email", u.Email),
			)
			return nil, httperror.NewForBadRequestWithSingleField("backup_code", "does do not match with record")
		}

		//
		// Update the user's profile.
		//

		// Keep track of when user's account changes.
		u.OTPEnabled = false
		u.OTPVerified = false
		u.OTPValidated = false
		u.OTPSecret = ""
		u.OTPAuthURL = ""
		u.ModifiedAt = time.Now()
		if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
			impl.Logger.Error("failed updating user",
				slog.String("user_id", userID.Hex()),
				slog.Any("err", err))
			return nil, err
		}

		impl.Logger.Debug("successfully reset 2fa via backup code recovery",
			slog.String("user_id", userID.Hex()))

		return u, nil
	}

	// Start a transaction
	u, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.String("user_id", userID.Hex()),
			slog.Any("error", err))
		return nil, err
	}

	// Login the user.
	return impl.loginWithUser(ctx, u.(*u_d.User))
}
