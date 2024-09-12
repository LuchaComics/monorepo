// eth package used to interact with the Ethereum blockchain and provide an
// easy to use interface for performing our smart contract operations. The
// name of the smart contract is "Collectible Protection Services Submission
// Token".
package ethereum

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/big"

	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

type EthereumWallet struct {
	AccountAddress string
	PrivateKey     string
	PublicKey      string
}

type EthereumAdapter interface {
	NewWalletFromMnemonic(mnemonic string) (*EthereumWallet, error)
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
func NewAdapter(logger *slog.Logger, nodeURL string) EthereumAdapter {
	logger.Debug("ethereum blockchain adapter initializing...")

	logger.Debug("ethereum blockchain adapter initialized")
	return &ethBlockchain{
		logger:  logger,
		nodeURL: nodeURL,
	}
}

func (e *ethBlockchain) NewWalletFromMnemonic(mnemonic string) (*EthereumWallet, error) {
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(wallet)

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := wallet.PrivateKeyHex(account)
	if err != nil {
		log.Fatal(err)
	}

	publicKey, _ := wallet.PublicKeyHex(account)
	if err != nil {
		log.Fatal(err)
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
