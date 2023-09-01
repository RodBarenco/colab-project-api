package rsakeys

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
)

func DecryptWithPrivateKey(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(nil, privateKey, ciphertext)
}

func ReadPrivateKeyFromFile() (*rsa.PrivateKey, error) {
	privateKeyFile, err := os.Open("private_key.der")
	if err != nil {
		return nil, fmt.Errorf("error opening private key file: %v", err)
	}
	defer privateKeyFile.Close()

	privateKeyInfo, err := privateKeyFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting file info: %v", err)
	}

	privateKeyBytes := make([]byte, privateKeyInfo.Size())
	_, err = privateKeyFile.Read(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %v", err)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing private key: %v", err)
	}

	return privateKey, nil
}

// helper aes to large files

func DecryptAES(key []byte, payload []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating AES cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error creating GCM: %v", err)
	}

	nonceSize := gcm.NonceSize()
	if len(payload) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := payload[:nonceSize], payload[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting payload: %v", err)
	}

	return plaintext, nil
}
