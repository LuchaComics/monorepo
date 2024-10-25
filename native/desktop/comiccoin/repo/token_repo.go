package repo

import (
	"fmt"
	"log/slog"
	"sort"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/blockchain/signature"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/ethereum/go-ethereum/common"
)

type TokenRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewTokenRepo(cfg *config.Config, logger *slog.Logger, db disk.Storage) *TokenRepo {
	return &TokenRepo{cfg, logger, db}
}

func (r *TokenRepo) Upsert(token *domain.Token) error {
	bBytes, err := token.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(fmt.Sprintf("%v", token.ID), bBytes); err != nil {
		return err
	}
	return nil
}

func (r *TokenRepo) GetByID(id uint64) (*domain.Token, error) {
	bBytes, err := r.dbClient.Get(fmt.Sprintf("%v", id))
	if err != nil {
		return nil, err
	}
	b, err := domain.NewTokenFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.Any("id", id),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *TokenRepo) ListAll() ([]*domain.Token, error) {
	res := make([]*domain.Token, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		token, err := domain.NewTokenFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		res = append(res, token)

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (r *TokenRepo) ListByOwner(owner *common.Address) ([]*domain.Token, error) {
	res := make([]*domain.Token, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		token, err := domain.NewTokenFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		if token.Owner == owner {
			res = append(res, token)
		}

		// Return nil to indicate success
		return nil
	})

	return res, err
}
func (r *TokenRepo) DeleteByID(id uint64) error {
	err := r.dbClient.Delete(fmt.Sprintf("%v", id))
	if err != nil {
		return err
	}
	return nil
}

func (r *TokenRepo) HashState() (string, error) {
	tokens, err := r.ListAll()
	if err != nil {
		return "", err
	}

	// Sort and hash our tokens.
	sort.Sort(byToken(tokens))

	// Serialize the accounts to JSON
	tokensBytes := make([]byte, 0)
	for _, tok := range tokens {
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
		tokBytes, err := tok.Serialize()
		if err != nil {
			return "", err
		}
		tokensBytes = append(tokensBytes, tokBytes...)
	}

	// Hash the deterministic serialized tokens.
	return signature.Hash(tokensBytes), nil
}

func (r *TokenRepo) OpenTransaction() error {
	return r.dbClient.OpenTransaction()
}

func (r *TokenRepo) CommitTransaction() error {
	return r.dbClient.CommitTransaction()
}

func (r *TokenRepo) DiscardTransaction() {
	r.dbClient.DiscardTransaction()
}

// =============================================================================

// byToken provides sorting support by the token id value.
type byToken []*domain.Token

// Len returns the number of transactions in the list.
func (ba byToken) Len() int {
	return len(ba)
}

// Less helps to sort the list by token id in ascending order to keep the
// tokens in the right order of processing.
func (ba byToken) Less(i, j int) bool {
	return ba[i].ID < ba[j].ID
}

// Swap moves tokens in the order of the token id value.
func (ba byToken) Swap(i, j int) {
	ba[i], ba[j] = ba[j], ba[i]
}
