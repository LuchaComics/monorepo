package nftsmartcontract

import (
	"math/big"
)

type EthereumNFTSmartContractAdapter interface {
	DeployFromPrivateKey(accountPrivateKey string) (string, error)
	Mint(accountPrivateKey string, contractAddress string, toAddress string) error
	GetBalance(address string) (*big.Int, error)
	GetTokenURI(contractAddress string, tokenId *big.Int) (string, error)
}
