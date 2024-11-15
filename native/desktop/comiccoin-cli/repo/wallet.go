package repo

import (
	"context"
	"log/slog"
	"strings"

	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"
)

type WalletRepo struct {
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewWalletRepo(logger *slog.Logger, db disk.Storage) *WalletRepo {
	return &WalletRepo{logger, db}
}

func (r *WalletRepo) Upsert(ctx context.Context, wallet *domain.Wallet) error {
	bBytes, err := wallet.Serialize()
	if err != nil {
		return err
	}
	addr := strings.ToLower(wallet.Address.String())
	r.logger.Debug("upserting wallet",
		slog.Any("addr", addr))
	if err := r.dbClient.Set(addr, bBytes); err != nil {
		return err
	}
	return nil
}

func (r *WalletRepo) GetByAddress(ctx context.Context, address *common.Address) (*domain.Wallet, error) {
	// First attempt.
	addr := strings.ToLower(address.String())
	r.logger.Debug("getting by address...",
		slog.Any("addr", addr))
	bBytes, err := r.dbClient.Get(addr)
	if err != nil {
		r.logger.Debug("Failed getting by address",
			slog.Any("addr", addr),
			slog.Any("error", err))
		return nil, err
	}

	// Second attempt.
	if bBytes == nil {
		addr := addr[2:] // Skip the `0x`
		r.logger.Debug("retrying again with modified address...",
			slog.Any("addr", addr))

		bBytes, err = r.dbClient.Get(addr)
		if err != nil {
			r.logger.Debug("Failed getting by address",
				slog.Any("addr", addr),
				slog.Any("error", err))
			return nil, err
		}
	}

	b, err := domain.NewWalletFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.Any("address", address),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *WalletRepo) ListAll(ctx context.Context) ([]*domain.Wallet, error) {
	res := make([]*domain.Wallet, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		wallet, err := domain.NewWalletFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		res = append(res, wallet)

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (r *WalletRepo) ListAllAddresses(ctx context.Context) ([]*common.Address, error) {
	res := make([]*common.Address, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		wallet, err := domain.NewWalletFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		res = append(res, wallet.Address)

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (r *WalletRepo) DeleteByAddress(ctx context.Context, address *common.Address) error {
	err := r.dbClient.Delete(strings.ToLower(address.String()))
	if err != nil {
		return err
	}
	return nil
}

func (r *WalletRepo) OpenTransaction() error {
	return r.dbClient.OpenTransaction()
}

func (r *WalletRepo) CommitTransaction() error {
	return r.dbClient.CommitTransaction()
}

func (r *WalletRepo) DiscardTransaction() {
	r.dbClient.DiscardTransaction()
}
