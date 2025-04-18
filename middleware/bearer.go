package middleware

import (
	"context"

	"github.com/mqdvi-dp/go-common/errs"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/types"
)

func (m *middleware) Bearer(ctx context.Context, tokenType, token string) (types.TokenClaim, error) {
	tc, err := m.validator.ValidateToken(ctx, tokenType, token)
	if err != nil {
		tracer.SetError(ctx, err)
		logger.Log.Errorf(ctx, "failed to validate bearer_token: %s", err)

		return types.TokenClaim{}, errs.TOKEN_EXPIRED
	}

	return tc, nil
}
