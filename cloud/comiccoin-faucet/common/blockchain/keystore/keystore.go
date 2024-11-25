package keystore

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/securestring"
)

type KeystoreAdapter interface {
	CreateWallet(password *sstring.SecureString) (common.Address, []byte, error)
	OpenWallet(binaryData []byte, password *sstring.SecureString) (*keystore.Key, error)
}

type keystoreAdapterImpl struct{}

func NewAdapter() KeystoreAdapter {
	return &keystoreAdapterImpl{}
}

func (impl *keystoreAdapterImpl) CreateWallet(password *sstring.SecureString) (common.Address, []byte, error) {
	return CreateWalletInMemory(password)
}

func (impl *keystoreAdapterImpl) OpenWallet(binaryData []byte, password *sstring.SecureString) (*keystore.Key, error) {
	return OpenWalletFromMemory(binaryData, password)
}

func CreateWalletInMemory(password *sstring.SecureString) (common.Address, []byte, error) {
	log.Println("KeystoreAdapter: CreateWalletInMemory: Starting ...")

	// Generate an ECDSA private key.
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Printf("KeystoreAdapter: GenerateKey: Failed with error: %v\n", err)
		return common.Address{}, nil, fmt.Errorf("failed to generate key: %v", err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		log.Printf("KeystoreAdapter: uuid.NewRandom: Failed with error: %v\n", err)
		return common.Address{}, nil, fmt.Errorf("Could not create random uuid: %v", err)
	}

	// Wrap the key into a keystore.Key struct.
	key := &keystore.Key{
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
		Id:         id,
	}

	// Encrypt the key with reduced scrypt parameters.
	keyJSON, err := keystore.EncryptKey(key, password.String(), 2<<10, 8)
	if err != nil {
		log.Printf("KeystoreAdapter: EncryptKey: Failed with error: %v\n", err)
		return common.Address{}, nil, fmt.Errorf("failed to encrypt key: %v", err)
	}

	log.Println("KeystoreAdapter: CreateWalletInMemory: Success")
	return key.Address, keyJSON, nil
}

// OpenWalletFromMemory decrypts the given key JSON in memory using the provided password.
func OpenWalletFromMemory(binaryData []byte, password *sstring.SecureString) (*keystore.Key, error) {
	log.Println("KeystoreAdapter: OpenWalletFromMemory: Starting ...")

	// Decrypt the key JSON with the password.
	key, err := keystore.DecryptKey(binaryData, password.String())
	if err != nil {
		log.Printf("KeystoreAdapter: DecryptKey: Failed with error: %v\n", err)
		return nil, fmt.Errorf("failed to decrypt key: %v", err)
	}

	log.Println("KeystoreAdapter: OpenWalletFromMemory: Success")
	return key, nil
}
