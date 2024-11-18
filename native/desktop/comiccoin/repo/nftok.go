package repo

import (
	"fmt"
	"log/slog"
	"math/big"

	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type NonFungibleTokenRepo struct {
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewNonFungibleTokenRepo(logger *slog.Logger, dbClient disk.Storage) *NonFungibleTokenRepo {
	return &NonFungibleTokenRepo{logger, dbClient}
}

func (r *NonFungibleTokenRepo) Upsert(token *domain.NonFungibleToken) error {
	bBytes, err := token.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(fmt.Sprintf("%v", token.TokenID), bBytes); err != nil {
		return err
	}
	return nil
}

func (r *NonFungibleTokenRepo) GetByTokenID(tokenID *big.Int) (*domain.NonFungibleToken, error) {
	bBytes, err := r.dbClient.Get(fmt.Sprintf("%v", tokenID))
	if err != nil {
		return nil, err
	}
	b, err := domain.NewNonFungibleTokenFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.Any("token_id", tokenID),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *NonFungibleTokenRepo) ListAll() ([]*domain.NonFungibleToken, error) {
	res := make([]*domain.NonFungibleToken, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		token, err := domain.NewNonFungibleTokenFromDeserialize(value)
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

func (r *NonFungibleTokenRepo) ListWithFilterByTokenIDs(tokIDs []*big.Int) ([]*domain.NonFungibleToken, error) {
	res := make([]*domain.NonFungibleToken, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		token, err := domain.NewNonFungibleTokenFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		// Apply our filter to the results.
		for _, tokID := range tokIDs {
			if tokID == token.TokenID {
				res = append(res, token)
			}
		}

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (r *NonFungibleTokenRepo) DeleteByTokenID(tokenID *big.Int) error {
	if err := r.dbClient.Delete(fmt.Sprintf("%v", tokenID)); err != nil {
		return err
	}
	return nil
}

func (r *NonFungibleTokenRepo) OpenTransaction() error {
	if err := r.dbClient.OpenTransaction(); err != nil {
		return err
	}
	return nil
}

func (r *NonFungibleTokenRepo) CommitTransaction() error {
	if err := r.dbClient.CommitTransaction(); err != nil {
		return err
	}
	return nil
}

func (r *NonFungibleTokenRepo) DiscardTransaction() {
	r.dbClient.DiscardTransaction()
}
