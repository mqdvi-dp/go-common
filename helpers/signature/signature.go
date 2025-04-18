package signature

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"

	"github.com/mqdvi-dp/go-common/logger"
)

func Token(ctx context.Context, clientKey, privateKeyEncoded, timestamp string) (string, error) {
	// parse privateKeyEncoded into *rsa.PrivateKey
	privateKey, err := getPrivateKeyPkcs8FromPem(ctx, privateKeyEncoded)
	if err != nil {
		return "", err
	}

	// get digest
	digest := digestToken(ctx, clientKey, timestamp)
	// hash the digest
	sig := hashSHA256([]byte(digest))

	// sign digest with privateKey to generate value of signature
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, sig[:])
	if err != nil {
		logger.Log.Errorf(ctx, "failed when signPkcs1v15: %s", err)
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(sign)
	logger.Log.Printf(ctx, "signature encoded: %s", signature)
	return signature, nil
}

func Transaction(ctx context.Context, privateKeyEncoded, httpMethod, endpoint, timestamp string, body interface{}) (string, error) {
	// parse privateKeyEncoded into *rsa.PrivateKey
	privateKey, err := getPrivateKeyPkcs8FromPem(ctx, privateKeyEncoded)
	if err != nil {
		return "", err
	}

	// get digest
	digest, err := digest(ctx, httpMethod, endpoint, timestamp, body)
	if err != nil {
		return "", err
	}

	// hash the digest
	sig := hashSHA256([]byte(digest))

	// sign digest with privateKey to generate value of signature
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, sig[:])
	if err != nil {
		logger.Log.Errorf(ctx, "failed when signPkcs1v15: %s", err)
		return "", err
	}

	signature := base64.StdEncoding.EncodeToString(sign)
	logger.Log.Printf(ctx, "signature encoded: %s", signature)
	return signature, nil
}
