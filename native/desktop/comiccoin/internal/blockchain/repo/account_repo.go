package repo

import (
	"log/slog"
	"sort"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db"
)

type AccountRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient dbase.Database
}

func NewAccountRepo(cfg *config.Config, logger *slog.Logger, db dbase.Database) *AccountRepo {
	return &AccountRepo{cfg, logger, db}
}

func (r *AccountRepo) Upsert(account *domain.Account) error {
	bBytes, err := account.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Setf(bBytes, "account-%v", account.ID); err != nil {
		return err
	}
	return nil
}

func (r *AccountRepo) GetByID(id string) (*domain.Account, error) {
	bBytes, err := r.dbClient.Getf("account-%v", id)
	if err != nil {
		return nil, err
	}
	b, err := domain.NewAccountFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.String("id", id),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *AccountRepo) ListAll() ([]*domain.Account, error) {
	res := make([]*domain.Account, 0)
	seekThenIterateKey := ""
	err := r.dbClient.Iterate("account-", seekThenIterateKey, func(key, value []byte) error {
		account, err := domain.NewAccountFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		res = append(res, account)

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (r *AccountRepo) DeleteByID(id string) error {
	err := r.dbClient.Deletef("account-%v", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *AccountRepo) HashState() (string, error) {
	accounts, err := r.ListAll()
	if err != nil {
		return "", err
	}

	sort.Sort(byAccount(accounts))
	return signature.Hash(accounts), nil
}

// =============================================================================

// byAccount provides sorting support by the account id value.
type byAccount []*domain.Account

// Len returns the number of transactions in the list.
func (ba byAccount) Len() int {
	return len(ba)
}

// Less helps to sort the list by account id in ascending order to keep the
// accounts in the right order of processing.
func (ba byAccount) Less(i, j int) bool {
	return ba[i].ID < ba[j].ID
}

// Swap moves accounts in the order of the account id value.
func (ba byAccount) Swap(i, j int) {
	ba[i], ba[j] = ba[j], ba[i]
}
