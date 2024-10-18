package repo

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"
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
			slog.Any("account_id", address),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *WalletRepo) DeleteByAddress(address *common.Address) error {
	err := r.dbClient.Delete(address.String())
	if err != nil {
		return err
	}
	return nil
}