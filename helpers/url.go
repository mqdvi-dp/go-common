package helpers

import (
	"fmt"
	"net/url"
	"strings"
)

// URL checking value, start with https or http
func URL(val string, schemes ...string) error {
	u, err := url.Parse(val)
	if err != nil {
		return err
	}

	if u.Scheme == "" {
		return ErrInvalidURL
	}

	if u.Host == "" {
		return ErrDomain
	}

	// valid the schema if exists on params
	if len(schemes) > 0 {
		for _, scheme := range schemes {
			if scheme == u.Scheme {
				return nil
			}
		}

		return fmt.Errorf("invalid url scheme. url scheme should be %s", strings.Join(schemes, " or "))
	}

	return nil
}
