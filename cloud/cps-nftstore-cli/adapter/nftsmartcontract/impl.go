package nftsmartcontract

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-cli/api"
	c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-cli/config"
)

type nftSmartContractAdapter struct {
	client *ethclient.Client
}

func NewAdapter(cfg *c.Conf) EthereumNFTSmartContractAdapter {
	log.Println("ethereum nft smart contract adapter initializing...")

	client, err := ethclient.Dial(cfg.EthServer.NodeURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("ethereum nft smart contract adapter initialized successfully")
	return &nftSmartContractAdapter{
		client: client,
	}
}

func (b *nftSmartContractAdapter) DeployFromPrivateKey(accountPrivateKey string) (string, error) {
	auth, err := b.getAccountAuth(accountPrivateKey)
	if err != nil {
		return "", err
	}

	privateKey, err := crypto.HexToECDSA(accountPrivateKey[2:])
	if err != nil {
		return "", fmt.Errorf("failed to convert value `%s` from hex to ecdsa: %v", accountPrivateKey[2:], err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("%s", "invalid key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//deploying smart contract
	deployedContractAddress, tx, instance, err := api.DeployApi(auth, b.client, fromAddress) //api is redirected from api directory from our contract go file
	if err != nil {
		return "", fmt.Errorf("failed deploying api: %v", err)
	}
	_ = tx
	_ = instance
	return deployedContractAddress.Hex(), nil
}

func (b *nftSmartContractAdapter) Mint(accountPrivateKey string, contractAddress string, toAddress string) error {
	auth, err := b.getAccountAuth(accountPrivateKey)
	if err != nil {
		return err
	}

	// privateKey, err := crypto.HexToECDSA(accountPrivateKey[2:])
	// if err != nil {
	// 	return "", fmt.Errorf("failed to convert value `%s` from hex to ecdsa: %v", accountPrivateKey[2:], err)
	// }

	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	return "", fmt.Errorf("%s", "invalid key")
	// }

	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	contractAddressCommon := common.HexToAddress(contractAddress)
	contractInstance, err := api.NewApi(contractAddressCommon, b.client)
	if err != nil {
		log.Fatalf("An error occurred while establishing a connection with the smart contract : %v", err)
	}

	to := common.HexToAddress(toAddress)
	tx, err := contractInstance.SafeMint(auth, to)
	if err != nil {
		return fmt.Errorf("failed minting: %v", err)
	}
	// fmt.Println("tx:", tx)
	// fmt.Println("tx:to:", tx.To())
	// fmt.Printf("tx sent: %s", tx.Hash().Hex()) // tx sent: 0x8d490e535678e9a24360e955d75b27ad307bdfb97a1dca51d0f3035dcee3e870
	_ = tx

	return nil
}

func (b *nftSmartContractAdapter) GetTokenURI(contractAddress string, tokenId *big.Int) (string, error) {
	contractAddressCommon := common.HexToAddress(contractAddress)
	contractInstance, err := api.NewApi(contractAddressCommon, b.client)
	if err != nil {
		log.Fatalf("An error occurred while establishing a connection with the smart contract : %v", err)
	}

	res, err := contractInstance.TokenURI(nil, tokenId)
	if err != nil {
		return "", fmt.Errorf("failed getting token uri: %v", err)
	}

	return res, nil
}

func (b *nftSmartContractAdapter) GetBalance(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := b.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (b *nftSmartContractAdapter) getAccountAuth(accountPrivateKey string) (*bind.TransactOpts, error) {
	// NOTE: We need to remove the `0x` from our private key and get `[2:]` slice.

	privateKey, err := crypto.HexToECDSA(accountPrivateKey[2:])
	if err != nil {
		return nil, fmt.Errorf("failed to convert value `%s` from hex to ecdsa: %v", accountPrivateKey[2:], err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := b.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := b.client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasPrice = gasPrice

	return auth, nil
}
