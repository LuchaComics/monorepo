package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/domain"
)

const (
	blockchainStateURL string = "/api/v1/blockchain-state?chain_id=${CHAIN_ID}"
)

type BlockchainStateDTORepo struct {
	config *config.Config
	logger *slog.Logger
}

func NewBlockchainStateDTORepo(cfg *config.Config, logger *slog.Logger) domain.BlockchainStateDTORepository {

	return &BlockchainStateDTORepo{cfg, logger}
}

func (repo *BlockchainStateDTORepo) GetFromCentralAuthorityByChainID(ctx context.Context, chainID uint16) (*domain.BlockchainStateDTO, error) {
	modifiedAccountDetailURL := strings.ReplaceAll(blockchainStateURL, "${CHAIN_ID}", fmt.Sprintf("%v", repo.config.Blockchain.ChainID))
	httpEndpoint := fmt.Sprintf("%s:%s", repo.config.Blockchain.CentralAuthorityHTTPAddress, modifiedAccountDetailURL)

	r, err := http.NewRequest("GET", httpEndpoint, nil)
	if err != nil {
		log.Fatalf("failed to setup get request: %v", err)
	}

	r.Header.Add("Content-Type", "application/json")

	repo.logger.Debug("Submitting to HTTP JSON API",
		slog.String("url", httpEndpoint),
		slog.String("method", "GET"))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Fatalf("failed to do post request: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		log.Fatalf("http endpoint does not exist for: %v", httpEndpoint)
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

	repo.logger.Debug("Blockchain state retrieved",
		slog.Any("chain_id", respPayload.ChainID),
		slog.Any("latest_block_number", respPayload.LatestBlockNumber),
		slog.Any("latest_hash", respPayload.LatestHash),
		slog.Any("latest_token_id", respPayload.LatestTokenID),
		slog.Any("account_hash_state", respPayload.AccountHashState),
		slog.Any("token_hash_state", respPayload.TokenHashState),
	)

	return respPayload, nil
}
