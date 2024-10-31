package repo

import (
	"fmt"
	"log/slog"

	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type NonFungibleTokenRepo struct {
	logger                  *slog.Logger
	dbByTokenIDClient       disk.Storage
	dbByMetadataURIIDClient disk.Storage
}

func NewNonFungibleTokenRepo(logger *slog.Logger, dbByTokenIDClient disk.Storage, dbByMetadataURIIDClient disk.Storage) *NonFungibleTokenRepo {
	return &NonFungibleTokenRepo{logger, dbByTokenIDClient, dbByMetadataURIIDClient}
}

func (r *NonFungibleTokenRepo) Upsert(token *domain.NonFungibleToken) error {
	bBytes, err := token.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbByTokenIDClient.Set(fmt.Sprintf("%v", token.TokenID), bBytes); err != nil {
		return err
	}
	if err := r.dbByMetadataURIIDClient.Set(token.MetadataURI, bBytes); err != nil {
		return err
	}
	return nil
}

func (r *NonFungibleTokenRepo) GetByTokenID(tokenID uint64) (*domain.NonFungibleToken, error) {
	bBytes, err := r.dbByTokenIDClient.Get(fmt.Sprintf("%v", tokenID))
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
	err := r.dbByTokenIDClient.Iterate(func(key, value []byte) error {
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

func (r *NonFungibleTokenRepo) DeleteByTokenID(tokenID uint64) error {
	token, err := r.GetByTokenID(tokenID)
	if err != nil {
		return err
	}
	if err := r.dbByTokenIDClient.Delete(fmt.Sprintf("%v", tokenID)); err != nil {
		return err
	}
	if err := r.dbByMetadataURIIDClient.Delete(token.MetadataURI); err != nil {
		return err
	}

	return nil
}

func (r *NonFungibleTokenRepo) OpenTransaction() error {
	if err := r.dbByTokenIDClient.OpenTransaction(); err != nil {
		return err
	}
	if err := r.dbByMetadataURIIDClient.OpenTransaction(); err != nil {
		return err
	}
	return nil
}

func (r *NonFungibleTokenRepo) CommitTransaction() error {
	if err := r.dbByTokenIDClient.CommitTransaction(); err != nil {
		return err
	}
	if err := r.dbByMetadataURIIDClient.CommitTransaction(); err != nil {
		return err
	}
	return nil
}

func (r *NonFungibleTokenRepo) DiscardTransaction() {
	r.dbByTokenIDClient.DiscardTransaction()
	r.dbByMetadataURIIDClient.DiscardTransaction()
}
