package repo

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db"
)

type WalletRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient dbase.Database
}

func NewWalletRepo(cfg *config.Config, logger *slog.Logger, db dbase.Database) *WalletRepo {
	return &WalletRepo{cfg, logger, db}
}

func (r *WalletRepo) Upsert(wallet *domain.Wallet) error {
	bBytes, err := wallet.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Setf(bBytes, "wallet-%v", wallet.AccountID); err != nil {
		return err
	}
	return nil
}

func (r *WalletRepo) GetByAccountID(accountID string) (*domain.Wallet, error) {
	bBytes, err := r.dbClient.Getf("wallet-%v", accountID)
	if err != nil {
		return nil, err
	}
	b, err := domain.NewWalletFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.String("account_id", accountID),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *WalletRepo) DeleteByAccountID(accountID string) error {
	err := r.dbClient.Deletef("wallet-%v", accountID)
	if err != nil {
		return err
	}
	return nil
}
