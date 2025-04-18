package signature

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/mqdvi-dp/go-common/logger"
)

func ParseRsaPublicKeyFromPemStr(ctx context.Context, publicKeyValue string) (*rsa.PublicKey, error) {
	var err error

	block, _ := pem.Decode([]byte(publicKeyValue))
	if block == nil {
		logger.Log.Errorf(ctx, "failed to parse pem block containing the key")
		return nil, err
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		logger.Log.Errorf(ctx, "failed to ParsePKIXPublicKey")
		return nil, err
	}

	switch pub := publicKey.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break
	}
	return nil, nil
}
