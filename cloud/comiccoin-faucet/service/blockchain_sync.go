package service

import (
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/blockchain/signature"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BlockchainSyncWithBlockchainAuthorityService struct {
	logger                                               *slog.Logger
	getGenesisBlockDataUseCase                           *usecase.GetGenesisBlockDataUseCase
	upsertGenesisBlockDataUseCase                        *usecase.UpsertGenesisBlockDataUseCase
	getGenesisBlockDataDTOFromBlockchainAuthorityUseCase *usecase.GetGenesisBlockDataDTOFromBlockchainAuthorityUseCase
	getBlockchainStateUseCase                            *usecase.GetBlockchainStateUseCase
	upsertBlockchainStateUseCase                         *usecase.UpsertBlockchainStateUseCase
	getBlockchainStateDTOFromBlockchainAuthorityUseCase  *usecase.GetBlockchainStateDTOFromBlockchainAuthorityUseCase
	getBlockDataUseCase                                  *usecase.GetBlockDataUseCase
	upsertBlockDataUseCase                               *usecase.UpsertBlockDataUseCase
	getBlockDataDTOFromBlockchainAuthorityUseCase        *usecase.GetBlockDataDTOFromBlockchainAuthorityUseCase
	getAccountUseCase                                    *usecase.GetAccountUseCase
	upsertAccountUseCase                                 *usecase.UpsertAccountUseCase
	upsertTokenIfPreviousTokenNonceGTEUseCase            *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase
	tenantGetByIDUseCase                                 *usecase.TenantGetByIDUseCase
	tenantUpdateUseCase                                  *usecase.TenantUpdateUseCase
}

func NewBlockchainSyncWithBlockchainAuthorityService(
	logger *slog.Logger,
	uc1 *usecase.GetGenesisBlockDataUseCase,
	uc2 *usecase.UpsertGenesisBlockDataUseCase,
	uc3 *usecase.GetGenesisBlockDataDTOFromBlockchainAuthorityUseCase,
	uc4 *usecase.GetBlockchainStateUseCase,
	uc5 *usecase.UpsertBlockchainStateUseCase,
	uc6 *usecase.GetBlockchainStateDTOFromBlockchainAuthorityUseCase,
	uc7 *usecase.GetBlockDataUseCase,
	uc8 *usecase.UpsertBlockDataUseCase,
	uc9 *usecase.GetBlockDataDTOFromBlockchainAuthorityUseCase,
	uc10 *usecase.GetAccountUseCase,
	uc11 *usecase.UpsertAccountUseCase,
	uc12 *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase,
	uc13 *usecase.TenantGetByIDUseCase,
	uc14 *usecase.TenantUpdateUseCase,
) *BlockchainSyncWithBlockchainAuthorityService {
	return &BlockchainSyncWithBlockchainAuthorityService{logger, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8, uc9, uc10, uc11, uc12, uc13, uc14}
}

func (s *BlockchainSyncWithBlockchainAuthorityService) Execute(sessCtx mongo.SessionContext, chainID uint16, tenantID primitive.ObjectID) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if chainID == 0 {
		e["chainID"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Validation failed",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// Step 2:
	// Get our genesis block, and if it doesn't exist then we need to fetch it
	// from the blockchain faucet for the particular `chainID`.
	//

	genesis, err := s.getGenesisBlockDataUseCase.Execute(sessCtx, chainID)
	if err != nil {
		s.logger.Error("Failed getting genesis block locally",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return err
	}
	if genesis == nil {
		s.logger.Debug("Genesis block d.n.e, fetching it now ...")
		genesisDTO, err := s.getGenesisBlockDataDTOFromBlockchainAuthorityUseCase.Execute(sessCtx, chainID)
		if err != nil {
			s.logger.Error("Failed getting genesis block remotely",
				slog.Any("chain_id", chainID),
				slog.Any("error", err))
			return err
		}
		if genesisDTO == nil {
			err := fmt.Errorf("Genesis block data does not exist for `chain_id`: %v", chainID)
			s.logger.Error("Failed getting genesis block remotely",
				slog.Any("chain_id", chainID),
				slog.Any("error", err))
			return err
		}

		// Convert from network format data to our local format.
		genesis = domain.GenesisBlockDataDTOToGenesisBlockData(genesisDTO)

		// Save the genesis block data to local database.
		if err := s.upsertGenesisBlockDataUseCase.Execute(sessCtx, genesis.Hash, genesis.Header, genesis.HeaderSignatureBytes, genesis.Trans, genesis.Validator); err != nil {
			s.logger.Error("Failed upserting genesis",
				slog.Any("chain_id", chainID),
				slog.Any("error", err))
			return err
		}

		// Process transactions for the genesis block.
		genesisCoinTx := genesis.Trans[0]
		genesisTokenTx := genesis.Trans[1]

		if err := s.upsertAccountUseCase.Execute(sessCtx, genesisCoinTx.From, genesisCoinTx.Value, big.NewInt(0)); err != nil {
			s.logger.Error("Failed upserting coinbase account.",
				slog.Any("error", err))
			return err
		}

		// Save our token to the local database ONLY if this transaction
		// is the most recent one. We track "most recent" transaction by
		// the nonce value in the token.
		upsertTokenErr := s.upsertTokenIfPreviousTokenNonceGTEUseCase.Execute(
			sessCtx,
			genesisTokenTx.GetTokenID(),
			genesisTokenTx.To,
			genesisTokenTx.TokenMetadataURI,
			genesisTokenTx.GetTokenNonce())
		if upsertTokenErr != nil {
			s.logger.Error("Failed upserting genesis block token",
				slog.Any("type", genesisTokenTx.Type),
				slog.Any("token_id", genesisTokenTx.GetTokenID()),
				slog.Any("owner", genesisTokenTx.To),
				slog.Any("metadataURI", genesisTokenTx.TokenMetadataURI),
				slog.Any("nonce", genesisTokenTx.GetTokenNonce()),
				slog.Any("error", upsertTokenErr))
			return upsertTokenErr
		}

		s.logger.Debug("Genesis block saved to local database from global blockchain network",
			slog.Any("chain_id", chainID),
			slog.Any("coinbase_address", genesisCoinTx.From.Hex()))
	}

	//
	// STEP 3:
	// Get the blockchain state we have *locally* and *remotely* and compare
	// the differences, if our local blockchain state matches what is on the
	// global blockchain network then we are done synching (because there is
	// nothin left to sync). If we don't even have a blockchain state then we need to
	// proceed to download the entire blockchain immediately. If there is any
	// discrepency between the global and local state then we proceed with
	// this function and update our local blockchain with the available data
	// on the global blockchain network.
	//

	globalBlockchainStateDTO, err := s.getBlockchainStateDTOFromBlockchainAuthorityUseCase.Execute(sessCtx, chainID)
	if err != nil {
		s.logger.Error("Failed getting global blockchain state",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return err
	}
	if globalBlockchainStateDTO == nil {
		err := fmt.Errorf("Failed getting global blockchain state for chainID: %v", chainID)
		s.logger.Error("Failed getting global blockchain state",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return err
	}
	// Convert from network format data to our local format.
	globalBlockchainState := domain.BlockchainStateDTOToBlockchainState(globalBlockchainStateDTO)

	// Fetch our local blockchain state.
	localBlockchainState, err := s.getBlockchainStateUseCase.Execute(sessCtx, chainID)
	if err != nil {
		s.logger.Error("Failed getting local blockchain state",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return err
	}

	// If our local blockchain state is empty then create it with using the genesis block.
	if localBlockchainState == nil {
		localBlockchainState = &domain.BlockchainState{
			ChainID:                genesis.Header.ChainID,
			LatestBlockNumberBytes: genesis.Header.NumberBytes,
			LatestHash:             genesis.Hash,
			LatestTokenIDBytes:     genesis.Header.LatestTokenIDBytes,
			TransactionFee:         genesis.Header.TransactionFee,
			AccountHashState:       genesis.Header.StateRoot,
			TokenHashState:         genesis.Header.TokensRoot,
		}
		if err := s.upsertBlockchainStateUseCase.Execute(sessCtx, localBlockchainState); err != nil {
			s.logger.Error("Failed upserting local blockchain state from genesis block data",
				slog.Any("chain_id", chainID),
				slog.Any("error", err))
			return err
		}
		s.logger.Debug("Local blockchain state set to genesis block data",
			slog.Any("chain_id", chainID))
	} else {
		if localBlockchainState.LatestHash == globalBlockchainState.LatestHash {
			s.logger.Debug("Local blockchain is in sync with global blockchain network",
				slog.Any("chain_id", chainID))
			return nil
		}
		s.logger.Debug("Local blockchain state is out of sync with global blockchain network",
			slog.Any("chain_id", chainID))
	}

	//
	// STEP 4:
	// Proceed to download all the missing block data from the global blockchain
	// network so our local blockchain will be in-sync.
	//

	if err := s.syncWithGlobalBlockchainNetwork(sessCtx, localBlockchainState, globalBlockchainState); err != nil {
		if localBlockchainState.LatestHash == globalBlockchainState.LatestHash {
			s.logger.Debug("Failed to sync with the global blockchain network",
				slog.Any("chain_id", chainID))
			return nil
		}
	}

	//
	// STEP 5:
	// Update our blockchain state to match the global blockchain network's state.
	//

	if err := s.upsertBlockchainStateUseCase.Execute(sessCtx, globalBlockchainState); err != nil {
		s.logger.Error("Failed upserting global blockchain state.",
			slog.Any("chain_id", chainID),
			slog.Any("error", err))
		return err
	}

	s.logger.Debug("Local blockchain was synced successfully with the global blockchain network",
		slog.Any("chain_id", chainID))

	//
	// STEP 6:
	// Update our faucet.
	//

	tenant, err := s.tenantGetByIDUseCase.Execute(sessCtx, tenantID)
	if err != nil {
		s.logger.Error("Failed getting tenant",
			slog.Any("chain_id", chainID),
			slog.Any("tenant_id", tenantID),
			slog.Any("error", err))
		return err
	}
	if tenant == nil {
		err := fmt.Errorf("Tenant does not exist at id: %v", tenantID)
		s.logger.Error("Failed getting tenant",
			slog.Any("chain_id", chainID),
			slog.Any("tenant_id", tenantID),
			slog.Any("error", err))
		return err
	}
	recentTenantAccount, err := s.getAccountUseCase.Execute(sessCtx, tenant.Account.Address)
	if err != nil {
		s.logger.Error("Failed getting recent tenant address",
			slog.Any("chain_id", chainID),
			slog.Any("tenant_id", tenantID),
			slog.Any("error", err))
		return err
	}
	if recentTenantAccount == nil {
		err := fmt.Errorf("Recent tenant address does not exist at address: %v", tenant.Account.Address)
		s.logger.Error("Failed getting tenant address",
			slog.Any("chain_id", chainID),
			slog.Any("tenant_id", tenantID),
			slog.Any("error", err))
		return err
	}
	tenant.Account = recentTenantAccount
	if err := s.tenantUpdateUseCase.Execute(sessCtx, tenant); err != nil {
		s.logger.Error("Failed updating tenant",
			slog.Any("chain_id", chainID),
			slog.Any("tenant_id", tenantID),
			slog.Any("error", err))
		return err
	}

	s.logger.Debug("Tenant account balance updated",
		slog.Any("chain_id", chainID),
		slog.Any("tenant_id", tenantID),
		slog.Any("balance", tenant.Account.Balance))

	return nil
}

func (s *BlockchainSyncWithBlockchainAuthorityService) syncWithGlobalBlockchainNetwork(sessCtx mongo.SessionContext, localBlockchainState, globalBlockchainState *domain.BlockchainState) error {
	s.logger.Debug("Beginning to sync with global blockchain network...")

	// Variable holds the position of the block data we are looking at in
	// the global blockchain network.
	currentHashIterator := globalBlockchainState.LatestHash

	// Variable used to trigger the sync operation.
	var isSyncOperationRunning bool = true

	// Continue to sync with global blockchain network until our
	// sync operation finishes.
	for isSyncOperationRunning {

		//
		// STEP 1: Fetch from the global blockchain network.
		//

		s.logger.Debug("Fetching block data from global blockchain network...",
			slog.Any("hash", currentHashIterator))

		blockDataDTO, err := s.getBlockDataDTOFromBlockchainAuthorityUseCase.Execute(sessCtx, currentHashIterator)
		if err != nil {
			s.logger.Debug("Failed to get block data from global blockchain network",
				slog.Any("hash", currentHashIterator))
			return err
		}

		// Artifical delay as to not overload the network resources.
		time.Sleep(1 * time.Second)

		blockData := domain.BlockDataDTOToBlockData(blockDataDTO)

		//
		// STEP 2: Save it to the local database.
		//

		s.logger.Debug("Downloaded block data from global blockchain network and saving to local database...",
			slog.Any("hash", currentHashIterator))

		if err := s.upsertBlockDataUseCase.Execute(sessCtx, blockData.Hash, blockData.Header, blockData.HeaderSignatureBytes, blockData.Trans, blockData.Validator); err != nil {
			s.logger.Debug("Failed to upsert block data ",
				slog.Any("hash", currentHashIterator))
			return err
		}

		//
		// STEP 3:
		// Process account coins and tokens from the transactions.
		//

		for _, blockTx := range blockData.Trans {

			// STEP 4:
			// Process accounts.
			//

			s.logger.Debug("Processing block tx...",
				slog.Any("type", blockTx.Type),
				slog.Any("nonce", blockTx.GetNonce()),
				slog.Any("timestamp", blockTx.TimeStamp))
			if err := s.processAccountForTransaction(sessCtx, blockData, &blockTx); err != nil {
				s.logger.Error("Failed processing transaction",
					slog.Any("error", err))
				return err
			}

			//
			// STEP 5:
			// Process tokens.
			//
			if blockTx.Type == domain.TransactionTypeToken {
				// Save our token to the local database ONLY if this transaction
				// is the most recent one. We track "most recent" transaction by
				// the nonce value in the token.
				err := s.upsertTokenIfPreviousTokenNonceGTEUseCase.Execute(
					sessCtx,
					blockTx.GetTokenID(),
					blockTx.To,
					blockTx.TokenMetadataURI,
					blockTx.GetTokenNonce())
				if err != nil {
					s.logger.Error("Failed upserting (if previous token nonce GTE then current)",
						slog.Any("type", blockTx.Type),
						slog.Any("token_id", blockTx.GetTokenID()),
						slog.Any("owner", blockTx.To),
						slog.Any("metadataURI", blockTx.TokenMetadataURI),
						slog.Any("nonce", blockTx.GetTokenNonce()),
						slog.Any("error", err))
					return err
				}
			}

			s.logger.Debug("Finished processing block tx",
				slog.Any("type", blockTx.Type),
				slog.Any("nonce", blockTx.GetNonce()),
				slog.Any("timestamp", blockTx.TimeStamp))
		}

		//
		// STEP 6:
		// Check to see if we haven't reached the last block data we have
		// in our local blockchain.
		//

		currentHashIterator = blockData.Header.PrevBlockHash
		isSyncOperationRunning = localBlockchainState.LatestHash != currentHashIterator
		isSyncOperationRunning = isSyncOperationRunning && (currentHashIterator != signature.ZeroHash) // consensus mechanism reached genesis block data, sync completed

		s.logger.Debug("Processed block data",
			slog.String("currentHashIterator", currentHashIterator),
			slog.Bool("isSyncOperationRunning", isSyncOperationRunning))
	}
	s.logger.Debug("Finished syncing with global blockchain network")
	return nil
}

func (s *BlockchainSyncWithBlockchainAuthorityService) processAccountForTransaction(sessCtx mongo.SessionContext, blockData *domain.BlockData, blockTx *domain.BlockTransaction) error {
	//
	// CASE 1 OF 2: 🎟️ Token Transaction
	//

	if blockTx.Type == domain.TransactionTypeToken {
		return s.processAccountForTokenTransaction(sessCtx, blockData, blockTx)
	}

	//
	// CASE 2 OF 2: 🪙 Coin Transaction
	//

	if blockTx.Type == domain.TransactionTypeCoin {
		return s.processAccountForCoinTransaction(sessCtx, blockData, blockTx)
	}

	return nil
}

func (s *BlockchainSyncWithBlockchainAuthorityService) processAccountForCoinTransaction(sessCtx mongo.SessionContext, blockData *domain.BlockData, blockTx *domain.BlockTransaction) error {
	//
	// STEP 1
	//

	if blockTx.From != nil {
		// DEVELOPERS NOTE:
		// We already *should* have a `From` account in our database, so we can
		acc, _ := s.getAccountUseCase.Execute(sessCtx, blockTx.From)
		if acc == nil {
			s.logger.Error("The `From` account does not exist in our database.",
				slog.Any("hash", blockTx.From))
			return fmt.Errorf("The `From` account does not exist in our database for hash: %v", blockTx.From.String())
		}

		acc.Balance -= blockTx.Value

		// Note: We do this to prevent reply attacks. (See notes in either `domain/accounts.go` or `service/genesis_init.go`)
		noince := acc.GetNonce()
		noince.Add(noince, big.NewInt(1))
		acc.NonceBytes = noince.Bytes()

		// DEVELOPERS NOTE:
		// Do not update this accounts `Nonce`, we need to only update the
		// `Nonce` to the receiving account, i.e. the `To` account.

		if err := s.upsertAccountUseCase.Execute(sessCtx, acc.Address, acc.Balance, acc.GetNonce()); err != nil {
			s.logger.Error("Failed upserting account.",
				slog.Any("error", err))
			return err
		}

		s.logger.Debug("New `From` account balance via blockchain faucet",
			slog.Any("account_address", acc.Address),
			slog.Any("balance", acc.Balance),
		)
	}

	//
	// STEP 2
	//

	if blockTx.To != nil {

		// Variable holds the value of coins to transfer to the account after
		// the transaction fee was collected by the Authority.
		var valueMinusFees uint64 = blockTx.Value - blockData.Header.TransactionFee

		acc, _ := s.getAccountUseCase.Execute(sessCtx, blockTx.To)
		if acc == nil {
			if err := s.upsertAccountUseCase.Execute(sessCtx, blockTx.To, 0, big.NewInt(0)); err != nil {
				s.logger.Error("Failed creating account.",
					slog.Any("error", err))
				return err
			}
			acc = &domain.Account{
				Address: blockTx.To,

				// Always start by zero, increment by 1 after mining successful.
				NonceBytes: big.NewInt(0).Bytes(),

				Balance: valueMinusFees,
			}
		} else {
			acc.Balance += valueMinusFees

			// Note: We do this to prevent reply attacks. (See notes in either `domain/accounts.go` or `service/genesis_init.go`)
			noince := acc.GetNonce()
			noince.Add(noince, big.NewInt(1))
			acc.NonceBytes = noince.Bytes()
		}

		if err := s.upsertAccountUseCase.Execute(sessCtx, acc.Address, acc.Balance, acc.GetNonce()); err != nil {
			s.logger.Error("Failed upserting account.",
				slog.Any("error", err))
			return err
		}

		s.logger.Debug("New `To` account balance via blockchain faucet",
			slog.Any("account_address", acc.Address),
			slog.Any("balance", acc.Balance),
		)
	}

	//
	// STEP 3
	//

	// Deposit the transaction fee back to the coinbase to be recirculated.
	proofOfAuthorityAccount, err := s.getAccountUseCase.Execute(sessCtx, &blockData.Header.Beneficiary)
	if err != nil {
		s.logger.Error("Failed getting proof of authority account.",
			slog.Any("error", err))
		return err
	}
	if proofOfAuthorityAccount == nil {
		s.logger.Error("Proof of authority account does not exist")
		return fmt.Errorf("Proof of authority account does not exist")
	}

	// Collect transaction fee from this coin transaction.
	proofOfAuthorityAccount.Balance += blockData.Header.TransactionFee

	// Note: We do this to prevent reply attacks. (See notes in either `domain/accounts.go` or `service/genesis_init.go`)
	pofNonce := proofOfAuthorityAccount.GetNonce()
	pofNonce.Add(pofNonce, big.NewInt(1))
	proofOfAuthorityAccount.NonceBytes = pofNonce.Bytes()

	if err := s.upsertAccountUseCase.Execute(sessCtx, proofOfAuthorityAccount.Address, proofOfAuthorityAccount.Balance, proofOfAuthorityAccount.GetNonce()); err != nil {
		s.logger.Error("Failed upserting account.",
			slog.Any("error", err))
		return err
	}
	s.logger.Debug("Authority collected transaction fee from coin transfer",
		slog.Any("authority_address", proofOfAuthorityAccount.Address),
		slog.Any("collected_fee", blockData.Header.TransactionFee),
		slog.Any("new_balance", proofOfAuthorityAccount.Balance),
	)

	return nil
}

func (s *BlockchainSyncWithBlockchainAuthorityService) processAccountForTokenTransaction(sessCtx mongo.SessionContext, blockData *domain.BlockData, blockTx *domain.BlockTransaction) error {
	//
	// STEP 1:
	// Check to see if we have an account for this particular token and
	// collect the transaction fee.
	//

	if blockTx.From != nil {
		acc, _ := s.getAccountUseCase.Execute(sessCtx, blockTx.From)
		if acc == nil {
			if err := s.upsertAccountUseCase.Execute(sessCtx, blockTx.To, 0, big.NewInt(0)); err != nil {
				s.logger.Error("Failed creating account.",
					slog.Any("error", err))
				return err
			}
			acc = &domain.Account{
				Address:    blockTx.To,
				NonceBytes: big.NewInt(0).Bytes(), // Always start by zero, increment by 1 after mining successful.
				Balance:    0,
			}
			if err := s.upsertAccountUseCase.Execute(sessCtx, acc.Address, acc.Balance, acc.GetNonce()); err != nil {
				s.logger.Error("Failed upserting account.",
					slog.Any("error", err))
				return err
			}
			s.logger.Debug("New `From` account balance via validator b/c of token",
				slog.Any("account_address", acc.Address),
				slog.Any("balance", acc.Balance),
			)
		} else {
			acc.Balance -= blockTx.Value // Note: The value is equal to the transaction fee.

			// Note: We do this to prevent reply attacks. (See notes in either `domain/accounts.go` or `service/genesis_init.go`)
			accNonce := acc.GetNonce()
			accNonce.Add(accNonce, big.NewInt(1))
			acc.NonceBytes = accNonce.Bytes()

			if err := s.upsertAccountUseCase.Execute(sessCtx, acc.Address, acc.Balance, acc.GetNonce()); err != nil {
				s.logger.Error("Failed upserting account.",
					slog.Any("error", err))
				return err
			}
		}
	}

	//
	// STEP 2:
	// Check to see if we have an account for this particular token, if not
	// then create it.  Do this from the `To` side of the transaction.
	//

	if blockTx.To != nil {
		acc, _ := s.getAccountUseCase.Execute(sessCtx, blockTx.To)
		if acc == nil {
			if err := s.upsertAccountUseCase.Execute(sessCtx, blockTx.To, 0, big.NewInt(0)); err != nil {
				s.logger.Error("Failed creating account.",
					slog.Any("error", err))
				return err
			}
			acc = &domain.Account{
				Address:    blockTx.To,
				NonceBytes: big.NewInt(0).Bytes(), // Always start by zero, increment by 1 after mining successful.
				Balance:    0,
			}
			if err := s.upsertAccountUseCase.Execute(sessCtx, acc.Address, acc.Balance, acc.GetNonce()); err != nil {
				s.logger.Error("Failed upserting account.",
					slog.Any("error", err))
				return err
			}

			s.logger.Debug("New `To` account via validator b/c of token",
				slog.Any("account_address", acc.Address),
				slog.Any("balance", acc.Balance),
			)
		}
	}

	//
	// STEP 3:
	// Deposit the transaction fee back to the coinbase to be recirculated.
	//

	proofOfAuthorityAccount, err := s.getAccountUseCase.Execute(sessCtx, &blockData.Header.Beneficiary)
	if err != nil {
		s.logger.Error("Failed getting proof of authority account.",
			slog.Any("error", err))
		return err
	}
	if proofOfAuthorityAccount == nil {
		s.logger.Error("Proof of authority account does not exist")
		return fmt.Errorf("Proof of authority account does not exist")
	}

	// Collect transaction fee from this token transaction.
	proofOfAuthorityAccount.Balance += blockTx.Value // Note: The value is equal to the transaction fee.

	// Note: We do this to prevent reply attacks. (See notes in either `domain/accounts.go` or `service/genesis_init.go`)
	pofNonce := proofOfAuthorityAccount.GetNonce()
	pofNonce.Add(pofNonce, big.NewInt(1))
	proofOfAuthorityAccount.NonceBytes = pofNonce.Bytes()

	if err := s.upsertAccountUseCase.Execute(sessCtx, proofOfAuthorityAccount.Address, proofOfAuthorityAccount.Balance, proofOfAuthorityAccount.GetNonce()); err != nil {
		s.logger.Error("Failed upserting account.",
			slog.Any("error", err))
		return err
	}
	s.logger.Debug("Authority collected transaction fee from token transfer or burn",
		slog.Any("authority_address", proofOfAuthorityAccount.Address),
		slog.Any("collected_fee", blockData.Header.TransactionFee),
		slog.Any("new_balance", proofOfAuthorityAccount.Balance),
	)

	return nil
}
