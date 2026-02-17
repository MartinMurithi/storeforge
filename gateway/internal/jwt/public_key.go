package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

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