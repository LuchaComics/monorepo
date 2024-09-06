package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

type BlockchainAdapter interface {
	GetClient() *ethclient.Client
	Balance(address string) (*big.Int, error)
	GetAccountAuth(accountPrivateKey string) (*bind.TransactOpts, error)
	DeploySmartContract(auth *bind.TransactOpts, contractAddress string) (string, error)
	DeploySmartContractFromPrivateKey(accountPrivateKey string) (string, error)
	Mint(accountPrivateKey string, contractAddress string, toAddress string) error
	GetTokenURI(contractAddress string, tokenId *big.Int) (string, error)
}
