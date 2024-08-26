package controller

import (
	"context"
	"log/slog"
	"time"

	tenant_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/tenant/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
)

func (impl *UserControllerImpl) updateRelatedStoreInBackground(u *user_s.User) error {
	ctx := context.Background() // Execute in background and not in foreground.

	////
	//// CASE 1: Related by `created_by_user_id`.
	////

	f := &tenant_s.TenantPaginationListFilter{
		Cursor:          "",
		CreatedByUserID: u.ID,
		PageSize:        1_000_000_000,
		SortField:       "created_at",
		SortOrder:       1,
	}
	ss, err := impl.TenantStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, s := range ss.Results {
		s.CreatedByUserName = u.Name
		s.ModifiedAt = time.Now()
		if err := impl.TenantStorer.UpdateByID(ctx, s); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated store",
			slog.Any("storet_id", s.ID),
			slog.Any("tenant_name", s.Name))
	}

	////
	//// CASE 2: Related by `modified_by_user_id`.
	////

	f = &tenant_s.TenantPaginationListFilter{
		Cursor:          "",
		CreatedByUserID: u.ID,
		PageSize:        1_000_000_000,
		SortField:       "created_at",
		SortOrder:       1,
	}
	ss, err = impl.TenantStorer.ListByFilter(ctx, f)
	if err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}
	for _, s := range ss.Results {
		s.ModifiedByUserName = u.Name
		s.ModifiedAt = time.Now()
		if err := impl.TenantStorer.UpdateByID(ctx, s); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("Updated store",
			slog.Any("storet_id", s.ID),
			slog.Any("tenant_name", s.Name))
	}

	return nil
}
