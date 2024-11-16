package repo

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

const (
	blockchainStateUpdateChangeEventsURL string = "/api/v1/blockchain-state/changes?chain_id=${CHAIN_ID}"
)

type BlockchainStateChangeEventDTOConfigurationProvider interface {
	GetAuthorityAddress() string
}

type blockchainStateChangeEventDTOConfigurationProviderImpl struct {
	authorityAddress string
}

func NewBlockchainStateChangeEventDTOConfigurationProvider(authorityAddress string) BlockchainStateChangeEventDTOConfigurationProvider {
	return &blockchainStateChangeEventDTOConfigurationProviderImpl{
		authorityAddress: authorityAddress,
	}
}

func (impl *blockchainStateChangeEventDTOConfigurationProviderImpl) GetAuthorityAddress() string {
	return impl.authorityAddress
}

type BlockchainStateChangeEventDTORepo struct {
	config BlockchainStateChangeEventDTOConfigurationProvider
	logger *slog.Logger
}

func NewBlockchainStateChangeEventDTORepo(
	config BlockchainStateChangeEventDTOConfigurationProvider,
	logger *slog.Logger,
) domain.BlockchainStateChangeEventDTORepository {

	return &BlockchainStateChangeEventDTORepo{
		config: config,
		logger: logger,
	}
}

func (repo *BlockchainStateChangeEventDTORepo) SubscribeToBlockchainAuthority(ctx context.Context, chainID uint16) (<-chan uint16, error) {
	modifiedURL := strings.ReplaceAll(blockchainStateUpdateChangeEventsURL, "${CHAIN_ID}", fmt.Sprintf("%v", chainID))
	httpEndpoint := fmt.Sprintf("%s:%s", repo.config.GetAuthorityAddress(), modifiedURL)

	// Make the HTTP request
	resp, err := http.Get(httpEndpoint)
	if err != nil {
		repo.logger.Error("failed to make HTTP request", slog.Any("err", err))
		return nil, err
	}

	// Ensure the response body is closed when the function exits
	go func() {
		<-ctx.Done()
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		repo.logger.Error("non-OK HTTP status code", slog.Any("status_code", resp.StatusCode))
		return nil, fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	// Create the output channel
	outputChan := make(chan uint16)

	// Start a goroutine to read from the HTTP response body
	go func() {
		defer close(outputChan)

		buf := make([]byte, 4096)
		for {
			select {
			case <-ctx.Done():
				repo.logger.Info("context cancelled, exiting")
				return
			default:
				// Read from the response body
				n, err := resp.Body.Read(buf)
				if err != nil {
					if err.Error() == "EOF" {
						repo.logger.Info("EOF reached, exiting")
					} else {
						repo.logger.Error("error reading response body", slog.Any("err", err))
					}
					return
				}

				// Process the received chunk
				str := string(buf[:n])
				lines := strings.Split(str, "\n")

				for _, line := range lines {
					// Parse "data:" lines
					if strings.HasPrefix(line, "data:") {
						data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
						repo.logger.Info("received data", slog.Any("data", data))

						// Convert the data to uint16
						value, err := strconv.ParseUint(data, 10, 16)
						if err != nil {
							repo.logger.Error("failed to parse uint16", slog.Any("data", data), slog.Any("err", err))
							continue
						}

						// Send the value to the output channel
						select {
						case outputChan <- uint16(value):
						case <-ctx.Done():
							repo.logger.Info("context cancelled, exiting")
							return
						}
					}
				}
			}
		}
	}()

	return outputChan, nil
}
