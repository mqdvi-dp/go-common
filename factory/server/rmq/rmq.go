package rmq

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/factory"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/monitoring"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/types"
	"github.com/mqdvi-dp/go-common/zone"

	"github.com/streadway/amqp"
)

type rabbitMqWorker struct {
	ctx        context.Context
	cancelFunc func()
	opt        option
	tz         *time.Location
	ch         *amqp.Channel
	shutdown   chan struct{}
	isShutdown bool
	semaphore  []chan struct{}
	wg         sync.WaitGroup
	channels   []reflect.SelectCase
	handlers   map[string]types.WorkerHandler
}

// NewWorker create new rabbitmq consumer
func NewWorker(service factory.ServiceFactory, opts ...OptionFunc) factory.AppServerFactory {
	if service.GetDependencies().GetBroker(constants.RabbitMQ) == nil {
		logger.Log.Fatalf("missing dependencies rabbitmq")
	}

	worker := &rabbitMqWorker{
		opt: getDefaultOption(),
		tz:  zone.TzJakarta(),
	}
	for _, opt := range opts {
		opt(&worker.opt)
	}

	if reflect.ValueOf(worker.opt.serviceName).IsZero() {
		worker.opt.serviceName = service.Name()
	}

	worker.ctx, worker.cancelFunc = context.WithCancel(context.Background())
	worker.ch = service.GetDependencies().GetBroker(constants.RabbitMQ).GetConfiguration().(*amqp.Channel)
	worker.shutdown = make(chan struct{}, 1)
	worker.handlers = make(map[string]types.WorkerHandler)

	if h := service.WorkerHandler(constants.RabbitMQ); h != nil {
		var hg types.WorkerHandlerGroup
		h.Register(&hg)

		for _, handler := range hg.Handlers {
			worker.opt.exchangeName = handler.ExchangeName
			worker.opt.queue = handler.Queue
			worker.opt.consumerGroup = handler.Channel
			if worker.opt.consumerGroup == "" {
				worker.opt.consumerGroup = service.Name()
			}

			// when pattern is set, set the data into variable
			// this is a highest hierarchy
			if handler.Pattern != "" {
				worker.opt.exchangeName, worker.opt.queue, worker.opt.consumerGroup = ParseQueueKey(handler.Pattern)
			}

			logger.Yellow(fmt.Sprintf(`⇨ [RABBITMQ-CONSUMER] (queue): %-15s`, `"`+worker.opt.queue+`"`))

			queueChan, err := setupQueueConfig(worker.ch, worker.opt.consumerGroup, worker.opt.exchangeName, worker.opt.queue)
			if err != nil {
				panic(err)
			}

			worker.channels = append(
				worker.channels, reflect.SelectCase{
					Dir: reflect.SelectRecv, Chan: reflect.ValueOf(queueChan),
				},
			)
			worker.handlers[worker.opt.queue] = handler
			worker.semaphore = append(worker.semaphore, make(chan struct{}, 1))
		}
	}
	logger.YellowBold(fmt.Sprintf("\x1b[34;1m⇨ RabbitMQ consumer running with %d queue\n", len(worker.channels)))
	return worker
}

func (r *rabbitMqWorker) Name() string {
	return string(constants.RabbitMQ)
}

func (r *rabbitMqWorker) Shutdown(_ context.Context) {
	r.shutdown <- struct{}{}
	r.isShutdown = true
	var runningJob int
	for _, semp := range r.semaphore {
		runningJob += len(semp)
	}

	if runningJob != 0 {
		fmt.Printf("\x1b[34;1mRabbitMQ Worker:\x1b[0m waiting %d job until done...\x1b[0m\n", runningJob)
	}

	r.wg.Wait()
	_ = r.ch.Close()
	r.cancelFunc()
}

func (r *rabbitMqWorker) Serve() {
	for {
		select {
		case <-r.shutdown:
			return
		default:
		}

		chosen, value, ok := reflect.Select(r.channels)
		if !ok {
			continue
		}

		// execute handler
		if msg, ok := value.Interface().(amqp.Delivery); ok {
			r.semaphore[chosen] <- struct{}{}
			if r.isShutdown {
				return
			}

			r.wg.Add(1)
			go func(message amqp.Delivery, index int) {
				r.processMessage(message)
				r.wg.Done()
				<-r.semaphore[index]
			}(msg, chosen)
		}
	}
}

func (r *rabbitMqWorker) processMessage(message amqp.Delivery) {
	start := time.Now().In(r.tz)

	if r.ctx.Err() != nil {
		logger.Red(fmt.Sprintf("rabbitmq_consumer > ctx root err: %s", r.ctx.Err()))
		return
	}

	ctx := r.ctx
	selectedHandler := r.handlers[message.RoutingKey]

	header := make(map[string]interface{})
	for key, val := range message.Headers {
		vals, err := convert.InterfaceToString(val)
		if err != nil {
			continue
		}

		header[key] = vals
	}

	var err error
	trace, ctx := tracer.StartTraceWithContext(
		ctx,
		fmt.Sprintf("RabbitMqConsumer:%s", strings.ReplaceAll(convert.StringToTitle(message.RoutingKey), " ", "")),
	)

	// implement logging
	// init logger data
	ol := &logger.Logger{
		StartTime:     start.Format(time.RFC3339),
		RequestId:     uuid.NewString(),
		HandlerType:   logger.RabbitMQ,
		Service:       r.opt.serviceName,
		Endpoint:      fmt.Sprintf("Queue %s", r.opt.queue),
		RequestBody:   string(message.Body),
		RequestHeader: fmt.Sprintf("Exchange: %s | Routing Key: %s | Header: %v", message.Exchange, message.RoutingKey, header),
	}

	defer func() {
		if re := recover(); re != nil {
			err = fmt.Errorf("%s", re)
		}

		sc := http.StatusOK

		if err != nil {
			trace.SetError(err)
			_ = message.Reject(true)
			_ = message.Nack(true, true)

			sc = http.StatusInternalServerError
			ol.ErrorMessage = fmt.Sprintf("%s", err)
		} else {
			_ = message.Ack(true)

			ol.ResponseBody = "success"
		}

		since := time.Since(start)
		ol.StatusCode = sc
		ol.ExecutionTime = since.Seconds()

		// finish trace and logging
		trace.SetTag("trace_id", tracer.GetTraceId(ctx))
		trace.Finish()
		ol.Finalize(ctx)
		monitoring.RecordPrometheus(sc, constants.RabbitMQ.String(), ol.Endpoint, since)
	}()

	// implement locking logging stdout
	var lock = logger.NewLocker(ctx)
	// set to context with logger.LogKey as a context key
	ctx = context.WithValue(ctx, logger.LogKey, lock)

	trace.SetTag("exchange", message.Exchange)
	trace.SetTag("routing_key", message.RoutingKey)
	trace.SetTag("body", message.Body)
	trace.SetTag("header", header)

	logger.YellowItalic(fmt.Sprintf("\x1b[35;3mRabbitMQ Consumer: message consumed, topic = %s\x1b[0m", message.RoutingKey))

	ec := &types.EventContext{}
	ec.SetContext(ctx)
	ec.SetWorkerType(string(constants.RabbitMQ))
	ec.SetTopic(message.RoutingKey)
	ec.SetKey(message.Exchange)
	ec.SetHeader(header)
	_, _ = ec.Write(message.Body)

	if err = selectedHandler.HandlerFunc(ec); err != nil {
		ec.SetError(err)
		trace.SetError(err)
	}
}
