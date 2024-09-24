package wallet

import (
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	// "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain"
)

const keystoreDirName = "keystore"

func GetKeystoreDirPath(dataDir string) string {
	return filepath.Join(dataDir, keystoreDirName)
}

func NewKeystoreAccount(dataDir, password string) (common.Address, string, error) {
	ks := keystore.NewKeyStore(GetKeystoreDirPath(dataDir), keystore.StandardScryptN, keystore.StandardScryptP)
	acc, err := ks.NewAccount(password)
	if err != nil {
		return common.Address{}, "", err
	}

	return acc.Address, acc.URL.Path, nil
}

// func SignTransactionWithKeystoreAccount(tx blockchain.Transaction, acc common.Address, pwd, keystoreDir string) (blockchain.SignedTransaction, error) {
// 	ks := keystore.NewKeyStore(keystoreDir, keystore.StandardScryptN, keystore.StandardScryptP)
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
// 	key, err := keystore.DecryptKey(ksAccountJson, pwd)
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

// func NewRandomKey() (*keystore.Key, error) {
// 	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	id := uuid.NewRandom()
// 	key := &keystore.Key{
// 		Id:         id,
// 		Address:    crypto.PubkeyToAddress(privateKeyECDSA.PublicKey),
// 		PrivateKey: privateKeyECDSA,
// 	}
//
// 	return key, nil
// }
