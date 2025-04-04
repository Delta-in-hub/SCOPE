package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// keyInHex must be a hex string, size in bytes must be 32
// Encrypt encrypts plaintext using AES-256 GCM.
// It returns the ciphertext (nonce + encrypted data) or an error.
func Encrypt(plaintext []byte, keyInHex string) ([]byte, error) {

	// Decode the hex key and check for errors
	key, err := hex.DecodeString(keyInHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex key: %w", err)
	}

	// Validate key size *after* successful decoding
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: must be 32 bytes (256 bits), got %d bytes", len(key))
	}

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		// This error typically occurs for invalid key sizes,
		// but we've already checked. Still good practice to handle.
		return nil, fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create a nonce with the standard size
	nonce := make([]byte, aesGCM.NonceSize())
	// Fill nonce with random data
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		// Error reading random data is serious
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data using aesGCM.Seal
	// Prepend the nonce to the ciphertext
	ciphertext := aesGCM.Seal(nonce /* prefix */, nonce /* nonce */, plaintext, nil /* additional data */)

	return ciphertext, nil // Return ciphertext and nil error on success
}

// Decrypt decrypts ciphertext using AES-256 GCM.
// It returns the plaintext or an error.
func Decrypt(ciphertextWithNonce []byte, keyInHex string) ([]byte, error) {
	key, err := hex.DecodeString(keyInHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex key: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: must be 32 bytes (256 bits), got %d bytes", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher block: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertextWithNonce) < nonceSize {
		return nil, fmt.Errorf("invalid ciphertext: too short to contain nonce")
	}

	// Split nonce and actual ciphertext
	nonce := ciphertextWithNonce[:nonceSize]
	ciphertext := ciphertextWithNonce[nonceSize:]

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil /* dst */, nonce, ciphertext, nil /* additional data */)
	if err != nil {
		// This error occurs if authentication fails (tampered data or wrong key/nonce)
		return nil, fmt.Errorf("failed to decrypt or authenticate data: %w", err)
	}

	return plaintext, nil
}
