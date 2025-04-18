package middleware

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/tracer"
)

// getPublicKey parse the public key from environment into struct rsa.PublicKey
func getPublicKey(ctx context.Context, publicKeyEncoded string) (*rsa.PublicKey, error) {
	var err error
	if publicKeyEncoded == "" {
		// get default public key from environment
		publicKeyEncoded = env.GetString("KLIKOO_PUBLIC_KEY")
		if publicKeyEncoded == "" {
			err = fmt.Errorf("public key is not found in context and environment")
			tracer.SetError(ctx, err)

			return nil, err
		}
	}

	// decode the public key
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyEncoded)
	if err != nil {
		tracer.SetError(ctx, fmt.Errorf("failed to decode public key. make sure the public key is base64: %s", err))
		return nil, err
	}

	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		err = fmt.Errorf("failed to decode PEM public key")
		tracer.SetError(ctx, err)

		return nil, err
	}

	// parse from bytes into public.Key instance
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		tracer.SetError(ctx, fmt.Errorf("failed to parse publicKey bytes. %s", err))
		return nil, err
	}

	// return the value
	return publicKey.(*rsa.PublicKey), nil
}
