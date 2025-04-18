package nsq

import (
	"context"
	"fmt"
	"log"

	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/factory"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/types"
	"github.com/nsqio/go-nsq"
)

type nsqWorker struct {
	opt             option
	service         factory.ServiceFactory
	cancelFunc      func()
	config          *nsq.Config
	engine          map[string]*nsq.Consumer
	consumerHandler map[string]*consumerHandler
}

// New create new nsq consumer
func New(service factory.ServiceFactory, opts ...OptionFunc) factory.AppServerFactory {
	worker := &nsqWorker{
		opt:             getDefaultOption(),
		service:         service,
		engine:          make(map[string]*nsq.Consumer),
		consumerHandler: make(map[string]*consumerHandler),
	}

	for _, opt := range opts {
		opt(&worker.opt)
	}

	if hdlr := service.WorkerHandler(constants.NSQ); hdlr != nil {
		var hg types.WorkerHandlerGroup
		hdlr.Register(&hg)

		for _, handler := range hg.Handlers {
			var channel = handler.Channel
			if channel == "" {
				channel = "default-channel"
			}
			worker.opt.channel = channel

			var consumerHandler consumerHandler
			consumerHandler.opt = &worker.opt
			consumerHandler.opt.serviceName = service.Name()
			consumerHandler.handlerFuncs = handler
			consumerHandler.topic = handler.Topic
			consumerHandler.channel = channel
			consumerHandler.isAutoAck = handler.AutoACK
			consumerHandler.ready = make(chan struct{})

			worker.consumerHandler[consumerHandler.topic] = &consumerHandler

			logger.Yellow(fmt.Sprintf(`[NSQ-CONSUMER] (topic): %-15s --> (channel): %-15s`, `"`+consumerHandler.topic+`"`, `"`+consumerHandler.channel+`"`))
		}
		logger.YellowBold(fmt.Sprintf("⇨ NSQ Consumer running with %d queue", len(worker.consumerHandler)))
	}

	cfg := nsq.NewConfig()
	cfg.MaxRequeueDelay = worker.opt.maxQueueDelay
	cfg.DefaultRequeueDelay = worker.opt.queueDelay
	cfg.MaxAttempts = uint16(worker.opt.maxAttempts)
	cfg.MaxInFlight = worker.opt.maxInflight

	worker.config = cfg
	return worker
}

func (n *nsqWorker) Serve() {
	ctx, cancel := context.WithCancel(context.Background())
	n.cancelFunc = cancel

	if len(n.consumerHandler) < 1 {
		return
	}

	for topic := range n.consumerHandler {
		consumer, err := nsq.NewConsumer(topic, n.consumerHandler[topic].channel, n.config)
		if err != nil {
			log.Fatalln(err)
		}
		consumer.AddHandler(n.consumerHandler[topic])
		if err := consumer.ConnectToNSQLookupds(n.consumerHandler[topic].opt.brokerHosts); err != nil {
			panic(fmt.Errorf("error start nsq %s", err))
		}

		if ctx.Err() != nil {
			return
		}

		n.engine[topic] = consumer
	}
}

func (n *nsqWorker) Shutdown(_ context.Context) {
	defer logger.RedBold("Stopping NSQ Broker")

	n.cancelFunc()
	for _, engine := range n.engine {
		engine.Stop()
		// stop all broker host
		for _, brokerHost := range n.opt.brokerHosts {
			engine.DisconnectFromNSQLookupd(brokerHost)
		}
	}
}

func (n *nsqWorker) Name() string {
	return constants.NSQ.String()
}
