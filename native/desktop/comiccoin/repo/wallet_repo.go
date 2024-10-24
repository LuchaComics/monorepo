package repo

import (
	"log/slog"

	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/ethereum/go-ethereum/common"
)

type WalletRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewWalletRepo(cfg *config.Config, logger *slog.Logger, db disk.Storage) *WalletRepo {
	return &WalletRepo{cfg, logger, db}
}

func (r *WalletRepo) Upsert(wallet *domain.Wallet) error {
	bBytes, err := wallet.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(wallet.Address.String(), bBytes); err != nil {
		return err
	}
	return nil
}

func (r *WalletRepo) GetByAddress(address *common.Address) (*domain.Wallet, error) {
	bBytes, err := r.dbClient.Get(address.String())
	if err != nil {
		return nil, err
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

func (r *WalletRepo) ListAll() ([]*domain.Wallet, error) {
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

func (r *WalletRepo) DeleteByAddress(address *common.Address) error {
	err := r.dbClient.Delete(address.String())
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
