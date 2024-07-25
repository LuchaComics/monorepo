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

func (impl *StoreControllerImpl) updateRelatedUsersInBackground(ns *store_s.Store) error {
	ctx := context.Background() // Execute in background and not in foreground.

	f := &user_s.UserPaginationListFilter{
		Cursor:    "",
		StoreID:   ns.ID,
		PageSize:  1_000_000_000,
		SortField: "created_at",
		SortOrder: 1,
	}
	uu, err := impl.UserStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, u := range uu.Results {
		u.StoreName = ns.Name
		u.StoreLevel = ns.Level
		u.ModifiedAt = time.Now()
		if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated user",
			slog.Any("ID", u.ID),
			slog.Any("StoreName", u.StoreName),
			slog.Any("StoreLevel", u.StoreLevel))
	}
	return nil
}

func (impl *StoreControllerImpl) updateRelatedComicSubmissionsInBackground(ns *store_s.Store) error {
	ctx := context.Background() // Execute in background and not in foreground.

	f := &submission_s.ComicSubmissionPaginationListFilter{
		Cursor:    "",
		StoreID:   ns.ID,
		PageSize:  1_000_000_000,
		SortField: "created_at",
		SortOrder: -1,
	}
	uu, err := impl.ComicSubmissionStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, u := range uu.Results {
		u.StoreName = ns.Name
		u.ModifiedAt = time.Now()
		if err := impl.ComicSubmissionStorer.UpdateByID(ctx, u); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated comic submission",
			slog.Any("ID", u.ID),
			slog.Any("StoreName", u.StoreName))
	}
	return nil
}

func (impl *StoreControllerImpl) updateRelatedAttachmentsInBackground(s *store_s.Store) error {
	ctx := context.Background() // Execute in background and not in foreground.

	////
	//// CASE 1: Related by `store_id`.
	////

	f := &attachment_s.AttachmentPaginationListFilter{
		Cursor:    "",
		StoreID:   s.ID,
		PageSize:  1_000_000_000,
		SortField: "created_at",
		SortOrder: attachment_s.SortOrderDescending,
	}
	aa, err := impl.AttachmentStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, a := range aa.Results {
		a.StoreName = s.Name
		a.ModifiedAt = time.Now()
		if err := impl.AttachmentStorer.UpdateByID(ctx, a); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated attachment",
			slog.Any("attachment_id", a.ID))
	}
	return nil
}

func (impl *StoreControllerImpl) updateRelatedCreditsInBackground(s *store_s.Store) error {
	ctx := context.Background() // Execute in background and not in foreground.

	////
	//// CASE 1: Related by `store_id`.
	////

	f := &credit_s.CreditPaginationListFilter{
		Cursor:    "",
		StoreID:   s.ID,
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
		c.StoreName = s.Name
		c.ModifiedAt = time.Now()
		if err := impl.CreditStorer.UpdateByID(ctx, c); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated credit",
			slog.Any("credit_id", c.ID))
	}

	return nil
}

func (impl *StoreControllerImpl) updateRelatedReceiptsInBackground(s *store_s.Store) error {
	ctx := context.Background() // Execute in background and not in foreground.

	////
	//// CASE 1: Related by `store_id`.
	////

	f := &receipt_s.ReceiptPaginationListFilter{
		Cursor:    "",
		StoreID:   s.ID,
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
		c.StoreName = s.Name
		c.ModifiedAt = time.Now()
		if err := impl.ReceiptStorer.UpdateByID(ctx, c); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated receipt",
			slog.Any("receipt_id", c.ID))
	}

	return nil
}

func (impl *StoreControllerImpl) updateRelateUserPurchasesInBackground(s *store_s.Store) error {
	ctx := context.Background() // Execute in background and not in foreground.

	////
	//// CASE 1: Related by `store_id`.
	////

	f := &userpurchase_s.UserPurchasePaginationListFilter{
		Cursor:    "",
		StoreID:   s.ID,
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
		up.StoreName = s.Name
		up.ModifiedAt = time.Now()
		if err := impl.UserPurchaseStorer.UpdateByID(ctx, up); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated user purchase",
			slog.Any("user_purchase_id", up.ID))
	}

	return nil
}
