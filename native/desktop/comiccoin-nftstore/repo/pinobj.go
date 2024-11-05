package repo

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/domain"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"
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
	bBytes, err := pinobj.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbByCIDClient.Set(pinobj.CID, bBytes); err != nil {
		return err
	}
	if err := r.dbByRequestIDClient.Set(fmt.Sprintf("%v", pinobj.RequestID), bBytes); err != nil {
		return err
	}
	return nil
}

func (r *PinObjectRepo) GetByCID(cid string) (*domain.PinObject, error) {
	bBytes, err := r.dbByCIDClient.Get(cid)
	if err != nil {
		return nil, err
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
		return err
	}
	if err := r.dbByCIDClient.Delete(pinobj.CID); err != nil {
		return err
	}
	if err := r.dbByRequestIDClient.Delete(fmt.Sprintf("%v", requestID)); err != nil {
		return err
	}

	return nil
}

func (r *PinObjectRepo) DeleteByCID(cid string) error {
	pinobj, err := r.GetByCID(cid)
	if err != nil {
		return err
	}
	if err := r.dbByCIDClient.Delete(cid); err != nil {
		return err
	}
	if err := r.dbByRequestIDClient.Delete(fmt.Sprintf("%v", pinobj.RequestID)); err != nil {
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
