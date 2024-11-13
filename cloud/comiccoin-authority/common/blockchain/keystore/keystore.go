package keystore

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

type KeystoreAdapter interface {
	CreateWallet(password string) (common.Address, []byte, error)
	OpenWallet(binaryData []byte, password string) (*keystore.Key, error)
}

type keystoreAdapterImpl struct{}

func NewAdapter() KeystoreAdapter {
	return &keystoreAdapterImpl{}
}

func (impl *keystoreAdapterImpl) CreateWallet(password string) (common.Address, []byte, error) {
	return createWalletWithTempFile(password)
}

func (impl *keystoreAdapterImpl) OpenWallet(binaryData []byte, password string) (*keystore.Key, error) {
	return decryptWalletFromBinary(binaryData, password)
}

func createWalletWithTempFile(password string) (common.Address, []byte, error) {
	tempDir, err := ioutil.TempDir("", "keystore")
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	ks := keystore.NewKeyStore(tempDir, keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.NewAccount(password)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to create account: %v", err)
	}

	keyJSON, err := ioutil.ReadFile(acc.URL.Path)
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed reading keystore JSON: %v", err)
	}

	return acc.Address, keyJSON, nil
}

func decryptWalletFromBinary(binaryData []byte, password string) (*keystore.Key, error) {
	tempFile, err := ioutil.TempFile("", "keystore-*")
	if err != nil {
		return nil, fmt.Errorf("failed creating temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(binaryData); err != nil {
		return nil, fmt.Errorf("failed writing to temp file: %v", err)
	}
	tempFile.Close()

	key, err := keystore.DecryptKey(binaryData, password)
	if err != nil {
		return nil, fmt.Errorf("failed decrypting key: %v", err)
	}
	return key, nil
}

// func SignTransactionWithKeystoreAccount(tx blockchain.Transaction, acc common.Address, pwd, keystoreDir string) (blockchain.SignedTransaction, error) {
// 	ks := kstore.NewKeyStore(keystoreDir, kstore.StandardScryptN, kstore.StandardScryptP)
// 	ksAccount, err := ks.Find(accounts.Account{Address: acc})
// 	if err != nil {
// 		return blockchain.SignedTransaction{}, err
// 	}
//
// 	ksAccountJson, err := ioutil.ReadFile(ksAccount.URL.Path)
// 	if err != nil {
// 		return blockchain.SignedTransaction{}, err
// 	}
//
// 	key, err := kstore.DecryptKey(ksAccountJson, pwd)
// 	if err != nil {
// 		return blockchain.SignedTransaction{}, err
// 	}
//
// 	signedTransaction, err := SignTransaction(tx, key.PrivateKey)
// 	if err != nil {
// 		return blockchain.SignedTransaction{}, err
// 	}
//
// 	return signedTransaction, nil
// }
//
// func SignTransaction(tx blockchain.Transaction, privKey *ecdsa.PrivateKey) (blockchain.SignedTransaction, error) {
// 	rawTransaction, err := tx.Encode()
// 	if err != nil {
// 		return blockchain.SignedTransaction{}, err
// 	}
//
// 	sig, err := Sign(rawTransaction, privKey)
// 	if err != nil {
// 		return blockchain.SignedTransaction{}, err
// 	}
//
// 	return blockchain.NewSignedTransaction(tx, sig), nil
// }
//
// func Sign(msg []byte, privKey *ecdsa.PrivateKey) (sig []byte, err error) {
// 	msgHash := sha256.Sum256(msg)
//
// 	return crypto.Sign(msgHash[:], privKey)
// }
//
// func Verify(msg, sig []byte) (*ecdsa.PublicKey, error) {
// 	msgHash := sha256.Sum256(msg)
//
// 	recoveredPubKey, err := crypto.SigToPub(msgHash[:], sig)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to verify message signature. %s", err.Error())
// 	}
//
// 	return recoveredPubKey, nil
// }

// func NewRandomKey() (*kstore.Key, error) {
// 	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	id := uuid.NewRandom()
// 	key := &kstore.Key{
// 		Id:         id,
// 		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
// 		PrivateKey: privateKeyECDSA,
// 	}
//
// 	return key, nil
// }
