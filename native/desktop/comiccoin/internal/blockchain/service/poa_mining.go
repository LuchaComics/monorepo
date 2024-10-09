package service

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/merkle"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
)

type ProofOfAuthorityMiningService struct {
	config                                  *config.Config
	logger                                  *slog.Logger
	kmutex                                  kmutexutil.KMutexProvider
	getKeyService                           *GetKeyService
	getAccountsHashStateUseCase             *usecase.GetAccountsHashStateUseCase
	listAllPendingBlockTransactionUseCase   *usecase.ListAllPendingBlockTransactionUseCase
	getBlockchainLastestHashUseCase         *usecase.GetBlockchainLastestHashUseCase
	setBlockchainLastestHashUseCase         *usecase.SetBlockchainLastestHashUseCase
	getBlockDataUseCase                     *usecase.GetBlockDataUseCase
	createBlockDataUseCase                  *usecase.CreateBlockDataUseCase
	proofOfWorkUseCase                      *usecase.ProofOfWorkUseCase
	broadcastProposedBlockDataDTOUseCase    *usecase.BroadcastProposedBlockDataDTOUseCase
	deleteAllPendingBlockTransactionUseCase *usecase.DeleteAllPendingBlockTransactionUseCase
	getAccountUseCase                       *usecase.GetAccountUseCase
	upsertAccountUseCase                    *usecase.UpsertAccountUseCase
}

func NewProofOfAuthorityMiningService(
	config *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	getKeyService *GetKeyService,
	uc1 *usecase.GetAccountsHashStateUseCase,
	uc2 *usecase.ListAllPendingBlockTransactionUseCase,
	uc3 *usecase.GetBlockchainLastestHashUseCase,
	uc4 *usecase.SetBlockchainLastestHashUseCase,
	uc5 *usecase.GetBlockDataUseCase,
	uc6 *usecase.CreateBlockDataUseCase,
	uc7 *usecase.ProofOfWorkUseCase,
	uc8 *usecase.BroadcastProposedBlockDataDTOUseCase,
	uc9 *usecase.DeleteAllPendingBlockTransactionUseCase,
	uc10 *usecase.GetAccountUseCase,
	uc11 *usecase.UpsertAccountUseCase,
) *ProofOfAuthorityMiningService {
	return &ProofOfAuthorityMiningService{config, logger, kmutex, getKeyService, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8, uc9, uc10, uc11}
}

func (s *ProofOfAuthorityMiningService) Execute(ctx context.Context) error {
	// s.logger.Debug("starting mining service...")
	// defer s.logger.Debug("finished mining service")

	//
	// STEP 1:
	// Lock this function - this is important - because it will fetch all the
	// latest pending block transactions, so there needs to be a lockdown of
	// this service that when it runs it will no longer accept any more calls
	// from the system. Therefore we are using a key-based mutex to lock down
	// this service to act as a singleton runtime usage.
	//

	// Lock the mining service until it has completed executing (or errored).
	s.kmutex.Acquire("mining-service")
	defer s.kmutex.Release("mining-service")

	pendingBlockTxs, err := s.listAllPendingBlockTransactionUseCase.Execute()
	if err != nil {
		s.logger.Error("failed listing pending block transactions",
			slog.Any("error", err))
		return nil
	}
	if len(pendingBlockTxs) <= 0 {
		// s.logger.Debug("skipped mining: has no pending block transactions")
		return nil
	}

	s.logger.Info("executing mining for pending block transactions",
		slog.Int("count", len(pendingBlockTxs)),
	)

	//
	// STEP 2:
	// Lookup the most recent block (data) in our blockchain
	//

	prevBlockDataHash, err := s.getBlockchainLastestHashUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed to get last hash in database",
			slog.Any("error", err))
		return fmt.Errorf("Failed to get last hash in database: %v", err)
	}
	if prevBlockDataHash == "" {
		s.logger.Error("Blockchain not initialized error")
		return fmt.Errorf("Error: %v", "blockchain not initialized")
	}
	prevBlockData, err := s.getBlockDataUseCase.Execute(string(prevBlockDataHash))
	if err != nil {
		s.logger.Error("Error getting block data frin database",
			slog.Any("error", err))
		return fmt.Errorf("Failed to get lastest block data from database: %v", err)
	}
	if prevBlockData == nil {
		s.logger.Error("Block data does not exist in database",
			slog.Any("hash", prevBlockDataHash))
		return fmt.Errorf("Block data does not exist in database for hash: %v", prevBlockDataHash)
	}
	poaValidator := prevBlockData.Validator

	//
	// STEP 3:
	// Setup our new block.
	//

	gasPrice := uint64(s.config.Blockchain.GasPrice)
	unitsOfGas := uint64(s.config.Blockchain.UnitsOfGas)
	trans := make([]domain.BlockTransaction, 0)
	for _, pendingBlockTx := range pendingBlockTxs {
		// Create our block.
		blockTx := domain.BlockTransaction{
			SignedTransaction: *pendingBlockTx.MempoolTransaction.ToSignedTransaction(),
			TimeStamp:         uint64(time.Now().UTC().UnixMilli()),
			GasPrice:          gasPrice,
			GasUnits:          unitsOfGas,
		}
		trans = append(trans, blockTx)
	}

	// Construct a merkle tree from the transaction for this block. The root
	// of this tree will be part of the block to be mined.
	tree, err := merkle.NewTree(trans)
	if err != nil {
		s.logger.Error("Failed to create merkle tree",
			slog.Any("error", err))
		return fmt.Errorf("Failed to create merkle tree: %v", err)
	}

	// Iterate through all the accounts from this local machine, sort them and
	// then hash all of them - this hash represents our `stateRoot` which is
	// in essence a snapshot of the current accounts and their balances. Why is this
	// important?
	//
	// At the start of creating a new block to be mined, a hash of this map is
	// created and stored in the block under the StateRoot field. This allows
	// each node to validate the current state of the peer’s accounts database
	// as part of block validation.
	//
	// It’s critically important that the order of the account balances are
	// exact when hashing the data. The Go spec does not define the order of
	// map iteration and leaves it up to the compiler. Since Go 1.0,
	// the compiler has chosen to have map iteration be random. This function
	//  sorts the accounts and their balances into a slice first and then
	// performs a hash of that slice.
	//
	// When a new block is received, the node can take a hash of their current
	// accounts database and match that to the StateRoot field in the
	// block header. If these hash values don’t match, then there is fraud
	//  going on by the peer and their block would be rejected.	//
	//
	// SPECIAL THANKS TO:
	// https://www.ardanlabs.com/blog/2022/05/blockchain-04-fraud-detection.html
	//
	stateRoot, err := s.getAccountsHashStateUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed to create merkle tree",
			slog.Any("error", err))
		return fmt.Errorf("Failed to create merkle tree: %v", err)
	}

	// Construct the genesis block.
	block := domain.Block{
		Header: &domain.BlockHeader{
			Number:        prevBlockData.Header.Number + 1,
			PrevBlockHash: string(prevBlockDataHash),
			TimeStamp:     uint64(time.Now().UTC().UnixMilli()),
			Beneficiary:   prevBlockData.Header.Beneficiary,
			Difficulty:    s.config.Blockchain.Difficulty,
			MiningReward:  s.config.Blockchain.MiningReward,
			StateRoot:     stateRoot,
			TransRoot:     tree.RootHex(), //
			Nonce:         0,              // Will be identified by the POW algorithm.
		},
		MerkleTree: tree,
	}

	//
	// STEP 4:
	// Execute the proof of work to find our nounce to meet the hash difficulty.
	//

	nonce, powErr := s.proofOfWorkUseCase.Execute(ctx, &block, s.config.Blockchain.Difficulty)
	if powErr != nil {
		s.logger.Error("Failed to mine block",
			slog.Any("error", powErr))
		return fmt.Errorf("Failed to mine block: %v", powErr)
	}

	block.Header.Nonce = nonce

	// Convert into saving for our database and transmitting over network.
	blockData := domain.NewBlockData(block)

	//
	// STEP 5
	// Our proof-of-authority signs this block data's header.
	//

	coinbaseAccountKey, err := s.getKeyService.Execute(
		s.config.Blockchain.ProofOfAuthorityAccountAddress,
		s.config.Blockchain.ProofOfAuthorityWalletPassword)
	if err != nil {
		log.Fatalf("failed getting account wallet key: %v", err)
		return fmt.Errorf("failed getting account wallet key: %v", err)
	}
	if coinbaseAccountKey == nil {
		return fmt.Errorf("failed getting account wallet key: %v", "does not exist")
	}

	coinbasePrivateKey := coinbaseAccountKey.PrivateKey
	blockData.Validator = poaValidator
	blockDataHeaderSignature, err := poaValidator.SignBlockHeader(coinbasePrivateKey, blockData.Header)
	if err != nil {
		return fmt.Errorf("Failed to sign block header: %v", err)
	}
	blockData.HeaderSignature = blockDataHeaderSignature

	//
	// STEP 6:
	// Save to database.
	//

	s.logger.Info("PoA mining completed",
		slog.String("hash", blockData.Hash),
		slog.Uint64("block_number", blockData.Header.Number),
		slog.String("prev_block_hash", blockData.Header.PrevBlockHash),
		slog.Uint64("timestamp", blockData.Header.TimeStamp),
		slog.String("beneficiary", blockData.Header.Beneficiary.String()),
		slog.Uint64("difficulty", uint64(blockData.Header.Difficulty)),
		slog.Uint64("mining_reward", blockData.Header.MiningReward),
		slog.String("state_root", blockData.Header.StateRoot),
		slog.String("trans_root", blockData.Header.TransRoot),
		slog.Uint64("nonce", blockData.Header.Nonce),
		slog.Any("trans", blockData.Trans))

	// Save new block data.
	if err := s.createBlockDataUseCase.Execute(blockData.Hash, blockData.Header, blockData.Trans); err != nil {
		s.logger.Error("Failed to save block data to blockchain",
			slog.Any("error", powErr))
		return fmt.Errorf("Failed to save block data to blockchain: %v", powErr)
	}

	// Delete purposed block data as it has already been processed.
	if err := s.deleteAllPendingBlockTransactionUseCase.Execute(); err != nil {
		s.logger.Error("Failed to deleta all pending block data from blockchain",
			slog.Any("error", powErr))
		return fmt.Errorf("Failed to deleta all pending block data from blockchain: %v", powErr)
	}

	//
	// STEP 7:
	// Broadcast to the distributed / P2P blockchain network our new proposed
	// block data.
	//

	// Convert to our new datastructure.
	purposeBlockData := &domain.ProposedBlockData{
		Hash:   blockData.Hash,
		Header: blockData.Header,
		Trans:  blockData.Trans,
	}

	if err := s.broadcastProposedBlockDataDTOUseCase.Execute(ctx, purposeBlockData); err != nil {
		s.logger.Error("Failed to broadcast proposed block data",
			slog.Any("error", powErr))
		return fmt.Errorf("Failed to broadcast proposed block data: %v", powErr)
	}

	s.logger.Info("PoA mining service broadcasted new block data to propose to the network",
		slog.Uint64("nonce", block.Header.Nonce))

	//
	// STEP 8:
	// Update the account in our in-memory database.
	//

	for _, blockTx := range blockData.Trans {
		if err := s.processAccountForBlockTransaction(blockData, &blockTx); err != nil {
			s.logger.Error("Failed processing transaction",
				slog.Any("error", err))
			return err
		}
	}

	return nil
}

// TODO: (1) Create somesort of `processAccountForBlockTransaction` service and (2) replace it with this.
func (s *ProofOfAuthorityMiningService) processAccountForBlockTransaction(blockData *domain.BlockData, blockTx *domain.BlockTransaction) error {
	// DEVELOPERS NOTE:
	// Please remember that when this function executes, there already is an
	// in-memory database of accounts populated and maintained by this node.
	// Therefore the code in this function is executed on a ready database.

	//
	// STEP 1
	//

	if blockTx.From != nil {
		// DEVELOPERS NOTE:
		// We already *should* have a `From` account in our database, so we can
		acc, _ := s.getAccountUseCase.Execute(blockTx.From)
		if acc == nil {
			s.logger.Error("The `From` account does not exist in our database.",
				slog.Any("hash", blockTx.From))
			return fmt.Errorf("The `From` account does not exist in our database for hash: %v", blockTx.From.String())
		}
		acc.Balance -= blockTx.Value

		// DEVELOPERS NOTE:
		// Do not update this accounts `Nonce`, we need to only update the
		// `Nonce` to the receiving account, i.e. the `To` account.

		if err := s.upsertAccountUseCase.Execute(acc.Address, acc.Balance, acc.Nonce); err != nil {
			s.logger.Error("Failed upserting account.",
				slog.Any("error", err))
			return err
		}

		s.logger.Debug("New `From` account balance via mining",
			slog.Any("account_address", acc.Address),
			slog.Any("balance", acc.Balance),
			slog.Any("tx_hash", blockTx.Hash),
		)
	}

	//
	// STEP 2
	//

	if blockTx.To != nil {
		// DEVELOPERS NOTE:
		// It is perfectly normal that our account would possibly not exist
		// so we would need to create a new Account record in our local
		// in-memory database.
		acc, _ := s.getAccountUseCase.Execute(blockTx.To)
		if acc == nil {
			if err := s.upsertAccountUseCase.Execute(blockTx.To, 0, 0); err != nil {
				s.logger.Error("Failed creating account.",
					slog.Any("error", err))
				return err
			}
			acc = &domain.Account{
				Address: blockTx.To,

				// Since we are iterating in reverse in the blockchain, we are
				// starting at the latest block data and then iterating until
				// we reach a genesis; therefore, if this account is created then
				// this is their most recent transaction so therefore we want to
				// save the nonce.
				Nonce: blockData.Header.Nonce,
			}
		}
		acc.Balance += blockTx.Value

		if err := s.upsertAccountUseCase.Execute(acc.Address, acc.Balance, acc.Nonce); err != nil {
			s.logger.Error("Failed upserting account.",
				slog.Any("error", err))
			return err
		}

		s.logger.Debug("New `To` account balance via mining",
			slog.Any("account_address", acc.Address),
			slog.Any("balance", acc.Balance),
			slog.Any("tx_hash", blockTx.Hash),
		)
	}

	return nil
}
