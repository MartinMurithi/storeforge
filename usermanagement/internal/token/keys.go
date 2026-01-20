package token

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// LoadPrivateKey loads a PEM-encoded RSA private key from file
func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	const op = "keys.LoadPrivateKey"

	if path == "" {
		return nil, fmt.Errorf("%s: private key file path was not provided", op)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("%s: error reading private key from PEM file %w", op, err)
	}

	block, _ := pem.Decode(data)

	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("%s: invalid private key PEM", op)
	}

	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("%s: error parsing PKCS#8 private key %w", op, err)
	}

	rsaKey, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("%s: not an RSA private key", op)
	}

	return rsaKey, nil
}

// LoadPublicKey, loads a PEM encoded RSA public key from a file
func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	const op = "keys.LoadPublicKey"

	data, err := os.ReadFile(path)

	if err != nil {
		return nil, fmt.Errorf("%s: error reading public key from file %w", op, err)
	}

	block, _ := pem.Decode(data)

	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid public key PEM %w", err)
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, fmt.Errorf("error parsing public key %w", err)
	}

	pub, ok := pubInterface.(*rsa.PublicKey)

	if !ok {
		return nil, fmt.Errorf("no an RSA public key")
	}

	return pub, nil

}
