package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-cli/api"
	c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-cli/config"
)

type ethBlockchain struct {
	client *ethclient.Client
}

func NewAdapter(cfg *c.Conf) BlockchainAdapter {
	log.Println("blockchain adapter initializing...")

	client, err := ethclient.Dial(cfg.EthServer.NodeURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("blockchain adapter initialized successfully")
	return &ethBlockchain{
		client: client,
	}
}

func (b *ethBlockchain) GetClient() *ethclient.Client {
	return b.client
}

func (b *ethBlockchain) Balance(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := b.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (b *ethBlockchain) GetAccountAuth(accountPrivateKey string) (*bind.TransactOpts, error) {
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

func (b *ethBlockchain) DeploySmartContract(auth *bind.TransactOpts, contractAddress string) (string, error) {
	contractAddressHex := common.HexToAddress(contractAddress)

	//deploying smart contract
	deployedContractAddress, tx, instance, err := api.DeployApi(auth, b.client, contractAddressHex) //api is redirected from api directory from our contract go file
	if err != nil {
		return "", fmt.Errorf("failed deploying api: %v", err)
	}
	_ = tx
	_ = instance
	return deployedContractAddress.Hex(), nil
}

func (b *ethBlockchain) DeploySmartContractFromPrivateKey(accountPrivateKey string) (string, error) {
	auth, err := b.GetAccountAuth(accountPrivateKey)
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

func (b *ethBlockchain) Mint(accountPrivateKey string, contractAddress string, toAddress string) error {
	auth, err := b.GetAccountAuth(accountPrivateKey)
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

func (b *ethBlockchain) GetTokenURI(contractAddress string, tokenId *big.Int) (string, error) {
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

func (b *ethBlockchain) Transfer(fromPrivateKeyStr string, toAddressHex string) error {
	privateKey, err := crypto.HexToECDSA(fromPrivateKeyStr)
	if err != nil {
		return fmt.Errorf("failed to convert value `%s` from hex to ecdsa: %v", privateKey, err)
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
	value := big.NewInt(1000000000000000000) // in wei (1 eth)
	auth.GasPrice = gasPrice
	gasLimit := uint64(21000)

	toAddress := common.HexToAddress(toAddressHex)

	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := b.client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = b.client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())

	return nil
}
