package controller

import (
	"context"
	"log/slog"
	"time"

	// attachment_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/attachment/datastore"
	submission_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	// credit_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	// receipt_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/datastore"
	// store_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	// userpurchase_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/datastore"
)

func (impl *CustomerControllerImpl) updateRelatedComicsInBackground(u *user_s.User) error {
	ctx := context.Background() // Execute in background and not in foreground.

	////
	//// CASE 1: Related by `customer_id`.
	////

	f := &submission_s.ComicSubmissionPaginationListFilter{
		Cursor:     "",
		CustomerID: u.ID,
		PageSize:   1_000_000_000,
		SortField:  "created_at",
		SortOrder:  -1,
	}
	cscs, err := impl.ComicSubmissionStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, cs := range cscs.Results {
		cs.CustomerFirstName = u.FirstName
		cs.CustomerLastName = u.LastName
		cs.ModifiedAt = time.Now()
		if err := impl.ComicSubmissionStorer.UpdateByID(ctx, cs); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated comic submission",
			slog.Any("ID", u.ID),
			slog.Any("StoreName", u.StoreName))
	}

	return nil
}
