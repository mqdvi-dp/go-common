package signature

import "crypto/sha256"

func hashSHA256(value []byte) [32]byte {
	return sha256.Sum256(value)
}
