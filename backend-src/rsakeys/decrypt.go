package rsakeys

import (
	"crypto/rsa"
	"crypto/x509"
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
