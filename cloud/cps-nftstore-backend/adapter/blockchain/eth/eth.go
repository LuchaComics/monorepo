// eth package used to interact with the Ethereum blockchain and provide an
// easy to use interface for performing our smart contract operations. The
// name of the smart contract is "Collectible Protection Services Submission
// Token".
package ethereum

import (
	"context"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

type EthereumWallet struct {
	AccountAddress string
	PrivateKey     string
	PublicKey      string
}

type EthereumAdapter interface {
	ConnectToNodeAtURL(nodeURL string) error
	NewWalletFromMnemonic(mnemonic string) (*EthereumWallet, error)
	GetOwnersBalance(context context.Context) (*big.Int, error)
	Mint(toAddress string) error
	GetTokenURI(tokenId *big.Int) (string, error)
	Shutdown()
}

type ethBlockchain struct {
	logger  *slog.Logger
	nodeURL string
	client  *ethclient.Client
}

// NewAdapter function connects to an Ethereum node and provides an interface
// for our application to use to make smart contract interactions. The
// configuration variables required are:
//
// 1. CPS_NFTSTORE_BACKEND_ETH_NODE_URL: This
func NewAdapter(logger *slog.Logger) EthereumAdapter {
	logger.Debug("ethereum blockchain adapter initializing...")

	logger.Debug("ethereum blockchain adapter initialized")
	return &ethBlockchain{
		logger: logger,
	}
}

func (e *ethBlockchain) ConnectToNodeAtURL(nodeURL string) error {
	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		return err
	}
	if client != nil {
		e.client = client
		e.nodeURL = nodeURL
	}
	return nil
}

func (e *ethBlockchain) NewWalletFromMnemonic(mnemonic string) (*EthereumWallet, error) {
	wallet, newErr := hdwallet.NewFromMnemonic(mnemonic)
	if newErr != nil {
		e.logger.Error("failed creating new wallet from mnemonic", slog.Any("error", newErr))
		return nil, newErr
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, depriveErr := wallet.Derive(path, false)
	if depriveErr != nil {
		e.logger.Error("failed depriving", slog.Any("error", depriveErr))
		return nil, depriveErr
	}

	privateKey, getErr := wallet.PrivateKeyHex(account)
	if getErr != nil {
		e.logger.Error("failed getting private key hex", slog.Any("error", getErr))
		return nil, getErr
	}

	publicKey, getErr := wallet.PublicKeyHex(account)
	if getErr != nil {
		e.logger.Error("failed getting public key hex", slog.Any("error", getErr))
		return nil, getErr
	}

	return &EthereumWallet{
		AccountAddress: account.Address.Hex(),
		PrivateKey:     privateKey,
		PublicKey:      publicKey,
	}, nil
}

func (e *ethBlockchain) GetOwnersBalance(context context.Context) (*big.Int, error) {
	return nil, nil //TODO
}

func (e *ethBlockchain) Mint(toAddress string) error {
	return nil //TODO
}

func (e *ethBlockchain) GetTokenURI(tokenId *big.Int) (string, error) {
	return "", nil //TODO
}

func (e *ethBlockchain) Shutdown() {
	// Do nothing...
}
