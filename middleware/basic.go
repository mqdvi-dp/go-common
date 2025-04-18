package middleware

import (
	"context"

	"github.com/mqdvi-dp/go-common/errs"
	"github.com/mqdvi-dp/go-common/logger"
)

func (m *middleware) Basic(ctx context.Context, token string) error {
	logger.Log.Errorf(ctx, "Unauthorized basic token")
	return errs.UNAUTHORIZED
}
