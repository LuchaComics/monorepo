// eth package used to interact with the Ethereum blockchain and provide an
// easy to use interface for performing our smart contract operations. The
// name of the smart contract is "Collectible Protection Services Submission
// Token".
package ethereum

import (
	"context"
	"log/slog"
	"math/big"
)

type EthereumBlockchainAdapter interface {
	GetOwnersBalance(context context.Context) (*big.Int, error)
	Mint(toAddress string) error
	GetTokenURI(tokenId *big.Int) (string, error)
	Shutdown()
}

type ethBlockchain struct {
	logger               *slog.Logger
	nodeURL              string
	smartContractAddress string
	ownerAddress         string
	ownerPrivateKey      string
}

// NewAdapter function connects to an Ethereum node and provides an interface
// for our application to use to make smart contract interactions. The
// configuration variables required are:
//
// 1. CPS_NFTSTORE_BACKEND_ETH_NODE_URL: This
func NewAdapter(cfg *c.Conf, logger *slog.Logger) EthereumBlockchainAdapter {
	logger.Debug("ethereum blockchain adapter initializing...")

	logger.Debug("ethereum blockchain adapter initialized")
	return &ethBlockchain{
		logger:               logger,
		nodeURL:              c.EthereumBlockchain.NodeURL,
		smartContractAddress: c.EthereumBlockchain.SmartContractAddress,
		ownerAddress:         c.EthereumBlockchain.OwnerAddress,
		ownerPrivateKey:      c.EthereumBlockchain.OwnerPrivateKey,
	}
}

func (e *ethBlockchain) GetOwnersBalance(context context.Context) (*big.Int, error) {
	return nil, nil //TODO
}

func (e *ethBlockchain) Mint(toAddress string) error {
	return nil, nil //TODO
}

func (e *ethBlockchain) GetTokenURI(tokenId *big.Int) (string, error) {
	return "", nil //TODO
}

func (e *ethBlockchain) Shutdown() {
	// Do nothing...
}
