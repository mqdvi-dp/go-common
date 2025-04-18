package kafka

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/factory"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/types"
)

type kafkaWorker struct {
	option          option
	engine          sarama.ConsumerGroup
	service         factory.ServiceFactory
	consumerHandler *consumerHandler
	cancelFunc      func()
}

// New create new kafka consumer
func New(svc factory.ServiceFactory, opts ...OptionFunc) factory.AppServerFactory {
	kw := &kafkaWorker{
		service: svc,
		option:  getDefaultOption(),
	}
	kw.option.serviceName = svc.Name()
	for _, opt := range opts {
		opt(&kw.option)
	}

	// sarama configurations
	scfg := sarama.NewConfig()
	scfg.ClientID = kw.option.serviceName
	scfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{kw.option.balanceStrategy}
	scfg.Consumer.Retry.Backoff = kw.option.retryBackoff
	scfg.Consumer.Offsets.Initial = kw.option.offsetInitial
	scfg.Consumer.Offsets.Retry.Max = kw.option.maxRetry
	scfg.Consumer.Return.Errors = true

	// if sasl is enabled
	if kw.option.saslEnabled {
		scfg.Net.SASL.Enable = true
		scfg.Net.SASL.User = kw.option.saslUser
		scfg.Net.SASL.Password = kw.option.saslPassword
	}

	consumerEngine, err := sarama.NewConsumerGroup(kw.option.brokerHosts, kw.option.consumerGroup, scfg)
	if err != nil {
		logger.Log.Fatalf("failed to creating kafka consumer group client: %s", err)
	}
	kw.engine = consumerEngine

	var consumerHandler = new(consumerHandler)
	consumerHandler.handlerFuncs = make(map[string]types.WorkerHandler)
	if h := kw.service.WorkerHandler(constants.Kafka); h != nil {
		var hg types.WorkerHandlerGroup
		h.Register(&hg)

		for _, handler := range hg.Handlers {
			if _, ok := consumerHandler.handlerFuncs[handler.Topic]; !ok {
				consumerHandler.disableTrace = handler.DisableTrace
				consumerHandler.handlerFuncs[handler.Topic] = handler
				consumerHandler.topics = append(consumerHandler.topics, handler.Topic)

				logger.Yellow(fmt.Sprintf(`[KAFKA-CONSUMER] (topic): %-15s --> (group): %-15s`, `"`+handler.Topic+`"`, `"`+kw.option.consumerGroup+`"`))
			}
		}

		logger.YellowBold(fmt.Sprintf("⇨ KAFKA Consumer running with %d queue", len(consumerHandler.handlerFuncs)))
	}

	consumerHandler.ready = make(chan struct{}, kw.option.maxGoroutines)
	consumerHandler.option = &kw.option
	consumerHandler.messagePool = sync.Pool{
		New: func() any {
			return types.NewEventContext(bytes.NewBuffer(make([]byte, 0, 256)))
		},
	}

	// need producer?
	if kw.option.isNeedProducer {
		if svc.GetDependencies().GetBroker(constants.Kafka) == nil {
			logger.Log.Fatalf("missing dependency kafka")
		}

		consumerHandler.broker = svc.GetDependencies().GetBroker(constants.Kafka)
		consumerHandler.publisher = consumerHandler.broker.GetPublisher()
		consumerHandler.retrier = newRetrier()
	}
	kw.consumerHandler = consumerHandler
	return kw
}

func (kw *kafkaWorker) Serve() {
	ctx, cancel := context.WithCancel(context.Background())
	kw.cancelFunc = cancel

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := kw.engine.Consume(ctx, kw.consumerHandler.topics, kw.consumerHandler); err != nil {
				logger.Red("Error from kafka consumer: " + err.Error())

				if sErr, ok := err.(sarama.KError); ok {
					switch sErr {
					case sarama.ErrInvalidTopic:
						logger.Log.Fatal(sErr)
					}
				}
				// waiting
				time.Sleep(time.Second)
			}

			if ctx.Err() != nil {
				return
			}
			kw.consumerHandler.ready = make(chan struct{}, kw.option.maxGoroutines)
		}
	}()

	<-kw.consumerHandler.ready
	wg.Wait()
}

func (kw *kafkaWorker) Shutdown(context.Context) {
	defer logger.RedBold("Stopping KAFKA Broker")

	kw.cancelFunc()
	kw.engine.Close()
}

func (kw *kafkaWorker) Name() string {
	return constants.Kafka.String()
}
