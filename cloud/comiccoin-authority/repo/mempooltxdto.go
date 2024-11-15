package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

const (
	mempoolTransactionURL string = "/api/v1/mempool-transactions"
)

type MempoolTransactionDTOConfigurationProvider interface {
	GetAuthorityAddress() string
}

type mempoolTransactionDTOConfigurationProviderImpl struct {
	authorityAddress string
}

func NewMempoolTransactionDTOConfigurationProvider(authorityAddress string) MempoolTransactionDTOConfigurationProvider {
	return &mempoolTransactionDTOConfigurationProviderImpl{
		authorityAddress: authorityAddress,
	}
}

func (impl *mempoolTransactionDTOConfigurationProviderImpl) GetAuthorityAddress() string {
	return impl.authorityAddress
}

type MempoolTransactionDTORepo struct {
	config MempoolTransactionDTOConfigurationProvider
	logger *slog.Logger
}

func NewMempoolTransactionDTORepo(
	config MempoolTransactionDTOConfigurationProvider,
	logger *slog.Logger,
) *MempoolTransactionDTORepo {

	return &MempoolTransactionDTORepo{
		config: config,
		logger: logger,
	}
}

func (repo *MempoolTransactionDTORepo) SubmitToBlockchainAuthority(ctx context.Context, dto *domain.MempoolTransactionDTO) error {
	httpEndpoint := fmt.Sprintf("%s:%s", repo.config.GetAuthorityAddress(), mempoolTransactionURL)
	jsonData, err := json.Marshal(dto)
	if err != nil {
		repo.logger.Error("Marshalling error",
			slog.Any("err", err),
		)
		return err
	}

	repo.logger.Debug("Submitting to HTTP JSON API",
		slog.String("url", httpEndpoint),
		slog.String("method", "POST"))

	resp, err := http.Post(httpEndpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		repo.logger.Error("Post error",
			slog.Any("err", err),
		)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		err := fmt.Errorf("http endpoint does not exist for: %v", httpEndpoint)
		repo.logger.Error("Failed posting to blockchain authority",
			slog.Any("err", err),
		)
		return err
	}

	if resp.StatusCode == http.StatusBadRequest {
		e := make(map[string]string)
		var rawJSON bytes.Buffer
		teeReader := io.TeeReader(resp.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

		// Try to decode the response as a string first
		var jsonStr string
		err := json.NewDecoder(teeReader).Decode(&jsonStr)
		if err != nil {
			repo.logger.Error("decoding string error",
				slog.Any("err", err),
				slog.String("json", rawJSON.String()),
			)
			return err
		}

		// Now try to decode the string into a map
		err = json.Unmarshal([]byte(jsonStr), &e)
		if err != nil {
			repo.logger.Error("decoding map error",
				slog.Any("err", err),
				slog.String("json", jsonStr),
			)
			return err
		}

		repo.logger.Debug("Parsed error response",
			slog.Any("errors", e),
		)
		return err
	}

	repo.logger.Debug("Mempool transaction submitted to blockchain authority",
		slog.Any("url", httpEndpoint),
		slog.Int("status", resp.StatusCode),
	)

	return nil
}
