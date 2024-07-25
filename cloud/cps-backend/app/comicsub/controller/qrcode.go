package controller

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/skip2/go-qrcode"

	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (c *ComicSubmissionControllerImpl) GetQRCodePNGImage(ctx context.Context, payload string) ([]byte, error) {
	// Generate the QR code for the specific URL and return the `png` binary
	// file bytes.
	var png []byte
	png, err := qrcode.Encode(payload, qrcode.Medium, 256)

	return png, err
}

func (c *ComicSubmissionControllerImpl) GetQRCodePNGImageOfRegisteryURLByCPSRN(ctx context.Context, cpsrn string) ([]byte, error) {
	// Retrieve from our database the record for the specific id.
	submission, err := c.ComicSubmissionStorer.GetByCPSRN(ctx, cpsrn)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if submission == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", fmt.Sprintf("does not exist for: %v", cpsrn))
	}

	// Create our payload.
	payload := fmt.Sprintf("https://%s/cpsrn?v=%s", c.Emailer.GetFrontendDomainName(), cpsrn)

	// Generate the QR code for the specific URL and return the `png` binary
	// file bytes.
	var png []byte
	png, err = qrcode.Encode(payload, qrcode.Medium, 256)
	if err != nil {
		c.Logger.Error("encode error", slog.Any("error", err))
		return nil, err
	}

	c.Logger.Debug("qr code ready",
		slog.Any("payload", payload),
		slog.Any("cpsrn", cpsrn))

	return png, err
}
