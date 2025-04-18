package helpers

import "errors"

var (
	// ErrInvalidURL error message when validate the url and the value not start from https or http
	ErrInvalidURL = errors.New("invalid url. must start from 'https' or 'http'")

	// ErrURLSeperator error message when value don't have seperator between schema and domain
	ErrURLSeperator = errors.New("there is no seperator between schema and domain")

	// ErrDomain error message when host/domain is not exists
	ErrDomain = errors.New("host/domain not exists")
)
