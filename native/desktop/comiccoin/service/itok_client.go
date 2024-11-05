package service

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/kmutexutil"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

// IssuedTokenClientService represents a node service which will wait to
// receive any newly minted non-fungible tokens from the proof of authority.
type IssuedTokenClientService struct {
	config                                    *config.Config
	logger                                    *slog.Logger
	kmutex                                    kmutexutil.KMutexProvider
	receiveIssuedTokenDTOUseCase              *usecase.ReceiveIssuedTokenDTOUseCase
	loadGenesisBlockDataUseCase               *usecase.LoadGenesisBlockDataUseCase
	upsertTokenIfPreviousTokenNonceGTEUseCase *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase
}

func NewIssuedTokenClientService(
	cfg *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.ReceiveIssuedTokenDTOUseCase,
	uc2 *usecase.LoadGenesisBlockDataUseCase,
	uc3 *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase,
) *IssuedTokenClientService {
	return &IssuedTokenClientService{cfg, logger, kmutex, uc1, uc2, uc3}
}

func (s *IssuedTokenClientService) Execute(ctx context.Context) error {

	//
	// STEP 1
	// Wait to receive data (which also was validated) from the P2P network.
	//

	itok, sig, poaValidator, err := s.receiveIssuedTokenDTOUseCase.Execute(ctx)
	if err != nil {
		s.logger.Error("Failed receiving issued token from PoA.",
			slog.Any("error", err))
		return err
	}
	if itok == nil || sig == nil || poaValidator == nil {
		// Developer Note:
		// If we haven't received anything, that means we haven't connected to
		// the distributed / P2P network, so all we can do at the moment is to
		// pause the execution for 1 second and then retry again.
		time.Sleep(1 * time.Second)
		return nil
	}

	s.logger.Info("received itok dto from network",
		slog.Any("itok_id", itok.ID),
		slog.Any("sig", sig),
		slog.Any("val", poaValidator),
	)

	// Lock the mempool's database so we coordinate when we delete the mempool
	// and when we add mempool.
	s.kmutex.Acquire("itok-receive-service")
	defer s.kmutex.Release("itok-receive-service")

	//
	// STEP 2:
	// Confirm the validator of the issued token matches the validator of the
	// genesis block we have.
	//

	genesisBlockData, err := s.loadGenesisBlockDataUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed loading up genesis block from file",
			slog.Any("error", err))
		return fmt.Errorf("Failed loading up genesis block from file: %v", err)
	}
	if genesisBlockData == nil {
		s.logger.Error("genesis block d.n.e.")
		return fmt.Errorf("genesis block does not exist")
	}

	// Developers Note:
	// This is a super important step to enforce the authority being used by
	// the correct party. This code verifies that the the public key of the
	// authority matches the public key set on the genesis block because the
	// user has opend up the actual authorities wallet.
	if bytes.Equal(genesisBlockData.Validator.PublicKeyBytes, poaValidator.PublicKeyBytes) == false {
		s.logger.Error("Failed comparing public keys of validators",
			slog.Any("genesis_val_pk", genesisBlockData.Validator.PublicKeyBytes),
			slog.Any("itok_val_pk", poaValidator.PublicKeyBytes))
		return fmt.Errorf("Failed comparing public keys: %s", "they do not match")
	}

	//
	// STEP 3:
	// Confirm the signature matches the validator's signature.
	//

	if poaValidator.Verify(sig, itok) == false {
		s.logger.Error("validator failed validating: authority signature is invalid")
		return fmt.Errorf("validator failed validating: %v", "authority signature is invalid")
	}

	//
	// STEP 4:
	// Save to our local database.
	//

	if err := s.upsertTokenIfPreviousTokenNonceGTEUseCase.Execute(itok.ID, nil, itok.MetadataURI, 0); err != nil {
		s.logger.Error("Failed saving issued token",
			slog.Any("error", err))
		return err
	}

	s.logger.Info("saved issued token",
		slog.Any("itok_id", itok.ID),
	)

	return nil
}
