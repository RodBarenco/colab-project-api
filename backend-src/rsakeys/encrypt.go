package rsakeys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
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
