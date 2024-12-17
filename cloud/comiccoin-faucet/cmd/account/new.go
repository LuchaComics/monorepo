package account

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"time"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/logger"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/jwt"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/password"
	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/securestring"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/database/mongodb"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/database/mongodbcache"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/repo"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/service"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

var (
	flagFirstName           string
	flagLastName            string
	flagEmail               string
	flagPassword            string
	flagPasswordConfirm     string
	flagPhone               string
	flagCountry             string
	flagCountryOther        string
	flagTimezone            string
	flagAgreeTermsOfService bool
	flagAgreePromotions     bool
	flagRole                int8
)

func NewAccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "new",
		Short: "Creates a new user account",
		Run: func(cmd *cobra.Command, args []string) {
			doRunNewAccount()
		},
	}

	cmd.Flags().StringVar(&flagFirstName, "first-name", "", "The first name")
	cmd.MarkFlagRequired("first-name")
	cmd.Flags().StringVar(&flagLastName, "last-name", "", "The last name")
	cmd.MarkFlagRequired("last-name")
	cmd.Flags().StringVar(&flagEmail, "email", "", "The email")
	cmd.MarkFlagRequired("email")
	cmd.Flags().StringVar(&flagPassword, "password", "", "The password for user")
	cmd.MarkFlagRequired("password")
	cmd.Flags().StringVar(&flagPasswordConfirm, "password-confirm", "", "The password confirm to verify user password is correct")
	cmd.MarkFlagRequired("password-confirm")
	cmd.Flags().StringVar(&flagPhone, "phone", "", "The phone of the user")
	cmd.MarkFlagRequired("phone")
	cmd.Flags().StringVar(&flagCountry, "country", "", "The country the user belongs to")
	cmd.MarkFlagRequired("country")
	cmd.Flags().StringVar(&flagCountryOther, "country-other", "", "The other country the user belongs to")
	// cmd.MarkFlagRequired("country-other")
	cmd.Flags().StringVar(&flagTimezone, "timezone", "", "The timezone the user belongs to")
	cmd.MarkFlagRequired("timezone")
	cmd.Flags().BoolVar(&flagAgreeTermsOfService, "agree-tos", false, "The agree to terms of service")
	cmd.MarkFlagRequired("agree-tos")
	cmd.Flags().BoolVar(&flagAgreePromotions, "agree-promotions", false, "The agree to promotions")
	cmd.MarkFlagRequired("agree-promotions")
	cmd.Flags().Int8Var(&flagRole, "role", 0, "The user role")
	cmd.MarkFlagRequired("role")
	return cmd
}

func doRunNewAccount() {
	//
	// Load up our dependencies and configuration
	//

	// Common
	logger := logger.NewProvider()
	// kmutex := kmutexutil.NewKMutexProvider()
	cfg := config.NewProviderUsingEnvironmentVariables()
	dbClient := mongodb.NewProvider(cfg, logger)
	// keystore := keystore.NewAdapter()
	passp := password.NewProvider()
	// blackp := blacklist.NewProvider()
	jwtp := jwt.NewProvider(cfg)
	cache := mongodbcache.NewCache(cfg, logger, dbClient)
	// emailer := mailgun.NewEmailer(cfg, logger)
	// templatedEmailer := templatedemailer.NewTemplatedEmailer(logger, emailer)
	// cloudstore := cloudstorage.NewCloudStorage(cfg, logger)

	//
	// Repository
	//

	tenantRepo := repo.NewTenantRepository(cfg, logger, dbClient)
	userRepo := repo.NewUserRepository(cfg, logger, dbClient)

	//
	// Use-case
	//

	// Tenant
	tenantGetByIDUseCase := usecase.NewTenantGetByIDUseCase(
		cfg,
		logger,
		tenantRepo)

	// User
	userGetByEmailUseCase := usecase.NewUserGetByEmailUseCase(
		cfg,
		logger,
		userRepo)
	userCreateUseCase := usecase.NewUserCreateUseCase(
		cfg,
		logger,
		userRepo)
	userUpdateUseCase := usecase.NewUserUpdateUseCase(
		cfg,
		logger,
		userRepo)

	//
	// Service
	//
	userCreateService := service.NewUserCreateService(
		cfg,
		logger,
		passp,
		cache,
		jwtp,
		tenantGetByIDUseCase,
		userGetByEmailUseCase,
		userCreateUseCase,
		userUpdateUseCase,
	)

	// Minor formatting of input.
	pass, err := sstring.NewSecureString(flagPassword)
	if err != nil {
		log.Fatalf("Failed securing: %v\n", err)
	}
	// defer pass.Wipe() // Developers Note: Commented out b/c they are causing the hang in the program to exit?
	passConfirm, err := sstring.NewSecureString(flagPasswordConfirm)
	if err != nil {
		log.Fatalf("Failed securing: %v\n", err)
	}
	// defer passRepeated.Wipe() // Developers Note: Commented out b/c they are causing the hang in the program to exit?
	if pass.String() != passConfirm.String() {
		log.Fatalf("doRunNewAccount(): Passwords do not match: pass: %s passConfirm: %s\n", pass, passConfirm)
	}

	////
	//// Start the transaction.
	////

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	session, err := dbClient.StartSession()
	if err != nil {
		logger.Error("start session error",
			slog.Any("error", err))
		log.Fatalf("Failed executing: %v\n", err)
	}
	defer session.EndSession(ctx)

	logger.Debug("Starting MongoDB transaction...")

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		logger.Debug("Transaction started")
		req := &service.UserCreateRequestIDO{
			FirstName:           flagFirstName,
			LastName:            flagLastName,
			Email:               flagEmail,
			Password:            flagPassword,
			PasswordConfirm:     flagPasswordConfirm,
			Phone:               flagPhone,
			Country:             flagCountry,
			Timezone:            flagTimezone,
			AgreeTermsOfService: flagAgreeTermsOfService,
			AgreePromotions:     flagAgreePromotions,
			Role:                flagRole,
		}

		// Execution
		account, err := userCreateService.Execute(sessCtx, req)
		if err != nil {
			sessCtx.AbortTransaction(ctx)
			return nil, err
		}
		if account == nil {
			err := errors.New("Account does not exist")
			sessCtx.AbortTransaction(ctx)
			return nil, err
		}

		logger.Debug("Committing transaction")
		if err := sessCtx.CommitTransaction(ctx); err != nil {
			return nil, err
		}
		logger.Debug("Transaction committed")

		return account, nil
	}
	res, err := session.WithTransaction(ctx, transactionFunc)
	logger.Debug("Finished mongodb transaction") //HELP: WHY DOES THE CODE NEVER GET HERE?
	if err != nil {
		logger.Error("session failed error",
			slog.Any("error", err))
		log.Fatalf("Failed creating account: %v\n", err)
	}

	resp := res.(*service.UserCreateResponseIDO)

	logger.Debug("Account created",
		slog.Any("resp", resp))
}
