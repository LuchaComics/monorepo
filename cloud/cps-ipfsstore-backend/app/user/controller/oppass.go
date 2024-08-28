package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	user_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/utils/httperror"
)

type UserOperationChangePasswordRequest struct {
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	Password         string             `bson:"password" json:"password"`
	PasswordRepeated string             `bson:"password_repeated" json:"password_repeated"`
}

func (impl *UserControllerImpl) validateOperationChangePasswordRequest(ctx context.Context, dirtyData *UserOperationChangePasswordRequest) error {
	e := make(map[string]string)

	if dirtyData.UserID.IsZero() {
		e["user_id"] = "missing value"
	}
	if dirtyData.Password == "" {
		e["password"] = "missing value"
	}
	if len(dirtyData.Password) > 255 {
		e["password"] = "too long"
	}
	if dirtyData.PasswordRepeated == "" {
		e["password_repeated"] = "missing value"
	}
	if len(dirtyData.PasswordRepeated) > 255 {
		e["password_repeated"] = "too long"
	}
	if dirtyData.Password != dirtyData.PasswordRepeated {
		e["password"] = "value does not match"
		e["password_repeated"] = "value does not match"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *UserControllerImpl) ChangePassword(ctx context.Context, req *UserOperationChangePasswordRequest) error {
	//
	// Get variables from our user authenticated session.
	//

	tid, _ := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
	role, _ := ctx.Value(constants.SessionUserRole).(int8)
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userName, _ := ctx.Value(constants.SessionUserName).(string)
	// ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)

	switch role {
	case user_s.UserRoleRoot:
		break
	default:
		impl.Logger.Error("you do not have permission to change password")
		return httperror.NewForForbiddenWithSingleField("message", "you do not have permission to change password")
	}

	//
	// Perform our validation and return validation error on any issues detected.
	//

	if err := impl.validateOperationChangePasswordRequest(ctx, req); err != nil {
		impl.Logger.Error("validation error", slog.Any("error", err))
		return err
	}

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		//
		// Fetch the original user.
		//

		u, err := impl.UserStorer.GetByID(sessCtx, req.UserID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if u == nil {
			return nil, httperror.NewForBadRequestWithSingleField("user_id", "user does not exist")
		}

		// Defensive Code: Tenancy protection
		if u.TenantID != tid && role != user_s.UserRoleRoot {
			return nil, httperror.NewForForbiddenWithSingleField("security", "you do not belong to this organization")
		}

		passwordHash, err := impl.Password.GenerateHashFromPassword(req.Password)
		if err != nil {
			impl.Logger.Error("hashing error", slog.Any("error", err))
			return nil, err
		}

		u.PasswordHash = passwordHash
		u.PasswordHashAlgorithm = impl.Password.AlgorithmName()
		u.ModifiedAt = time.Now()
		u.ModifiedByUserID = userID
		// u.ModifiedByUserName = userName
		// u.ModifiedFromIPAddress = ipAddress

		if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
			impl.Logger.Error("update error", slog.Any("err", err))
			return nil, err
		}

		return u, nil
	}

	// Start a transaction
	if _, err := session.WithTransaction(ctx, transactionFunc); err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return err
	}

	return nil
}

// func (impl *UserControllerImpl) createUserFromUser(sessCtx mongo.SessionContext, a *a_s.User) (*user_s.User, error) {
// 	//
// 	// Get variables from our user authenticated session.
// 	//
//
// 	tid, _ := sessCtx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
// 	// role, _ := sessCtx.Value(constants.SessionUserRole).(int8)
// 	userID, _ := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
// 	userName, _ := sessCtx.Value(constants.SessionUserName).(string)
// 	ipAddress, _ := sessCtx.Value(constants.SessionIPAddress).(string)
//
// 	impl.Logger.Warn("user account does not exist for user update, creating user account now",
// 		slog.Any("UserID", a.ID))
//
// 	//
// 	// Create user.
// 	//
//
// 	u := &user_s.User{
// 		ID:                      primitive.NewObjectID(),
// 		TenantID:                tid,
// 		FirstName:               a.FirstName,
// 		LastName:                a.LastName,
// 		Name:                    fmt.Sprintf("%s %s", a.FirstName, a.LastName),
// 		LexicalName:             fmt.Sprintf("%s, %s", a.LastName, a.FirstName),
// 		OrganizationName:        a.OrganizationName,
// 		OrganizationType:        a.OrganizationType,
// 		Email:                   a.Email,
// 		PasswordHashAlgorithm:   "DO BELOW...",
// 		PasswordHash:            "DO BELOW...",
// 		Role:                    user_s.UserRoleUser,
// 		WasEmailVerified:        true,
// 		EmailVerificationCode:   "",
// 		EmailVerificationExpiry: time.Now(),
// 		Phone:                   a.Phone,
// 		Country:                 a.Country,
// 		Region:                  a.Region,
// 		City:                    a.City,
// 		AgreeTOS:                true,
// 		AgreePromotionsEmail:    true,
// 		CreatedAt:               time.Now(),
// 		CreatedByUserID:         userID,
// 		CreatedByUserName:       userName,
// 		CreatedFromIPAddress:    ipAddress,
// 		ModifiedAt:              time.Now(),
// 		ModifiedByUserID:        userID,
// 		ModifiedByUserName:      userName,
// 		ModifiedFromIPAddress:   ipAddress,
// 		Status:                  user_s.UserStatusActive,
// 		Comments:                make([]*user_s.UserComment, 0),
// 		Salt:                    "",
// 		JoinedTime:              a.JoinDate,
// 		PrAccessCode:            "",
// 		PrExpiryTime:            time.Now(),
// 		PublicID:                0,
// 		Timezone:                "American/Toronto",
// 	}
//
// 	//
// 	// Temporary password.
// 	//
//
// 	// Generate a temporary password.
// 	temporaryPassword := primitive.NewObjectID().Hex()
//
// 	// Hash our password with the temporary password and attach to account.
// 	temporaryPasswordHash, err := impl.Password.GenerateHashFromPassword(temporaryPassword)
// 	if err != nil {
// 		impl.Logger.Error("hashing error", slog.Any("error", err))
// 		return nil, err
// 	}
// 	u.PasswordHashAlgorithm = impl.Password.AlgorithmName()
// 	u.PasswordHash = temporaryPasswordHash
//
// 	//
// 	// Insert the user into the database.
// 	//
//
// 	if err := impl.UserStorer.Create(sessCtx, u); err != nil {
// 		impl.Logger.Error("database create error", slog.Any("error", err))
// 		return nil, err
// 	}
// 	impl.Logger.Warn("user account created for user during update", slog.Any("UserID", a.ID))
//
// 	return u, nil
// }
