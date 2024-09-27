package keystore

import (
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	kstore "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
)

const keystoreDirName = "keystore"

func getKeystoreDirPath(dataDir string) string {
	return filepath.Join(dataDir, keystoreDirName)
}

func NewKeystore(dataDir, password string) (common.Address, string, error) {
	ks := keystore.NewKeyStore(getKeystoreDirPath(dataDir), kstore.StandardScryptN, kstore.StandardScryptP)
	acc, err := ks.NewAccount(password)
	if err != nil {
		return common.Address{}, "", err
	}

	return acc.Address, acc.URL.Path, nil
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
