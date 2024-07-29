package templatedemailer

import (
	"log/slog"

	mg "github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/emailer/mailgun"

	c "github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
)

// TemplatedEmailer Is adapter for responsive HTML email templates sender.
type TemplatedEmailer interface {
	GetBackendDomainName() string
	GetFrontendDomainName() string
	SendNewUserTemporaryPasswordEmail(email, firstName, temporaryPassword string) error
	SendBusinessVerificationEmail(email, verificationCode, firstName string) error
	SendCustomerVerificationEmail(email, verificationCode, firstName string) error
	SendForgotPasswordEmail(email, verificationCode, firstName string) error
	SendNewComicSubmissionEmailToStaff(staffEmails []string, submissionID string, storeName string, item string, cpsrn string, serviceTypeName string) error
	SendNewComicSubmissionEmailToRetailers(retailerEmails []string, submissionID string, storeName string, item string, cpsrn string, serviceTypeName string) error
	SendNewStoreEmailToStaff(staffEmails []string, storeID string) error
	SendRetailerStoreActiveEmailToRetailers(retailerEmails []string, storeName string) error
}

type templatedEmailer struct {
	UUID    uuid.Provider
	Logger  *slog.Logger
	Emailer mg.Emailer
}

func NewTemplatedEmailer(cfg *c.Conf, logger *slog.Logger, uuidp uuid.Provider, emailer mg.Emailer) TemplatedEmailer {
	// Defensive code: Make sure we have access to the file before proceeding any further with the code.
	logger.Debug("templated emailer initializing...")
	logger.Debug("templated emailer initialized")

	return &templatedEmailer{
		UUID:    uuidp,
		Logger:  logger,
		Emailer: emailer,
	}
}

func (impl *templatedEmailer) GetBackendDomainName() string {
	return impl.Emailer.GetBackendDomainName()
}

func (impl *templatedEmailer) GetFrontendDomainName() string {
	return impl.Emailer.GetFrontendDomainName()
}
