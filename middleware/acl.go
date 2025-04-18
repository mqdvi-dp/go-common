package middleware

import (
	"context"
	"strings"

	"github.com/mqdvi-dp/go-common/logger"
)

func (m *middleware) checkACLPermission(ctx context.Context, roleId, resource string) error {
	// is resource exists?
	if resource == "" {
		logger.Log.Printf(ctx, "resource is empty")
		return nil
	}

	// when resource is public, skip to the next condition
	if strings.EqualFold(resource, public) {
		return nil
	}

	err := m.validator.CheckPermission(ctx, roleId, resource)
	if err != nil {
		return err
	}

	return nil
}
