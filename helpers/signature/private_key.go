package signature

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/mqdvi-dp/go-common/logger"
)

// getPrivateKeyPkcs8FromPem convert privateKey string into *rsa.PrivateKey
// @privateKeyValue should be base64 encoded
func getPrivateKeyPkcs8FromPem(ctx context.Context, privateKeyValue string) (*rsa.PrivateKey, error) {
	pkv, err := base64.StdEncoding.DecodeString(privateKeyValue)
	if err != nil {
		logger.Log.Errorf(ctx, "privateKeyValue is not base64 encoding. error while decode privateKeyValue: %s", err)
		return nil, err
	}

	block, _ := pem.Decode(pkv)
	if block == nil {
		logger.Log.Errorf(ctx, "failed to decode PEM PrivateKey")
		return nil, err
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		logger.Log.Errorf(ctx, "privateKeyValue should be pkcs#8. failed to parsing privateKey PKCS#8: %s", err)
		return nil, err
	}

	return privateKey.(*rsa.PrivateKey), nil
}
