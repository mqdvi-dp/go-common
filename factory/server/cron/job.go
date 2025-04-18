package cron

import (
	"time"

	"github.com/mqdvi-dp/go-common/cronexpr"
	"github.com/mqdvi-dp/go-common/types"
)

// job model
type job struct {
	handlerName  string
	interval     string
	handler      types.WorkerHandler
	workerIndex  int
	ticker       *time.Ticker
	nextDuration *time.Duration
	schedule     cronexpr.Schedule
}
