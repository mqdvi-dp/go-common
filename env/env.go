package env

import (
	"os"
	"strings"
)

// getEnv get value with specific key from env variable
func getEnv(key string) (string, bool) {
	val, ok := os.LookupEnv(strings.ToUpper(key))
	return val, ok
}
