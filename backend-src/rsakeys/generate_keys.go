package rsakeys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"os"
	"time"
)

// KeyValidityFilePath é o caminho para o arquivo de validade da chave privada
const KeyValidityFilePath = "key_validity.json"

// GenerateKeys gera um novo par de chaves RSA e as salva nos arquivos DER
func genKeys(validUntil time.Time) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privateKeyFile, err := os.Create("private_key.der")
	if err != nil {
		return err
	}
	defer privateKeyFile.Close()

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	if _, err := privateKeyFile.Write(privateKeyBytes); err != nil {
		return err
	}

	// Create or overwrite the JSON file with private key validity
	keyValidity := KeyValidity{
		PrivateKeyValidUntil: validUntil.Format(time.RFC3339),
	}

	keyValidityFile, err := os.Create(KeyValidityFilePath)
	if err != nil {
		return err
	}
	defer keyValidityFile.Close()

	encoder := json.NewEncoder(keyValidityFile)
	if err := encoder.Encode(keyValidity); err != nil {
		return err
	}

	// Create the public key DER file
	publicKeyFile, err := os.Create("public_key.der")
	if err != nil {
		return err
	}
	defer publicKeyFile.Close()

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}

	if _, err := publicKeyFile.Write(publicKeyBytes); err != nil {
		return err
	}

	return nil
}
