package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/memory/redis"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

// BlockchainStateChangeSubscriptionService is responsible for subscribing to our
// remote Redis server and wait to receive published events from the consensus
// mechanism; in addition, this struct will take the received event and return
// the blockchain state struct. This service is responsible for continously
// maintaining a connection with Redis and accepting published events and
// then returning latest blockchain states.
type BlockchainStateChangeSubscriptionService struct {
	logger                          *slog.Logger
	blockchainStateSubscribeUseCase *usecase.BlockchainStateSubscribeUseCase
	blockchainStateSubscriber       redis.RedisSubscriber
}

func NewBlockchainStateChangeSubscriptionService(
	logger *slog.Logger,
	uc1 *usecase.BlockchainStateSubscribeUseCase,
) *BlockchainStateChangeSubscriptionService {
	subscriber := uc1.Execute(context.Background())
	return &BlockchainStateChangeSubscriptionService{logger, uc1, subscriber}
}

func (s *BlockchainStateChangeSubscriptionService) Execute(ctx context.Context) (*domain.BlockchainState, error) {
	ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)
	s.logger.Debug("Waiting to receive latest blockchain state changes from redis...",
		slog.Any("ip_address", ipAddress))
	bcStateBin, err := s.blockchainStateSubscriber.WaitUntilReceiveMessage(ctx)
	if err != nil {
		s.logger.Error("failed to wait and receive message from blockchain state subscriber",
			slog.Any("ip_address", ipAddress),
			slog.Any("error", err))
		return nil, err
	}
	if bcStateBin == nil {
		err := errors.New("Did not receive any blockchain state change event")
		s.logger.Error("failed to wait and receive message from blockchain state subscriber",
			slog.Any("ip_address", ipAddress),
			slog.Any("error", err))
		return nil, err
	}
	bcState, err := domain.NewBlockchainStateFromDeserialize(bcStateBin)
	if err != nil {
		s.logger.Error("failed to deserialize latest blockchain state from subscriber",
			slog.Any("ip_address", ipAddress),
			slog.Any("error", err))
		return nil, err
	}
	s.logger.Debug("Successfully received latest blockchain state change from redis",
		slog.Any("ip_address", ipAddress),
		slog.Any("chain_id", bcState.ChainID),
	)
	return bcState, nil
}

func (s *BlockchainStateChangeSubscriptionService) Terminate(ctx context.Context) error {
	ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)
	s.logger.Debug("Gracefully shutting down blockchain state change subscription service...",
		slog.Any("ip_address", ipAddress))
	if s.blockchainStateSubscriber != nil {
		if err := s.blockchainStateSubscriber.Close(); err != nil {
			s.logger.Error("failed to close blockchain state subscriber",
				slog.Any("ip_address", ipAddress),
				slog.Any("error", err))
			return err
		}
		s.blockchainStateSubscriber = nil
		s.logger.Debug("Successfully shut down blockchain state change subscription service.",
			slog.Any("ip_address", ipAddress))
	}
	return nil
}
