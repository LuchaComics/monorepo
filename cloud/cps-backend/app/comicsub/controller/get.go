package controller

import (
	"context"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (c *ComicSubmissionControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.ComicSubmission, error) {
	// DEVELOPERS NOTE:
	// Every submission creation is dependent on the `role` of the logged in
	// user in our system so we need to extract it right away.
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Retrieve from our database the record for the specific id.
	m, err := c.ComicSubmissionStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}

	// Variable controls if we need to update the record before we return it.
	var hasUpdate bool = false

	// The following will generate a pre-signed URL so user can download the it
	// if the object key exists. But before generating the pre-signed URL if
	// the expiry date has elapsed.
	if m.FindingsFormObjectKey != "" {
		nowt := time.Now() // Get the current datetime.

		if nowt.After(m.FindingsFormObjectURLExpiry) {
			// The following will generate a pre-signed URL so user can download the file.
			expiryDate := time.Now().Add(time.Minute * 15)
			downloadableURL, err := c.S3.GetDownloadablePresignedURL(ctx, m.FindingsFormObjectKey, time.Minute*15)
			if err != nil {
				c.Logger.Warn("s3 presign error", slog.Any("error", err))
				// Do not return an error, simply continue this function as there might
				// be a case were the file was removed on the s3 bucket by ourselves
				// or some other reason.
				return m, nil
			}
			m.FindingsFormObjectURL = downloadableURL
			m.FindingsFormObjectURLExpiry = expiryDate
			m.ModifiedAt = time.Now()
			hasUpdate = true
		}
	}
	if m.LabelObjectKey != "" {
		nowt := time.Now() // Get the current datetime.

		if nowt.After(m.LabelObjectURLExpiry) {
			// The following will generate a pre-signed URL so user can download the file.
			expiryDate := time.Now().Add(time.Minute * 15)
			downloadableURL, err := c.S3.GetDownloadablePresignedURL(ctx, m.LabelObjectKey, time.Minute*15)
			if err != nil {
				c.Logger.Warn("s3 presign error", slog.Any("error", err))
				// Do not return an error, simply continue this function as there might
				// be a case were the file was removed on the s3 bucket by ourselves
				// or some other reason.
				return m, nil
			}
			m.LabelObjectURL = downloadableURL
			m.LabelObjectURLExpiry = expiryDate
			m.ModifiedAt = time.Now()
			hasUpdate = true
		}
	}

	if hasUpdate {
		c.Logger.Debug("has update when getting")
		if err := c.ComicSubmissionStorer.UpdateByID(ctx, m); err != nil {
			c.Logger.Error("database update error", slog.Any("error", err))
			return nil, err
		}
	}

	//
	// Security - Censor label data if the logged in user is retailer. We do
	//            this because if the retailer gets our label then they can
	//            print it themeselves!
	//

	switch userRole {
	case u_d.UserRoleRetailer, u_d.UserRoleCustomer:
		m.LabelObjectKey = "[hidden]"
		m.LabelObjectURL = "[hidden]"
		m.LabelObjectURLExpiry = time.Now()
	}

	return m, err
}

func (c *ComicSubmissionControllerImpl) GetByCPSRN(ctx context.Context, cpsrn string) (*domain.ComicSubmission, error) {
	// Retrieve from our database the record for the specific cspn.
	m, err := c.ComicSubmissionStorer.GetByCPSRN(ctx, cpsrn)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		c.Logger.Warn("submission registry does not exist for cpsrn lookup validation error", slog.String("cpsrn", cpsrn))
		return nil, httperror.NewForBadRequestWithSingleField("message", "registry entry does not exist")
	}

	return m, err
}
