// eth package used to interact with the Ethereum blockchain and provide an
// easy to use interface for performing our smart contract operations. The
// name of the smart contract is "Collectible Protection Services Submission
// Token".
package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"log/slog"
	"math"
	"math/big"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-cli/api"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	GetBalance(accountAddress string) (*big.Float, error)
	DeploySmartContract(smartContract, privateKey, baseURI string) (string, error)
	Mint(smartContract, privateKey, smartContractAddressHex, toAddressHex string) error
	GetTokenURI(smartContractAddress string, tokenId uint64) (string, error)
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
	// Special thanks to:  https://github.com/miguelmota/go-ethereum-hdwallet/blob/master/example/keys.go
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

func (e *ethBlockchain) GetBalance(accountAddress string) (*big.Float, error) {
	// Special thanks to: https://goethereumbook.org/account-balance/
	account := common.HexToAddress(accountAddress)
	balance, balanceErr := e.client.BalanceAt(context.Background(), account, nil)
	if balanceErr != nil {
		e.logger.Error("failed getting balance", slog.Any("error", balanceErr))
		return nil, balanceErr
	}

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	return ethValue, nil
}

func (e *ethBlockchain) DeploySmartContract(smartContract, privateKey, baseURI string) (string, error) {
	// Defensive code: Make sure the programmer is using only the specified
	// smart contract or error. This adapter can only support the following
	// smart contract.
	if smartContract != "Collectible Protection Service Submissions" {
		e.logger.Error("wront smart contract used")
		return "", fmt.Errorf("unsupported smart contract: %v", smartContract)
	}

	//
	// STEP 1
	// Get our public and wallet address derived from the private key.
	// Afterwords generate the necessary variables for generating a transaction
	// request to the blockchain.
	//

	privateKeyECDSA, convertErr := crypto.HexToECDSA(privateKey)
	if convertErr != nil {
		e.logger.Error("failed to convert from hex to ecdsa", slog.Any("error", convertErr))
		return "", fmt.Errorf("failed to convert private key value from hex to ecdsa: %v", convertErr)
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		e.logger.Error("failed to get public key from private key")
		return "", fmt.Errorf("cannot assert type: %v", "publicKey is not of type *ecdsa.PublicKey")
	}

	// Figure out who we're sending the ETH to, in this case it is our wallet.
	walletAccountAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//
	// STEP 2
	// Generate `nounce` for our smart contract deployment.
	//

	// We need to get the account nonce. Every transaction requires a nonce.
	// A nonce by definition is a number that is only used once.
	// Every new transaction from an account must have a nonce that the
	// previous nonce incremented by 1
	nonce, nonceErr := e.client.PendingNonceAt(context.Background(), walletAccountAddress)
	if nonceErr != nil {
		e.logger.Error("failed getting pending nonce", slog.Any("error", nonceErr))
		return "", fmt.Errorf("failed getting pending nonce: %v", nonceErr)
	}

	//
	// STEP 3
	// Generate our request to the blockchain for deploying our smart contract.
	//

	auth := bind.NewKeyedTransactor(privateKeyECDSA)
	auth.Nonce = big.NewInt(int64(nonce))

	// The next step is to set the amount of ETH that we'll be transferring.
	// However we must convert ether to wei since that's what the Ethereum
	// blockchain uses. Ether supports up to 18 decimal places so 1 ETH is 1
	// plus 18 zeros. Because we are deploying a smart contract, we can
	// safely set this value to zero.
	auth.Value = big.NewInt(0) // in wei

	// We need to set a gas prices for our transaction; however, gas prices are
	// always fluctuating based on market demand and what users are willing to
	// pay, so hardcoding a gas price is sometimes not ideal. The go-ethereum
	// client provides the `SuggestGasPrice` function for getting the average
	// gas price based on x number of previous blocks.
	gasPrice, suggestErr := e.client.SuggestGasPrice(context.Background())
	if suggestErr != nil {
		e.logger.Error("failed getting suggested gas price", slog.Any("error", suggestErr))
		return "", fmt.Errorf("failed getting suggested gas price: %v", suggestErr)
	}
	auth.GasPrice = gasPrice // The gas price must be set in `wei`.

	// The gas limit for a standard ETH transfer measured in `units`.
	auth.GasLimit = uint64(6000000) // in units

	//
	// STEP 3
	// Submit our request to the blockchain to deploy our smart contract.
	//

	deployedContractAddress, tx, instance, deployErr := DeployApi(auth, e.client, walletAccountAddress)
	if deployErr != nil {
		e.logger.Error("failed deploying api",
			slog.Any("gas_price", auth.GasPrice),
			slog.Any("gas_limit", auth.GasLimit),
			slog.Any("error", deployErr))
		return "", fmt.Errorf("failed deploying api: %v", deployErr)
	}

	// We do not need to use these.
	_ = tx
	_ = instance

	// Extract the new address that was provided to use
	smartContractAddressHex := deployedContractAddress.Hex()

	//
	// STEP 4
	// Generate a new request, submit a `SetBaseURI` request to the blockchain.
	//

	newNonce, newOnceErr := e.client.PendingNonceAt(context.Background(), walletAccountAddress)
	if newOnceErr != nil {
		e.logger.Error("failed pending nonce", slog.Any("error", deployErr))
		return "", fmt.Errorf("failed pending nonce at error: %v", newOnceErr)
	}
	auth.Nonce = big.NewInt(int64(newNonce))

	_, setErr := instance.SetBaseURI(auth, baseURI)
	if setErr != nil {
		e.logger.Error("failed deploying",
			slog.Any("error", deployErr),
			slog.Any("gas_price", auth.GasPrice),
			slog.Any("gas_limit", auth.GasLimit))
		return "", fmt.Errorf("failed deploying api: %v", setErr)
	}

	return smartContractAddressHex, nil
}

func (e *ethBlockchain) Mint(smartContract, privateKey, smartContractAddressHex, toAddressHex string) error {
	// Defensive code: Make sure the programmer is using only the specified
	// smart contract or error. This adapter can only support the following
	// smart contract.
	if smartContract != "Collectible Protection Service Submissions" {
		e.logger.Error("wront smart contract used")
		return fmt.Errorf("unsupported smart contract: %v", smartContract)
	}

	//
	// STEP 1
	// Get our public and wallet address derived from the private key.
	// Afterwords generate the necessary variables for generating a transaction
	// request to the blockchain.
	//

	privateKeyECDSA, convertErr := crypto.HexToECDSA(privateKey)
	if convertErr != nil {
		e.logger.Error("failed to convert from hex to ecdsa", slog.Any("error", convertErr))
		return fmt.Errorf("failed to convert private key value from hex to ecdsa: %v", convertErr)
	}

	publicKey := privateKeyECDSA.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		e.logger.Error("failed to get public key from private key")
		return fmt.Errorf("cannot assert type: %v", "publicKey is not of type *ecdsa.PublicKey")
	}

	// Figure out who we're sending the ETH to, in this case it is our wallet.
	walletAccountAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//
	// STEP 2
	// Generate `nounce` for our request to the blockchain.
	//

	// We need to get the account nonce. Every transaction requires a nonce.
	// A nonce by definition is a number that is only used once.
	// Every new transaction from an account must have a nonce that the
	// previous nonce incremented by 1
	nonce, nonceErr := e.client.PendingNonceAt(context.Background(), walletAccountAddress)
	if nonceErr != nil {
		e.logger.Error("failed getting pending nonce", slog.Any("error", nonceErr))
		return fmt.Errorf("failed getting pending nonce: %v", nonceErr)
	}

	//
	// STEP 3
	// Generate our actual request to the blockchain.
	//

	auth := bind.NewKeyedTransactor(privateKeyECDSA)
	auth.Nonce = big.NewInt(int64(nonce))

	// The next step is to set the amount of ETH that we'll be transferring.
	// However we must convert ether to wei since that's what the Ethereum
	// blockchain uses. Ether supports up to 18 decimal places so 1 ETH is 1
	// plus 18 zeros. Because we are making a request to our smart contract
	// that doesn't require any payment, so set to zero.
	auth.Value = big.NewInt(0) // in wei

	// We need to set a gas prices for our transaction; however, gas prices are
	// always fluctuating based on market demand and what users are willing to
	// pay, so hardcoding a gas price is sometimes not ideal. The go-ethereum
	// client provides the `SuggestGasPrice` function for getting the average
	// gas price based on x number of previous blocks.
	gasPrice, suggestErr := e.client.SuggestGasPrice(context.Background())
	if suggestErr != nil {
		e.logger.Error("failed getting suggested gas price", slog.Any("error", suggestErr))
		return fmt.Errorf("failed getting suggested gas price: %v", suggestErr)
	}
	auth.GasPrice = gasPrice // The gas price must be set in `wei`.

	// The gas limit for a standard ETH transfer measured in `units`.
	auth.GasLimit = uint64(6000000) // in units

	//
	// STEP 3
	// Submit our request to the blockchain.
	//

	smartContractAddress := common.HexToAddress(smartContractAddressHex)
	toAddress := common.HexToAddress(toAddressHex)
	instance, err := api.NewApi(smartContractAddress, e.client)
	if err != nil {
		log.Fatalf("An error occurred while establishing a connection with the smart contract : %v", err)
	}

	tx, err := instance.SafeMint(auth, toAddress)
	if err != nil {
		return fmt.Errorf("failed minting: %v", err)
	}

	// Ignore our transaction.
	_ = tx

	return nil
}

func (e *ethBlockchain) GetTokenURI(smartContractAddressHex string, tokenId uint64) (string, error) {
	smartContractAddress := common.HexToAddress(smartContractAddressHex)
	smartContractInstance, newApiErr := NewApi(smartContractAddress, e.client)
	if newApiErr != nil {
		e.logger.Error("An error occurred while establishing a connection with the smart contract",
			slog.Any("error", newApiErr))
		return "", newApiErr
	}

	// Create a new big.Int and set its value using SetUint64
	tokenIdBigInt := new(big.Int).SetUint64(tokenId)

	// Execute the `GetTokenURI` function from our smart contract which is a
	// standard function provided by following the ERC-721 standard.
	tokenUri, err := smartContractInstance.TokenURI(nil, tokenIdBigInt)
	if err != nil {
		return "", fmt.Errorf("failed getting token uri: %v", err)
	}

	return tokenUri, nil
}

func (e *ethBlockchain) Shutdown() {
	// Do nothing...
}
