package mailgun

import (
	"context"
	"time"

	"log/slog"

	"github.com/mailgun/mailgun-go/v4"

	c "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/provider/uuid"
)

type Emailer interface {
	Send(ctx context.Context, sender, subject, recipient, htmlContent string) error
	GetSenderEmail() string
	GetDomainName() string // Deprecated
	GetBackendDomainName() string
	GetFrontendDomainName() string
	GetMaintenanceEmail() string
}

type mailgunEmailer struct {
	Mailgun          *mailgun.MailgunImpl
	UUID             uuid.Provider
	Logger           *slog.Logger
	senderEmail      string
	apiDomainName    string
	appDomainName    string
	maintenanceEmail string
}

func NewEmailer(cfg *c.Conf, logger *slog.Logger, uuidp uuid.Provider) Emailer {
	// Defensive code: Make sure we have access to the file before proceeding any further with the code.
	logger.Debug("mailgun emailer initializing...")
	mg := mailgun.NewMailgun(cfg.Emailer.Domain, cfg.Emailer.APIKey)
	logger.Debug("mailgun emailer was initialized.")

	mg.SetAPIBase(cfg.Emailer.APIBase) // Override to support our custom email requirements.

	return &mailgunEmailer{
		Mailgun:          mg,
		UUID:             uuidp,
		Logger:           logger,
		senderEmail:      cfg.Emailer.SenderEmail,
		apiDomainName:    cfg.AppServer.APIDomainName,
		appDomainName:    cfg.AppServer.AppDomainName,
		maintenanceEmail: cfg.Emailer.MaintenanceEmail,
	}
}

func (me *mailgunEmailer) Send(ctx context.Context, sender, subject, recipient, body string) error {
	me.Logger.Debug("sent email",
		slog.String("sender", sender),
		slog.String("subject", subject),
		slog.String("recipient", recipient))

	message := me.Mailgun.NewMessage(sender, subject, "", recipient)
	message.SetHtml(body)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// Send the message with a 10 second timeout
	_, id, err := me.Mailgun.Send(ctx, message)

	if err != nil {
		me.Logger.Error("emailer failed sending", slog.Any("err", err))
		return err
	}

	me.Logger.Debug("emailer sent with response", slog.Any("response id", id))

	return nil
}

func (me *mailgunEmailer) GetDomainName() string {
	return me.appDomainName
}

func (me *mailgunEmailer) GetSenderEmail() string {
	return me.senderEmail
}

func (me *mailgunEmailer) GetBackendDomainName() string {
	return me.apiDomainName
}

func (me *mailgunEmailer) GetFrontendDomainName() string {
	return me.appDomainName
}

func (me *mailgunEmailer) GetMaintenanceEmail() string {
	return me.maintenanceEmail
}
