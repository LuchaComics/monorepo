package repo

import (
	"fmt"
	"log/slog"

	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/storage"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/domain"
)

type PinObjectRepo struct {
	logger              *slog.Logger
	dbByCIDClient       disk.Storage
	dbByRequestIDClient disk.Storage
}

func NewPinObjectRepo(logger *slog.Logger, dbByCIDClient disk.Storage, dbByRequestIDClient disk.Storage) *PinObjectRepo {
	return &PinObjectRepo{logger, dbByCIDClient, dbByRequestIDClient}
}

func (r *PinObjectRepo) Upsert(pinobj *domain.PinObject) error {
	if pinobj == nil {
		r.logger.Warn("Nil detected")
		return nil
	}
	bBytes, err := pinobj.Serialize()
	if err != nil {
		r.logger.Error("Failed serializing pinobject", slog.Any("error", err))
		return err
	}
	if bBytes == nil {
		r.logger.Warn("Nil bytes after serialization detected")
		return nil
	}

	// DEVELOPERS NOTE:
	// We want to make sure the `CID` is always unique but the `RequestID` is
	// always unique on every API post call, therefore if an existing record
	// exists then we will default to use the existing records `RequestID`.
	fetched, err := r.GetByCID(pinobj.CID)
	if fetched != nil && err == nil {
		pinobj.RequestID = fetched.RequestID
	}

	if err := r.dbByCIDClient.Set(pinobj.CID, bBytes); err != nil {
		r.logger.Error("Failed setting by cid", slog.Any("error", err))
		return err
	}
	if err := r.dbByRequestIDClient.Set(fmt.Sprintf("%v", pinobj.RequestID), bBytes); err != nil {
		r.logger.Error("Failed setting by request id", slog.Any("error", err))
		return err
	}
	return nil
}

func (r *PinObjectRepo) GetByCID(cid string) (*domain.PinObject, error) {
	bBytes, err := r.dbByCIDClient.Get(cid)
	if err != nil {
		r.logger.Error("Failed getting from db", slog.Any("error", err))
		return nil, err
	}

	// If nothing exists then simply return nil, do not continue and error.
	if bBytes == nil {
		return nil, nil
	}

	b, err := domain.NewPinObjectFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.Any("cid", cid),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *PinObjectRepo) GetByRequestID(requestID uint64) (*domain.PinObject, error) {
	bBytes, err := r.dbByRequestIDClient.Get(fmt.Sprintf("%v", requestID))
	if err != nil {
		r.logger.Error("Failed getting from db", slog.Any("error", err))
		return nil, err
	}
	b, err := domain.NewPinObjectFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.Any("request_id", requestID),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *PinObjectRepo) ListAll() ([]*domain.PinObject, error) {
	res := make([]*domain.PinObject, 0)
	err := r.dbByCIDClient.Iterate(func(key, value []byte) error {
		pinobj, err := domain.NewPinObjectFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		res = append(res, pinobj)

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (r *PinObjectRepo) DeleteByRequestID(requestID uint64) error {
	pinobj, err := r.GetByRequestID(requestID)
	if err != nil {
		r.logger.Error("Failed getting from db", slog.Any("error", err))
		return err
	}
	if err := r.dbByCIDClient.Delete(pinobj.CID); err != nil {
		r.logger.Error("Failed deleting from db by cid", slog.Any("error", err))
		return err
	}
	if err := r.dbByRequestIDClient.Delete(fmt.Sprintf("%v", requestID)); err != nil {
		r.logger.Error("Failed getting from db by request id", slog.Any("error", err))
		return err
	}

	return nil
}

func (r *PinObjectRepo) DeleteByCID(cid string) error {
	pinobj, err := r.GetByCID(cid)
	r.logger.Error("Failed getting from db", slog.Any("error", err))
	if err != nil {
		return err
	}
	if err := r.dbByCIDClient.Delete(cid); err != nil {
		r.logger.Error("Failed deleting from db by cid", slog.Any("error", err))
		return err
	}
	if err := r.dbByRequestIDClient.Delete(fmt.Sprintf("%v", pinobj.RequestID)); err != nil {
		r.logger.Error("Failed getting from db by request id", slog.Any("error", err))
		return err
	}

	return nil
}

func (r *PinObjectRepo) OpenTransaction() error {
	if err := r.dbByCIDClient.OpenTransaction(); err != nil {
		return err
	}
	if err := r.dbByRequestIDClient.OpenTransaction(); err != nil {
		return err
	}
	return nil
}

func (r *PinObjectRepo) CommitTransaction() error {
	if err := r.dbByCIDClient.CommitTransaction(); err != nil {
		return err
	}
	if err := r.dbByRequestIDClient.CommitTransaction(); err != nil {
		return err
	}
	return nil
}

func (r *PinObjectRepo) DiscardTransaction() {
	r.dbByCIDClient.DiscardTransaction()
	r.dbByRequestIDClient.DiscardTransaction()
}
