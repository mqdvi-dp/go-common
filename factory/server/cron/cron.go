package cron

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/cronexpr"
	"github.com/mqdvi-dp/go-common/factory"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/types"
	"github.com/mqdvi-dp/go-common/zone"
)

const lockPattern = "%s:lock-cron-worker:%s"

type cronWorker struct {
	ctx                          context.Context
	cancelFunc                   func()
	tz                           *time.Location
	opt                          option
	service                      factory.ServiceFactory
	workers                      []reflect.SelectCase
	refreshWorkerNotif, shutdown chan struct{}
	semaphore                    []chan struct{}
	wg                           sync.WaitGroup
	activeJobs                   []*job
}

// NewWorker create new cron worker
func NewWorker(service factory.ServiceFactory, opts ...OptionFunc) factory.AppServerFactory {
	c := &cronWorker{
		service:            service,
		opt:                getDefaultOption(service),
		tz:                 zone.TzJakarta(),
		refreshWorkerNotif: make(chan struct{}),
		shutdown:           make(chan struct{}),
	}

	for _, opt := range opts {
		opt(&c.opt)
	}

	// add a shutdown channel to first index
	c.workers = append(
		c.workers, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(c.shutdown),
		},
	)

	// add a refresh worker channel to second index
	c.workers = append(
		c.workers, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(c.refreshWorkerNotif),
		},
	)
	// reset all cron workers before start from beginning
	c.opt.locker.Reset(fmt.Sprintf(lockPattern, c.service.Name(), "*"))

	if wh := service.WorkerHandler(constants.Scheduler); wh != nil {
		var hg types.WorkerHandlerGroup
		wh.Register(&hg)
		// add handler into job scheduler
		for _, handler := range hg.Handlers {
			if handler.Pattern == "" {
				logger.Log.Fatal("cron pattern not yet set. please set the pattern using, types.WorkerHandlerOptionPattern(cron.CreateSchedulerKey(param))")
			}
			funcName, interval := ParseSchedulerKey(handler.Pattern)

			j := job{
				handlerName: funcName,
				handler:     handler,
				interval:    interval,
			}

			if err := c.addJob(&j); err != nil {
				logger.Log.Fatalf("Cron Scheduler Worker: '%s' (interval: %s) %s", funcName, interval, err)
			}

			c.semaphore = append(c.semaphore, make(chan struct{}, c.opt.maxGoroutines))
			logger.Yellow(fmt.Sprintf(`⇨ [CRON-WORKER] (job name): "%s" (every): %-8s`, funcName, interval))
		}
	}

	fmt.Printf("\x1b[34;1m⇨ Cron worker running with %d jobs\x1b[0m\n\n", len(c.activeJobs))
	c.ctx, c.cancelFunc = context.WithCancel(context.Background())
	return c
}

func (c *cronWorker) Name() string {
	return string(constants.Scheduler)
}

func (c *cronWorker) Serve() {
	for _, j := range c.activeJobs {
		c.workers[j.workerIndex].Chan = reflect.ValueOf(j.ticker.C)
	}

	// running worker
	for {
		chosen, _, ok := reflect.Select(c.workers)
		if !ok {
			continue
		}

		// if a shutdown channel captured, break loop (no more job will run)
		if chosen == 0 {
			return
		}

		// notify for refresh worker
		if chosen == 1 {
			continue
		}

		chosen = chosen - 2
		j := c.activeJobs[chosen]
		c.registerNextInterval(j)

		if len(c.semaphore[j.workerIndex-2]) >= c.opt.maxGoroutines {
			continue
		}

		c.semaphore[j.workerIndex-2] <- struct{}{}
		c.wg.Add(1)
		go func(j *job) {
			defer func() {
				c.wg.Done()
				<-c.semaphore[j.workerIndex-2]
			}()

			if c.ctx.Err() != nil {
				logger.Red(fmt.Sprintf("cron_scheduler > context root err: %s", c.ctx.Err()))
				return
			}

			c.processJob(j)
		}(j)
	}
}

func (c *cronWorker) Shutdown(_ context.Context) {
	defer func() {
		fmt.Println("\x1b[33;1mStopping Cron Job Scheduler:\x1b[0m \x1b[32;1mSUCCESS\x1b[0m")
	}()

	if len(c.activeJobs) < 1 {
		return
	}

	c.stopAllJob()
	c.shutdown <- struct{}{}
	runningJob := 0
	for _, sem := range c.semaphore {
		runningJob += len(sem)
	}

	if runningJob != 0 {
		fmt.Printf("\x1b[34;1mCron Job Scheduler:\x1b[0m waiting %d job until done...\n", runningJob)
	}

	c.wg.Wait()
	c.cancelFunc()
	c.opt.locker.Reset(fmt.Sprintf(lockPattern, c.service.Name(), "*"))
}

func (c *cronWorker) processJob(j *job) {
	var err error
	start := time.Now().In(c.tz)
	ctx := c.ctx

	// lock for multiple worker (if running on multiple pods/instance)
	if c.opt.locker.IsLocked(c.getLockKey(j.handlerName)) {
		logger.Yellow(fmt.Sprintf("cron job > job %s is locked", j.handlerName))
		return
	}
	defer c.opt.locker.Unlock(c.getLockKey(j.handlerName))

	// implement logging
	ol := &logger.Logger{
		StartTime:   start.Format(time.RFC3339),
		RequestId:   uuid.NewString(),
		HandlerType: logger.Scheduler,
		Service:     c.service.Name(),
		Endpoint:    fmt.Sprintf("CRON %s", j.handlerName),
	}

	trace, ctx := tracer.StartTraceWithContext(ctx, fmt.Sprintf("CronScheduler:%s", j.handlerName))
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}

		sc := http.StatusOK
		if err != nil {
			trace.SetError(err)
			sc = http.StatusInternalServerError
			ol.ErrorMessage = fmt.Sprintf("%s", err)
		} else {
			ol.ResponseBody = "success"
		}

		trace.SetTag("trace_id", tracer.GetTraceId(ctx))
		ol.StatusCode = sc
		ol.ExecutionTime = time.Since(start).Seconds()
		trace.Finish()
		ol.Finalize(ctx)
	}()
	trace.SetTag("job_name", j.handlerName)

	// implement locking logging stdout
	var lock = logger.NewLocker(ctx)
	// set to context with logger.LogKey as a context key
	ctx = context.WithValue(ctx, logger.LogKey, lock)

	logger.YellowItalic(fmt.Sprintf("Cron Scheduler: executing task '%s' (interval: %s)", j.handlerName, j.interval))

	var ec types.EventContext
	ec.SetContext(ctx)
	ec.SetWorkerType(string(constants.Scheduler))
	ec.SetTopic(j.handlerName)
	ec.SetKey(j.handlerName)
	ec.SetHeader(map[string]interface{}{"interval": j.interval})

	if err := j.handler.HandlerFunc(&ec); err != nil {
		ec.SetError(err)
		trace.SetError(err)
	}
}

func (c *cronWorker) registerNextInterval(j *job) {
	if j.schedule != nil {
		j.ticker.Stop()
		j.ticker = time.NewTicker(j.schedule.NextInterval(time.Now().In(c.tz)))
		c.workers[j.workerIndex].Chan = reflect.ValueOf(j.ticker.C)
	} else if j.nextDuration != nil {
		j.ticker.Stop()
		j.ticker = time.NewTicker(*j.nextDuration)
		c.workers[j.workerIndex].Chan = reflect.ValueOf(j.ticker.C)
		j.nextDuration = nil
	}

	c.refreshWorker()
}

func (c *cronWorker) getLockKey(handlerName string) string {
	return fmt.Sprintf(lockPattern, c.service.Name(), handlerName)
}

func (c *cronWorker) refreshWorker() {
	go func() { c.refreshWorkerNotif <- struct{}{} }()
}

func (c *cronWorker) stopAllJob() {
	for _, j := range c.activeJobs {
		j.ticker.Stop()
	}
}

func (c *cronWorker) addJob(j *job) (err error) {
	if j.handler.HandlerFunc == nil {
		err = fmt.Errorf("handler func cannot empty")
		return
	}

	if reflect.ValueOf(j.handlerName).IsZero() {
		err = fmt.Errorf("handler name cannot empty")
		return
	}

	duration, nextDuration, err := cronexpr.ParseDuration(j.interval)
	if err != nil {
		j.schedule, err = cronexpr.Parse(j.interval)
		if err != nil {
			return
		}

		duration = j.schedule.NextInterval(time.Now().In(c.tz))
	}

	if nextDuration > 0 {
		j.nextDuration = &nextDuration
	}

	j.ticker = time.NewTicker(duration)
	j.workerIndex = len(c.workers)

	c.activeJobs = append(c.activeJobs, j)
	c.workers = append(
		c.workers, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(j.ticker.C),
		},
	)

	return nil
}
