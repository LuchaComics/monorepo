package repo

import (
	"log/slog"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/blockchain/signature"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/storage"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/domain"
)

type AccountRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewAccountRepo(cfg *config.Config, logger *slog.Logger, db disk.Storage) *AccountRepo {
	return &AccountRepo{cfg, logger, db}
}

func (r *AccountRepo) Upsert(account *domain.Account) error {
	bBytes, err := account.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(strings.ToLower(account.Address.String()), bBytes); err != nil {
		return err
	}
	return nil
}

func (r *AccountRepo) GetByAddress(addr *common.Address) (*domain.Account, error) {
	bBytes, err := r.dbClient.Get(strings.ToLower(addr.String()))
	if err != nil {
		return nil, err
	}
	b, err := domain.NewAccountFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.Any("addr", addr),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *AccountRepo) ListAll() ([]*domain.Account, error) {
	res := make([]*domain.Account, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
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

func (r *AccountRepo) DeleteByAddress(addr *common.Address) error {
	err := r.dbClient.Delete(strings.ToLower(addr.String()))
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

	// Variable used to only store the accounts which have a balance greater
	// then the value of zero.
	accountsWithBalance := make([]*domain.Account, 0)

	// Iterate through all the accounts and only save the accounts with balance.
	for _, account := range accounts {
		if account.Balance > 0 {
			accountsWithBalance = append(accountsWithBalance, account)
		}
	}

	// Sort the accounts by address
	sort.Sort(byAccount(accountsWithBalance))

	// Serialize the accounts to JSON
	accountsBytes := make([]byte, 0)
	for _, account := range accountsWithBalance {
		// DEVELOPERS NOTE:
		// In Go, the order of struct fields is determined by the order in which
		// they are defined in the struct. However, this order is not guaranteed
		// to be the same across different nodes or even different runs of the
		// same program.
		//
		// To fix this issue, you can use a deterministic serialization
		// algorithm, such as JSON or CBOR, to serialize the Account struct
		// before hashing it. This will ensure that the fields are always
		// serialized in the same order, regardless of the node or run.
		accountBytes, err := account.Serialize()
		if err != nil {
			return "", err
		}
		accountsBytes = append(accountsBytes, accountBytes...)
	}

	// Hash the deterministic serialized accounts
	res := signature.Hash(accountsBytes)
	return res, nil
}

func (r *AccountRepo) OpenTransaction() error {
	return r.dbClient.OpenTransaction()
}

func (r *AccountRepo) CommitTransaction() error {
	return r.dbClient.CommitTransaction()
}

func (r *AccountRepo) DiscardTransaction() {
	r.dbClient.DiscardTransaction()
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
	return strings.ToLower(ba[i].Address.String()) < strings.ToLower(ba[j].Address.String())
}

// Swap moves accounts in the order of the account id value.
func (ba byAccount) Swap(i, j int) {
	ba[i], ba[j] = ba[j], ba[i]
}