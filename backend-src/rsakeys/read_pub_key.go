package rsakeys

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"os"
)

func ReadPublicKeyFromFile(filename string) (string, error) {
	publicKeyFile, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer publicKeyFile.Close()

	publicKeyInfo, err := publicKeyFile.Stat()
	if err != nil {
		return "", err
	}

	publicKeyBytes := make([]byte, publicKeyInfo.Size())
	_, err = publicKeyFile.Read(publicKeyBytes)
	if err != nil {
		return "", err
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return "", err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("failed to cast public key to RSA public key")
	}

	publicKeyBytes = x509.MarshalPKCS1PublicKey(rsaPublicKey)
	publicKeyString := base64.StdEncoding.EncodeToString(publicKeyBytes)

	return publicKeyString, nil
}
