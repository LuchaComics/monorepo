package controller

import (
	"context"
	"log/slog"
	"time"

	attachment_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/attachment/datastore"
	submission_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	credit_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	receipt_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/datastore"
	store_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	userpurchase_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/datastore"
)

func (impl *UserControllerImpl) updateRelatedComicsInBackground(u *user_s.User) error {
	ctx := context.Background() // Execute in background and not in foreground.

	////
	//// CASE 1: Related by `inspector_id`.
	////

	f := &submission_s.ComicSubmissionPaginationListFilter{
		Cursor:      "",
		InspectorID: u.ID,
		PageSize:    1_000_000_000,
		SortField:   "created_at",
		SortOrder:   -1,
	}
	cscs, err := impl.ComicSubmissionStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, cs := range cscs.Results {
		cs.InspectorFirstName = u.FirstName
		cs.InspectorLastName = u.LastName
		cs.ModifiedAt = time.Now()
		if err := impl.ComicSubmissionStorer.UpdateByID(ctx, cs); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated comic submission",
			slog.Any("ID", u.ID),
			slog.Any("StoreName", u.StoreName))
	}

	////
	//// CASE 2: Related by `customer_id`.
	////

	f = &submission_s.ComicSubmissionPaginationListFilter{
		Cursor:     "",
		CustomerID: u.ID,
		PageSize:   1_000_000_000,
		SortField:  "created_at",
		SortOrder:  -1,
	}
	cscs, err = impl.ComicSubmissionStorer.ListByFilter(ctx, f)
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

func (impl *UserControllerImpl) updateRelatedAttachmentsInBackground(u *user_s.User) error {
	ctx := context.Background() // Execute in background and not in foreground.

	////
	//// CASE 1: Related by `created_by_user_id`.
	////

	f := &attachment_s.AttachmentPaginationListFilter{
		Cursor:          "",
		CreatedByUserID: u.ID,
		PageSize:        1_000_000_000,
		SortField:       "created_at",
		SortOrder:       attachment_s.SortOrderDescending,
	}
	aa, err := impl.AttachmentStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, a := range aa.Results {
		a.CreatedByUserName = u.Name
		a.ModifiedAt = time.Now()
		if err := impl.AttachmentStorer.UpdateByID(ctx, a); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated attachment",
			slog.Any("attachment_id", a.ID),
			slog.Any("StoreName", u.StoreName))
	}

	////
	//// CASE 2: Related by `modified_by_user_id`.
	////

	f = &attachment_s.AttachmentPaginationListFilter{
		Cursor:          "",
		CreatedByUserID: u.ID,
		PageSize:        1_000_000_000,
		SortField:       "created_at",
		SortOrder:       attachment_s.SortOrderDescending,
	}
	aa, err = impl.AttachmentStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, a := range aa.Results {
		a.ModifiedByUserName = u.Name
		a.ModifiedAt = time.Now()
		if err := impl.AttachmentStorer.UpdateByID(ctx, a); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated attachment",
			slog.Any("attachment_id", a.ID),
			slog.Any("StoreName", u.StoreName))
	}

	return nil
}

func (impl *UserControllerImpl) updateRelatedCreditsInBackground(u *user_s.User) error {
	ctx := context.Background() // Execute in background and not in foreground.

	////
	//// CASE 1: Related by `user_id`.
	////

	f := &credit_s.CreditPaginationListFilter{
		Cursor:    "",
		UserID:    u.ID,
		PageSize:  1_000_000_000,
		SortField: "created_at",
		SortOrder: 1,
	}
	cc, err := impl.CreditStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, c := range cc.Results {
		c.UserName = u.Name
		c.UserLexicalName = u.LexicalName
		c.ModifiedAt = time.Now()
		if err := impl.CreditStorer.UpdateByID(ctx, c); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated credit",
			slog.Any("credit_id", c.ID),
			slog.Any("StoreName", u.StoreName))
	}

	return nil
}

func (impl *UserControllerImpl) updateRelatedReceiptsInBackground(u *user_s.User) error {
	ctx := context.Background() // Execute in background and not in foreground.

	////
	//// CASE 1: Related by `user_id`.
	////

	f := &receipt_s.ReceiptPaginationListFilter{
		Cursor:    "",
		UserID:    u.ID,
		PageSize:  1_000_000_000,
		SortField: "created_at",
		SortOrder: 1,
	}
	cc, err := impl.ReceiptStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, c := range cc.Results {
		c.UserName = u.Name
		c.UserLexicalName = u.LexicalName
		c.ModifiedAt = time.Now()
		if err := impl.ReceiptStorer.UpdateByID(ctx, c); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated receipt",
			slog.Any("receipt_id", c.ID),
			slog.Any("StoreName", u.StoreName))
	}

	return nil
}

func (impl *UserControllerImpl) updateRelatedStoreInBackground(u *user_s.User) error {
	ctx := context.Background() // Execute in background and not in foreground.

	////
	//// CASE 1: Related by `created_by_user_id`.
	////

	f := &store_s.StorePaginationListFilter{
		Cursor:          "",
		CreatedByUserID: u.ID,
		PageSize:        1_000_000_000,
		SortField:       "created_at",
		SortOrder:       1,
	}
	ss, err := impl.StoreStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, s := range ss.Results {
		s.CreatedByUserName = u.Name
		s.ModifiedAt = time.Now()
		if err := impl.StoreStorer.UpdateByID(ctx, s); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated store",
			slog.Any("storet_id", s.ID),
			slog.Any("store_name", s.Name))
	}

	////
	//// CASE 2: Related by `modified_by_user_id`.
	////

	f = &store_s.StorePaginationListFilter{
		Cursor:          "",
		CreatedByUserID: u.ID,
		PageSize:        1_000_000_000,
		SortField:       "created_at",
		SortOrder:       1,
	}
	ss, err = impl.StoreStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, s := range ss.Results {
		s.ModifiedByUserName = u.Name
		s.ModifiedAt = time.Now()
		if err := impl.StoreStorer.UpdateByID(ctx, s); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated store",
			slog.Any("storet_id", s.ID),
			slog.Any("store_name", s.Name))
	}

	return nil
}

func (impl *UserControllerImpl) updateRelateUserPurchasesInBackground(u *user_s.User) error {
	ctx := context.Background() // Execute in background and not in foreground.

	////
	//// CASE 1: Related by `user_id`.
	////

	f := &userpurchase_s.UserPurchasePaginationListFilter{
		Cursor:    "",
		UserID:    u.ID,
		PageSize:  1_000_000_000,
		SortField: "created_at",
		SortOrder: 1,
	}
	upup, err := impl.UserPurchaseStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, up := range upup.Results {
		up.UserName = u.Name
		up.UserLexicalName = u.LexicalName
		up.ModifiedAt = time.Now()
		if err := impl.UserPurchaseStorer.UpdateByID(ctx, up); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated user purchase",
			slog.Any("user_purchase_id", up.ID),
			slog.Any("StoreName", u.StoreName))
	}

	return nil
}
