package wallets

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"golang.org/x/crypto/scrypt"
)

func encryptPrivateKey(privateKey, password string) (string, string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", "", err
	}

	// Derive a key from the password and salt
	key, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return "", "", err
	}

	// Encrypt the private key using AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", err
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", err
	}
	ciphertext := aesGCM.Seal(nil, nonce, []byte(privateKey), nil)

	// Return base64-encoded ciphertext and salt
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), base64.StdEncoding.EncodeToString(salt), nil
}

func decryptPrivateKey(encryptedKeyB64, saltB64, password string) (string, error) {
	// Decode the base64-encoded salt
	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return "", err
	}

	// Decode the base64-encoded encrypted key (nonce + ciphertext)
	encrypted, err := base64.StdEncoding.DecodeString(encryptedKeyB64)
	if err != nil {
		return "", err
	}

	// Derive the key from the password and salt using scrypt
	key, err := scrypt.Key([]byte(password), salt, 1<<15, 8, 1, 32)
	if err != nil {
		return "", err
	}

	// Prepare AES-GCM
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(encrypted) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := encrypted[:nonceSize], encrypted[nonceSize:]

	// Decrypt
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
