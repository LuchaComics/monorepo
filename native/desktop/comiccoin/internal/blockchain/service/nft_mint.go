package service

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type MintNFTService struct {
	config                                *config.Config
	logger                                *slog.Logger
	loadGenesisBlockDataUseCase           *usecase.LoadGenesisBlockDataUseCase
	getAccountUseCase                     *usecase.GetAccountUseCase
	getWalletUseCase                      *usecase.GetWalletUseCase
	walletDecryptKeyUseCase               *usecase.WalletDecryptKeyUseCase
	broadcastMempoolTransactionDTOUseCase *usecase.BroadcastMempoolTransactionDTOUseCase
}

func NewMintNFTService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.LoadGenesisBlockDataUseCase,
	uc2 *usecase.GetAccountUseCase,
	uc3 *usecase.GetWalletUseCase,
	uc4 *usecase.WalletDecryptKeyUseCase,
	uc5 *usecase.BroadcastMempoolTransactionDTOUseCase,
) *MintNFTService {
	return &MintNFTService{cfg, logger, uc1, uc2, uc3, uc4, uc5}
}

func (s *MintNFTService) Execute(
	ctx context.Context,
	poaAddr *common.Address,
	poaPassword string,
	toAddr *common.Address,
	nftMetadata *domain.NFTMetadata,
	metadataURI string,
) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if poaAddr == nil {
		e["poa_address"] = "missing value"
	}
	if poaPassword == "" {
		e["poa_password"] = "missing value"
	}
	if toAddr == nil {
		e["to"] = "missing value"
	}
	if nftMetadata == nil {
		e["nft_metadata"] = "missing value"
	} else {
		if nftMetadata.Image == "" {
			e["iamge"] = "missing value"
		}
		if nftMetadata.ExternalURL == "" {
			e["external_url"] = "missing value"
		}
		if nftMetadata.Description == "" {
			e["description"] = "missing value"
		}
		if nftMetadata.Name == "" {
			e["name"] = "missing value"
		}
		if len(nftMetadata.Attributes) > 0 {
			// TODO: Impl.
		}
		if nftMetadata.BackgroundColor == "" {
			e["background_color"] = "missing value"
		}
		if nftMetadata.AnimationURL == "" {
			e["animation_url"] = "missing value"
		}
		if nftMetadata.YoutubeURL == "" {
			e["youtube_url"] = "missing value"
		}
	}
	if metadataURI == "" {
		e["metadata_uri"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed creating new account",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get the account and extract the wallet private/public key.
	//

	wallet, err := s.getWalletUseCase.Execute(poaAddr)
	if err != nil {
		s.logger.Error("failed getting from database",
			slog.Any("from_account_address", poaAddr),
			slog.Any("error", err))
		return fmt.Errorf("failed getting from database: %s", err)
	}
	if wallet == nil {
		s.logger.Error("failed getting from database",
			slog.Any("from_account_address", poaAddr),
			slog.Any("error", "d.n.e."))
		return fmt.Errorf("failed getting from database: %s", "wallet d.n.e.")
	}

	key, err := s.walletDecryptKeyUseCase.Execute(wallet.Filepath, poaPassword)
	if err != nil {
		s.logger.Error("failed getting key",
			slog.Any("from_account_address", poaAddr),
			slog.Any("error", err))
		return fmt.Errorf("failed getting key: %s", err)
	}
	if key == nil {
		return fmt.Errorf("failed getting key: %s", "d.n.e.")
	}

	//
	// STEP 3:
	// Verify the account is validator.
	//

	genesisBlockData, err := s.loadGenesisBlockDataUseCase.Execute()
	if err != nil {
		s.logger.Error("failed getting account",
			slog.Any("from_account_address", poaAddr),
			slog.Any("error", err))
		return fmt.Errorf("failed getting account: %s", err)
	}
	if genesisBlockData == nil {
		return fmt.Errorf("failed getting genesis block data: %s", "d.n.e.")
	}
	validator := genesisBlockData.Validator

	publicKeyECDSA, err := validator.GetPublicKeyECDSA()
	if err != nil {
		s.logger.Error("failed unmarshalling validator public key",
			slog.Any("from_account_address", poaAddr),
			slog.Any("error", err))
		return fmt.Errorf("failed unmarshalling validator public key: %s", err)
	}
	log.Println("--->", validator.PublicKeyBytes)
	log.Println("--->", key.PrivateKey.Public())
	log.Println("--->", key.PrivateKey.PublicKey)
	log.Println("--->", publicKeyECDSA)

	// if account.Balance <= value {
	// 	s.logger.Warn("insufficient balance in account",
	// 		slog.Any("account_addr", poaAddr),
	// 		slog.Any("account_balance", account.Balance),
	// 		slog.Any("value", value))
	// 	return fmt.Errorf("insufficient balance: %d", account.Balance)
	// }

	// //
	// // STEP 4
	// // Create our pending transaction and sign it with the accounts private key.
	// //
	//
	// tx := &domain.Transaction{
	// 	ChainID: s.config.Blockchain.ChainID,
	// 	Nonce:   uint64(time.Now().Unix()),
	// 	From:    wallet.Address,
	// 	To:      to,
	// 	Value:   value,
	// 	Data:    data,
	// }
	//
	// stx, signingErr := tx.Sign(key.PrivateKey)
	// if signingErr != nil {
	// 	s.logger.Debug("Failed to sign the transaction",
	// 		slog.Any("error", signingErr))
	// 	return signingErr
	// }
	//
	// s.logger.Debug("Pending transaction signed successfully",
	// 	slog.Uint64("tx_nonce", stx.Nonce))
	//
	// mempoolTx := &domain.MempoolTransaction{
	// 	Transaction: stx.Transaction,
	// 	V:           stx.V,
	// 	R:           stx.R,
	// 	S:           stx.S,
	// }
	//
	// //
	// // STEP 3
	// // Send our pending signed transaction to our distributed mempool nodes
	// // in the blochcian network.
	// //
	//
	// if err := s.broadcastMempoolTransactionDTOUseCase.Execute(ctx, mempoolTx); err != nil {
	// 	s.logger.Error("Failed to broadcast to the blockchain network",
	// 		slog.Any("error", err))
	// 	return err
	// }
	//
	// s.logger.Info("Pending signed transaction submitted to blockchain",
	// 	slog.Uint64("tx_nonce", stx.Nonce))

	return nil
}
