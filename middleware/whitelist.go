package middleware

import (
	"context"
	"fmt"

	"github.com/mqdvi-dp/go-common/errs"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/tracer"
)

const prefixWhitelistDeviceId = "whitelist_device:%s"

func (m *middleware) checkWhitelist(ctx context.Context, hashedDeviceId, username string) (bool, error) {
	if m.redis == nil {
		err := fmt.Errorf("middleware: redis is nil")
		tracer.SetError(ctx, err)

		return false, errs.NewErrorWithCodeErr(err, errs.REDIS_CONNECTION_ERROR)
	}

	key := fmt.Sprintf(prefixWhitelistDeviceId, username)
	members, err := m.redis.DoSMembers(ctx, key)
	if err != nil {
		err = fmt.Errorf("middleware: redis DoSMembers error")
		tracer.SetError(ctx, err)

		return false, errs.NewErrorWithCodeErr(err, errs.REDIS_CONNECTION_ERROR)
	}

	if len(members) == 0 {
		logger.Log.Printf(ctx, "this user: %s has not been whitelisted: %s", username, hashedDeviceId)
		return true, nil
	}

	if !m.redis.DoSIsMember(ctx, key, hashedDeviceId) {
		logger.Log.Printf(ctx, "new device_id detected: %s", hashedDeviceId)
		return false, nil
	}

	logger.Log.Printf(ctx, "checkWhitelist bottom: %s", hashedDeviceId)

	return true, nil
}
