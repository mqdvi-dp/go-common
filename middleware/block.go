package middleware

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/errs"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/redis/go-redis/v9"
)

const dot = "."

const prefixBlockedIp = "ip_blocked:%s"

func (m *middleware) getIp(ctx context.Context, originalIp string) (string, error) {
	var err error
	if m.redis == nil {
		err = fmt.Errorf("middleware: redis is nil")
		tracer.SetError(ctx, err)

		return "", errs.NewErrorWithCodeErr(err, errs.REDIS_CONNECTION_ERROR)
	}

	// splits the ip
	ips := strings.Split(originalIp, dot)
	if len(ips) != 4 {
		err = fmt.Errorf("ip is invalid")
		tracer.SetError(ctx, err)

		return "", errs.NewErrorWithCodeErr(err, errs.INVALID_IP)
	}

	originalIp = strings.Join(ips[:3], dot)
	return originalIp, nil
}

func (m *middleware) isIpBlocked(ctx context.Context, originalIp string) error {
	var err error
	originalIp, err = m.getIp(ctx, originalIp)
	if err != nil {
		return err
	}

	result, err := m.redis.Get(ctx, fmt.Sprintf(prefixBlockedIp, originalIp))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}

		tracer.SetError(ctx, err)
		return errs.NewErrorWithCodeErr(err, errs.REDIS_CONNECTION_ERROR)
	}

	val, err := strconv.Atoi(result)
	if err != nil {
		tracer.SetError(ctx, err)

		return errs.NewErrorWithCodeErr(err, errs.GENERAL_ERROR)
	}

	if val >= env.GetInt("TOTAL_IP_BLOCKED", 1) {
		err = fmt.Errorf("ip blocked")
		tracer.SetError(ctx, err)

		return errs.NewErrorWithCodeErr(err, errs.IP_BLOCKED)
	}

	return nil
}

func (m *middleware) blockedIp(ctx context.Context, originalIp string) error {
	var err error
	err = m.isIpBlocked(ctx, originalIp)
	if err != nil {
		return err
	}

	originalIp, err = m.getIp(ctx, originalIp)
	if err != nil {
		return err
	}

	duration := env.GetDuration("DURATION_IP_BLOCKED", time.Hour*730)
	err = m.redis.Set(ctx, fmt.Sprintf(prefixBlockedIp, originalIp), "1", duration)
	if err != nil {
		tracer.SetError(ctx, err)

		return errs.NewErrorWithCodeErr(err, errs.REDIS_INSERT_FAILED)
	}

	return nil
}
