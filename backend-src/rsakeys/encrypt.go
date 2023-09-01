package rsakeys

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
)

// EncryptAndEncode realiza a encriptação usando a chave pública em formato base64
// e retorna o resultado em formato hexadecimal.
func EncryptAndEncode(base64PublicKey string, data []byte) (string, error) {
	// Decodifica a chave pública base64 para bytes
	pubKeyBytes, err := base64.StdEncoding.DecodeString(base64PublicKey)
	if err != nil {
		return "", err
	}

	// Parse da chave pública
	pubKey, err := x509.ParsePKCS1PublicKey(pubKeyBytes)
	if err != nil {
		return "", err
	}

	// Encriptação
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey, data)
	if err != nil {
		return "", err
	}

	// Conversão para formato hexadecimal
	encryptedHex := hex.EncodeToString(encryptedData)
	return encryptedHex, nil
}

// GenerateAESKeyForLargeFile generates a random AES key suitable for encrypting large files.
func GenerateAESKeyForLargeFile() ([]byte, error) {
	keySize := 32 // AES-256
	key := make([]byte, keySize)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AES key: %v", err)
	}
	return key, nil
}

// use the ase to encrypt
func EncryptAES(key []byte, plaintext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a new GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	encryptedHex := hex.EncodeToString(ciphertext)
	return encryptedHex, nil
}
