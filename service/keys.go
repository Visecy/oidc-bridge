package service

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"oidc-bridge/config"
)

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	keyData, err := os.ReadFile(config.AppConfig.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func LoadPublicKey() (*rsa.PublicKey, error) {
	keyData, err := os.ReadFile(config.AppConfig.PublicKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyData)
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PublicKey), nil
}
