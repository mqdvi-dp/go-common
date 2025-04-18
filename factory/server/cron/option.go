package cron

import (
	"github.com/mqdvi-dp/go-common/cronexpr"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/factory"
)

type option struct {
	debugMode     bool
	maxGoroutines int
	locker        cronexpr.Locker
}

type OptionFunc func(*option)

func getDefaultOption(service factory.ServiceFactory) option {
	opt := option{
		debugMode:     env.GetBool("DEBUG_MODE"),
		maxGoroutines: env.GetInt("CRON_MAX_GOROUTINES", 20),
	}

	opt.locker = cronexpr.NoopLocker{} // default
	if redisConn := service.GetDependencies().GetRedisDatabase(); redisConn != nil {
		opt.locker = cronexpr.NewRedisLocker(redisConn.Client())
	}

	return opt
}

func SetDebugMode(debugMode bool) OptionFunc {
	return func(o *option) {
		o.debugMode = debugMode
	}
}

func SetLocker(locker cronexpr.Locker) OptionFunc {
	return func(o *option) {
		o.locker = locker
	}
}
