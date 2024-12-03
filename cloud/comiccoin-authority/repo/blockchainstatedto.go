package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

const (
	blockchainStateURL string = "/api/v1/blockchain-state?chain_id=${CHAIN_ID}"
)

type BlockchainStateDTOConfigurationProvider interface {
	GetAuthorityAddress() string
}

type blockchainStateDTOConfigurationProviderImpl struct {
	authorityAddress string
}

func NewBlockchainStateDTOConfigurationProvider(authorityAddress string) BlockchainStateDTOConfigurationProvider {
	return &blockchainStateDTOConfigurationProviderImpl{
		authorityAddress: authorityAddress,
	}
}

func (impl *blockchainStateDTOConfigurationProviderImpl) GetAuthorityAddress() string {
	return impl.authorityAddress
}

type BlockchainStateDTORepo struct {
	config BlockchainStateDTOConfigurationProvider
	logger *slog.Logger
}

func NewBlockchainStateDTORepo(
	config BlockchainStateDTOConfigurationProvider,
	logger *slog.Logger,
) *BlockchainStateDTORepo {

	return &BlockchainStateDTORepo{
		config: config,
		logger: logger,
	}
}

func (repo *BlockchainStateDTORepo) GetFromBlockchainAuthorityByChainID(ctx context.Context, chainID uint16) (*domain.BlockchainStateDTO, error) {
	modifiedURL := strings.ReplaceAll(blockchainStateURL, "${CHAIN_ID}", fmt.Sprintf("%v", chainID))
	httpEndpoint := fmt.Sprintf("%s%s", repo.config.GetAuthorityAddress(), modifiedURL)

	r, err := http.NewRequest("GET", httpEndpoint, nil)
	if err != nil {
		repo.logger.Debug("failed to setup get request",
			slog.Any("error", err))
		return nil, err
	}

	r.Header.Add("Content-Type", "application/json")

	repo.logger.Debug("Submitting to HTTP JSON API",
		slog.String("url", httpEndpoint),
		slog.String("method", "GET"))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		repo.logger.Debug("failed to do post request",
			slog.Any("error", err))
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		err := fmt.Errorf("http endpoint does not exist for: %v", httpEndpoint)
		repo.logger.Debug("failed to do post request",
			slog.Any("error", err))
		return nil, err
	}

	if res.StatusCode == http.StatusBadRequest {
		e := make(map[string]string)
		var rawJSON bytes.Buffer
		teeReader := io.TeeReader(res.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

		// Try to decode the response as a string first
		var jsonStr string
		err := json.NewDecoder(teeReader).Decode(&jsonStr)
		if err != nil {
			repo.logger.Error("decoding string error",
				slog.Any("err", err),
				slog.String("json", rawJSON.String()),
			)
			return nil, err
		}

		// Now try to decode the string into a map
		err = json.Unmarshal([]byte(jsonStr), &e)
		if err != nil {
			repo.logger.Error("decoding map error",
				slog.Any("err", err),
				slog.String("json", jsonStr),
			)
			return nil, err
		}

		repo.logger.Debug("Parsed error response",
			slog.Any("errors", e),
		)
		return nil, err
	}

	var rawJSON bytes.Buffer
	teeReader := io.TeeReader(res.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

	respPayload := &domain.BlockchainStateDTO{}
	if err := json.NewDecoder(teeReader).Decode(&respPayload); err != nil {
		repo.logger.Error("decoding string error",
			slog.Any("err", err),
			slog.String("json", rawJSON.String()),
		)
		return nil, err
	}

	repo.logger.Debug("Blockchain state retrieved") // slog.Any("chain_id", respPayload.ChainID),
	// slog.Any("latest_block_number", respPayload.LatestBlockNumber),
	// slog.Any("latest_hash", respPayload.LatestHash),
	// slog.Any("latest_token_id", respPayload.LatestTokenID),
	// slog.Any("account_hash_state", respPayload.AccountHashState),
	// slog.Any("token_hash_state", respPayload.TokenHashState),

	return respPayload, nil
}
