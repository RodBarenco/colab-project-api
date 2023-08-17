package rsakeys

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

const keyValidityFilePath = "key_validity.json"

type KeyValidity struct {
	PrivateKeyValidUntil string `json:"private_key_valid_until"`
}

func EnsureKeysValid(validUntil time.Time) error {

	privateKeyExists := fileExists("private_key.der")
	publicKeyExists := fileExists("public_key.der")
	keyValidityExists := fileExists(keyValidityFilePath)

	if !privateKeyExists || !publicKeyExists {
		return genKeys(validUntil)
	}

	if !keyValidityExists {
		return errors.New("key validity information is missing")
	}

	// Read and parse key validity JSON
	keyValidity, err := readKeyValidity()
	if err != nil {
		return err
	}

	validity, err := time.Parse(time.RFC3339, keyValidity.PrivateKeyValidUntil)
	if err != nil {
		return err
	}

	if time.Now().After(validity) {
		return genKeys(validUntil)
	}

	return nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func readKeyValidity() (*KeyValidity, error) {
	file, err := os.Open(keyValidityFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var keyValidity KeyValidity
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&keyValidity); err != nil {
		return nil, err
	}

	return &keyValidity, nil
}
