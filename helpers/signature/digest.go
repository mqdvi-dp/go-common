package signature

import (
	"context"
	"fmt"
	"strings"

	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/logger"
)

const (
	formatDigestToken = "%s|%s"
	formatDigest      = "%s:%s:%s:%s"
)

func digestToken(ctx context.Context, clientKey, timestamp string) string {
	digest := fmt.Sprintf(formatDigestToken, clientKey, timestamp)
	logger.Log.Printf(ctx, "digest_token: %s", digest)

	return digest
}

func digest(ctx context.Context, httpMethod, url, timestamp string, body interface{}) (string, error) {
	bodyBytes, err := convert.InterfaceToBytes(body)
	if err != nil {
		logger.Log.Errorf(ctx, "failed when convert interface to bytes: %s", err)
		return "", err
	}
	bodyHashed := strings.ToLower(fmt.Sprintf("%x", hashSHA256(bodyBytes)))
	logger.Log.Printf(ctx, "body_hashed: %s", bodyHashed)

	digest := fmt.Sprintf(formatDigest, httpMethod, url, fmt.Sprintf("%x", hashSHA256(bodyBytes)), timestamp)
	logger.Log.Printf(ctx, "digest: %s", digest)
	return digest, nil
}
