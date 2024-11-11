package cryptowrapper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/scrypt"
)

const (
	saltSize  = 16
	nonceSize = 12
	keySize   = 32 // AES-256 requires a 32-byte key
)

// SymmetricKeyEncryption encrypts plaintext with the given password and returns base64-encoded ciphertext.
func SymmetricKeyEncryption(plaintext, password string) (string, error) {
	// Generate a random salt.
	salt := make([]byte, saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Derive a key from the password and salt using scrypt.
	key, err := deriveKey(password, salt)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block with the derived key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Generate a random nonce for GCM mode.
	nonce := make([]byte, nonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Use GCM mode for encryption.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Encrypt the plaintext.
	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)

	// Combine the salt, nonce, and ciphertext for storage.
	result := append(salt, nonce...)
	result = append(result, ciphertext...)

	// Return the result as a base64-encoded string.
	return base64.StdEncoding.EncodeToString(result), nil
}

// SymmetricKeyDecryption decrypts the base64-encoded ciphertext using the given password.
func SymmetricKeyDecryption(encodedCiphertext, password string) (string, error) {
	// Decode the base64-encoded input.
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	// Extract the salt, nonce, and ciphertext parts.
	salt := ciphertext[:saltSize]
	nonce := ciphertext[saltSize : saltSize+nonceSize]
	actualCiphertext := ciphertext[saltSize+nonceSize:]

	// Derive the key from the password and salt using scrypt.
	key, err := deriveKey(password, salt)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher block with the derived key.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher block: %w", err)
	}

	// Use GCM mode for decryption.
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Decrypt the ciphertext.
	plaintext, err := aesGCM.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// deriveKey derives a key from a password and salt using the scrypt key derivation function.
func deriveKey(password string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(password), salt, 1<<14, 8, 1, keySize)
}
