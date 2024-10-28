package repo

import (
	"fmt"
	"log/slog"

	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/domain"
)

type NFTRepo struct {
	logger                  *slog.Logger
	dbByTokenIDClient       disk.Storage
	dbByMetadataURIIDClient disk.Storage
}

func NewNFTRepo(logger *slog.Logger, dbByTokenIDClient disk.Storage, dbByMetadataURIIDClient disk.Storage) *NFTRepo {
	return &NFTRepo{logger, dbByTokenIDClient, dbByMetadataURIIDClient}
}

func (r *NFTRepo) Upsert(nft *domain.NFT) error {
	bBytes, err := nft.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbByTokenIDClient.Set(fmt.Sprintf("%v", nft.TokenID), bBytes); err != nil {
		return err
	}
	if err := r.dbByMetadataURIIDClient.Set(nft.MetadataURI, bBytes); err != nil {
		return err
	}
	return nil
}

func (r *NFTRepo) GetByTokenID(tokenID uint64) (*domain.NFT, error) {
	bBytes, err := r.dbByTokenIDClient.Get(fmt.Sprintf("%v", tokenID))
	if err != nil {
		return nil, err
	}
	b, err := domain.NewNFTFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.Any("token_id", tokenID),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *NFTRepo) ListAll() ([]*domain.NFT, error) {
	res := make([]*domain.NFT, 0)
	err := r.dbByTokenIDClient.Iterate(func(key, value []byte) error {
		nft, err := domain.NewNFTFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		res = append(res, nft)

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (r *NFTRepo) DeleteByTokenID(tokenID uint64) error {
	nft, err := r.GetByTokenID(tokenID)
	if err != nil {
		return err
	}
	if err := r.dbByTokenIDClient.Delete(fmt.Sprintf("%v", tokenID)); err != nil {
		return err
	}
	if err := r.dbByMetadataURIIDClient.Delete(nft.MetadataURI); err != nil {
		return err
	}

	return nil
}

func (r *NFTRepo) OpenTransaction() error {
	if err := r.dbByTokenIDClient.OpenTransaction(); err != nil {
		return err
	}
	if err := r.dbByMetadataURIIDClient.OpenTransaction(); err != nil {
		return err
	}
	return nil
}

func (r *NFTRepo) CommitTransaction() error {
	if err := r.dbByTokenIDClient.CommitTransaction(); err != nil {
		return err
	}
	if err := r.dbByMetadataURIIDClient.CommitTransaction(); err != nil {
		return err
	}
	return nil
}

func (r *NFTRepo) DiscardTransaction() {
	r.dbByTokenIDClient.DiscardTransaction()
	r.dbByMetadataURIIDClient.DiscardTransaction()
}
