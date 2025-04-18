package helpers

import (
	"math/rand"
	"time"

	"github.com/mqdvi-dp/go-common/convert"
)

const defaultChartset = "0123456789"

// randomStringOption is a function that modifies a randomString struct
type randomStringOption func(*randomString)

// randomString is a struct that holds the options for the random string
type randomString struct {
	length            int
	charset           string
	withUnixTimestamp bool
}

// WithLength is a randomStringOption that sets the length of the random string
func WithLength(length int) randomStringOption {
	return func(r *randomString) {
		r.length = length
	}
}

// WithCharset is a randomStringOption that sets the charset of the random string
func WithCharset(charset string) randomStringOption {
	return func(r *randomString) {
		r.charset = charset
	}
}

// WithUnixTimestamp is a randomStringOption that sets the unix timestamp to the beginning of the random string
func WithUnixTimestamp() randomStringOption {
	return func(r *randomString) {
		r.withUnixTimestamp = true
	}
}

// AlphaNumeric is a randomStringOption that sets the charset to the alphanumeric characters
func AlphaNumeric() randomStringOption {
	return WithCharset("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// Alphabetical is a randomStringOption that sets the charset to the alphabetical characters
func Alphabetical() randomStringOption {
	return WithCharset("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// Numeric is a randomStringOption that sets the charset to the numeric characters
func Numeric() randomStringOption {
	return WithCharset("0123456789")
}

// defaultRandomString is a function that returns a pointer to a randomString struct
func defaultRandomString() *randomString {
	return &randomString{
		length:  6,
		charset: defaultChartset,
	}
}

// GenerateRandomNumber is a function that generates a random string
func GenerateRandomNumber(options ...randomStringOption) string {
	r := defaultRandomString()
	for _, option := range options {
		option(r)
	}

	// create a byte slice of the specified length
	randomBytes := make([]byte, r.length)

	// populate the byte slice with random characters from the charset
	for i := 0; i < r.length; i++ {
		randomBytes[i] = r.charset[rand.Intn(len(r.charset))]
	}

	// if the withUnixTimestamp option is not set, return the random string as is
	if !r.withUnixTimestamp {
		return string(randomBytes)
	}

	// if the withUnixTimestamp option is set, return the unix timestamp and the random string
	return convert.IntToString(time.Now().Unix()) + string(randomBytes)
}
