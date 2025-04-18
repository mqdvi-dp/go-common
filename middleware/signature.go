package middleware

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/zone"
)

// compareSignature compare the signature from request with generated
// func compareSignature(ctx context.Context, messageFromRequest, generatedSignature string) bool {
// 	logger.Log.Printf(ctx, "signature from request: %s", messageFromRequest)
// 	logger.Log.Printf(ctx, "signature generated from server: %s", generatedSignature)

// 	return generatedSignature == messageFromRequest
// }

// compareSignatureRSA compare the signature from request with algorithm RSA and SHA512
func compareSignatureRSA(ctx context.Context, publicKey *rsa.PublicKey, messageFromRequest string, signature [64]byte) bool {
	// decode the signature from request
	signatureRequest, err := base64.StdEncoding.DecodeString(messageFromRequest)
	if err != nil {
		tracer.SetError(ctx, fmt.Errorf("signature from request is not base64: %s", err))
		return false
	}

	tracer.Log(ctx, "signature from request", fmt.Sprintf("signature from request: %s", fmt.Sprintf("%x", sha512.Sum512(signatureRequest))))
	tracer.Log(ctx, "signature generated from server", fmt.Sprintf("signature generated from server: %s", fmt.Sprintf("%x", signature)))

	// verify the signature from request
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA512, signature[:], signatureRequest)
	if err != nil {
		tracer.SetError(ctx, fmt.Errorf("failed to verify the signature from request: %s", err))
		return false
	}

	return true
}

// signatureFormat will creates a signature format
func signatureFormat(ctx context.Context, c *gin.Context) (string, error) {
	var (
		body         string
		req          = c.Request
		httpMethod   = strings.ToUpper(req.Method)
		endpoint     = req.URL.Path
		channel      = req.Header.Get(constants.ApplicationChannel)
		reqTimestamp = req.Header.Get(constants.ApplicationTimestamp)
		deviceId     = req.Header.Get(constants.ApplicationDevice)
		secretKey    = c.GetString(constants.ClientSecret)
	)

	if secretKey == "" {
		secretKey = env.GetString("KLIKOO_SECRET_KEY")
	}

	// only for http.MethodPost and http.MethodPut
	if req.Method == http.MethodPost || req.Method == http.MethodPut {
		bodyPayload, err := io.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			return "", err
		}

		// hash the payload
		if len(bodyPayload) > 0 {
			body = fmt.Sprintf("%x", hashSHA256(bodyPayload))
		}

		// put again
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyPayload))
	}

	// parse into integer
	timestampInt, err := strconv.Atoi(reqTimestamp)
	if err != nil {
		return "", err
	}

	now := time.Now().In(zone.TzJakarta())
	// parse the timestamp from request header
	t := time.Unix(int64(timestampInt), 0).In(zone.TzJakarta())
	maxTs := env.GetDuration("HTTP_MAX_DURATION_TIMESTAMP", 5*time.Minute)
	now = now.Add(maxTs)
	tracer.Log(ctx, "check current_time from server", fmt.Sprintf("now: %v and t: %v", now, t))
	// validate the timestamp
	if now.Before(t) {
		return "", fmt.Errorf("timestamp is expired")
	}
	ts := fmt.Sprintf("%d", t.Unix())

	// create the message of signature
	message := fmt.Sprintf(
		"%s:%s:%s:%s:%s:%s:%s",
		httpMethod, endpoint,
		body, deviceId,
		channel, ts, secretKey,
	)

	tracer.Log(ctx, "digest", fmt.Sprintf("message formatted: %s", message))
	return message, nil
}

func hashSHA256(value []byte) [32]byte {
	return sha256.Sum256(value)
}

func hashSHA512(value []byte) [64]byte {
	return sha512.Sum512(value)
}
