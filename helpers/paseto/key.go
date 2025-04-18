package paseto

import (
	"context"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/mqdvi-dp/go-common/logger"
)

func decodeKey(ctx context.Context, val string) ([]byte, error) {
	// decode the private key
	privateKeyBytes, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		logger.Log.Errorf(ctx, "failed to decode private key. make sure the private key is base64: %s", err)
		return nil, err
	}

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		err = fmt.Errorf("failed to decode PEM private key")
		logger.Log.Error(ctx, err)

		return nil, err
	}

	return block.Bytes, nil
}

// getPrivateKey private key for paseto login
func getPrivateKey(ctx context.Context, privateKeyValue string) (crypto.PrivateKey, error) {
	var err error
	// decode the value from env variable
	// make sure, the data is encoded
	privateKeyBytes, err := decodeKey(ctx, privateKeyValue)
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		logger.Log.Error(ctx, err)

		return nil, err
	}

	return privateKey, nil
}

func getPublicKey(ctx context.Context, publicKeyValue string) (crypto.PublicKey, error) {
	var err error
	// decode the value from env variable
	publicKeyBytes, err := decodeKey(ctx, publicKeyValue)
	if err != nil {
		logger.Log.Error(ctx, err)

		return nil, err
	}

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		logger.Log.Error(ctx, err)

		return nil, err
	}

	return publicKey, nil
}
