package middleware

import (
	"context"
	"strings"

	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/config/database/rdc"
	"github.com/mqdvi-dp/go-common/types"
)

type middleware struct {
	validator           abstract.AuthenticationValidator
	authTypeCheckerFunc map[string]func(context.Context, string) (types.TokenClaim, error)
	redis               rdc.Rdc
}

// New initiate middleware
func New(authValidator abstract.AuthenticationValidator, rds rdc.Rdc) *middleware {
	m := &middleware{validator: authValidator, redis: rds}

	m.authTypeCheckerFunc = map[string]func(context.Context, string) (types.TokenClaim, error){
		basic: func(ctx context.Context, token string) (types.TokenClaim, error) {
			return types.TokenClaim{RoleKey: public}, m.Basic(ctx, token)
		},
		bearer: func(ctx context.Context, token string) (types.TokenClaim, error) {
			return m.Bearer(ctx, bearer, token)
		},
	}

	return m
}

// extractAuthType will get and validation type of authorization
func (m *middleware) extractAuthType(auth string) (tokenType, token string, err error) {
	auths := strings.Split(auth, " ")
	if len(auths) == 2 {
		tokenType = auths[0]
		token = auths[1]
		// check if an auth type is implemented
		if _, ok := m.authTypeCheckerFunc[tokenType]; ok {
			return
		}
	}

	err = ErrInvalidAuthorization
	return
}
