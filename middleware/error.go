package middleware

import "errors"

var (
	// ErrUnauthorized error when request does not have an Authorization header or user did not have access
	ErrUnauthorized error = errors.New("unauthorized")

	// ErrInvalidSession give an error when session is invalid
	// not found in redis or auth types are wrong
	ErrInvalidSession error = errors.New("session invalid")

	// ErrInvalidAuthorization is an error when format authorization is not match or invalid
	ErrInvalidAuthorization error = errors.New("invalid authorization")

	// ErrAccessDenied is error when user did not have access on resources
	ErrAccessDenied error = errors.New("access denied")

	// ErrTooManyRequest is error when users reach a maximum request to access api
	ErrTooManyRequest error = errors.New("too many request. please calm down")

	// ErrInvalidRateLimit is error when value rate limit is invalid
	ErrInvalidRateLimit error = errors.New("invalid value rate limit")

	// ErrInvalidTimestamp is error when timestamp header is not valid
	ErrInvalidTimestamp = errors.New("invalid timestamp format")

	// ErrInvalidTimestampExpired is errors when it's been more than 10minutes
	ErrInvalidTimestampExpired = errors.New("signature expired")

	// ErrInvalidSignature is errors when signature not matches with server
	ErrInvalidSignature = errors.New("invalid signature")

	// ErrIpNotWhitelisted is errors when request is provided skip-signature, but ip is not in whitelists
	ErrIpNotWhitelisted = errors.New("invalid ip")
)
