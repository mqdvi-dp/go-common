package helpers

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
)

func MD5(value string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(value)))
}

func SHA1(value string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(value)))
}

func SHA256(value string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(value)))
}
