package controller

import (
	"context"
	"log/slog"
	"time"

	pinobject_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/pinobject/datastore"
	tenant_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
)

func (impl *TenantControllerImpl) updateRelatedUsersInBackground(ns *tenant_s.Tenant) error {
	ctx := context.Background() // Execute in background and not in foreground.

	f := &user_s.UserPaginationListFilter{
		Cursor:    "",
		TenantID:   ns.ID,
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
		u.TenantName = ns.Name
		u.TenantLevel = ns.Level
		u.ModifiedAt = time.Now()
		if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated user",
			slog.Any("ID", u.ID),
			slog.Any("TenantName", u.TenantName),
			slog.Any("TenantLevel", u.TenantLevel))
	}
	return nil
}

func (impl *TenantControllerImpl) updateRelatedPinObjectsInBackground(s *tenant_s.Tenant) error {
	ctx := context.Background() // Execute in background and not in foreground.

	////
	//// CASE 1: Related by `tenant_id`.
	////

	f := &pinobject_s.PinObjectPaginationListFilter{
		Cursor:    "",
		TenantID:   s.ID,
		PageSize:  1_000_000_000,
		SortField: "created_at",
		SortOrder: pinobject_s.SortOrderDescending,
	}
	aa, err := impl.PinObjectStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, a := range aa.Results {
		a.TenantName = s.Name
		a.ModifiedAt = time.Now()
		if err := impl.PinObjectStorer.UpdateByID(ctx, a); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated pinobject",
			slog.Any("pinobject_id", a.ID))
	}
	return nil
}
