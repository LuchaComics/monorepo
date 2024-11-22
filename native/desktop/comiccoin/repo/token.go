package repo

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"sort"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/signature"
	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/ethereum/go-ethereum/common"
)

type TokenRepo struct {
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewTokenRepo(logger *slog.Logger, db disk.Storage) domain.TokenRepository {
	return &TokenRepo{logger, db}
}

func (r *TokenRepo) Upsert(ctx context.Context, token *domain.Token) error {
	bBytes, err := token.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(fmt.Sprintf("%v", token.IDBytes), bBytes); err != nil {
		return err
	}
	return nil
}

func (r *TokenRepo) GetByID(ctx context.Context, id *big.Int) (*domain.Token, error) {
	bBytes, err := r.dbClient.Get(fmt.Sprintf("%v", id.Bytes()))
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

func (r *TokenRepo) ListByChainID(ctx context.Context, chainID uint16) ([]*domain.Token, error) {
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

		if token.ChainID == chainID {
			res = append(res, token)
		}

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (r *TokenRepo) ListByOwner(ctx context.Context, owner *common.Address) ([]*domain.Token, error) {
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

		if strings.ToLower(token.Owner.Hex()) == strings.ToLower(owner.Hex()) {
			res = append(res, token)
		}

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (r *TokenRepo) CountByOwner(ctx context.Context, owner *common.Address) (int64, error) {
	var counter int64
	err := r.dbClient.Iterate(func(key, value []byte) error {
		token, err := domain.NewTokenFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		if strings.ToLower(token.Owner.Hex()) == strings.ToLower(owner.Hex()) {
			counter++
		}

		// Return nil to indicate success
		return nil
	})

	return counter, err
}

func (r *TokenRepo) DeleteByID(ctx context.Context, id *big.Int) error {
	err := r.dbClient.Delete(fmt.Sprintf("%v", id))
	if err != nil {
		return err
	}
	return nil
}

func (r *TokenRepo) HashStateByChainID(ctx context.Context, chainID uint16) (string, error) {
	tokens, err := r.ListByChainID(ctx, chainID)
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

		// For defensive purposes, only append by chain ID.
		if tok.ChainID == chainID {
			tokensBytes = append(tokensBytes, tokBytes...)
		}
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
	return ba[i].GetID().Cmp(ba[j].GetID()) < 0
}

// Swap moves tokens in the order of the token id value.
func (ba byToken) Swap(i, j int) {
	ba[i], ba[j] = ba[j], ba[i]
}
