package kafka

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/monitoring"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/types"
)

type consumerHandler struct {
	option       *option
	topics       []string
	handlerFuncs map[string]types.WorkerHandler
	disableTrace bool
	ready        chan struct{}
	messagePool  sync.Pool
	broker       abstract.Broker
	publisher    abstract.Publisher
	retrier      *retrier
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			go c.processMessage(session, message)
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *consumerHandler) processMessage(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) {
	handler, ok := c.handlerFuncs[message.Topic]
	if !ok {
		return
	}

	ctx := session.Context()
	start := time.Now()

	header := map[string]interface{}{
		"offset":    strconv.Itoa(int(message.Offset)),
		"partition": strconv.Itoa(int(message.Partition)),
		"timestamp": message.Timestamp.Format(time.RFC3339),
	}
	for _, val := range message.Headers {
		header[strings.ToLower(string(val.Key))] = string(val.Value)
	}

	reqBody := message.Value
	if len(reqBody) > env.GetInt("MAX_BODY_SIZE", 1500) {
		reqBody = []byte(fmt.Sprintf("request body too long %d", len(reqBody)))
	}

	var username string
	if h, ok := header["username"]; ok {
		if value, ok := h.(string); ok {
			username = value
		}
	}

	// init logger data
	hd, _ := convert.InterfaceToString(header)
	ol := &logger.Logger{
		StartTime:     start.Format(time.RFC3339),
		RequestId:     uuid.NewString(),
		HandlerType:   logger.Kafka,
		Service:       c.option.serviceName,
		Endpoint:      fmt.Sprintf("topic: %s", message.Topic),
		RequestBody:   string(reqBody),
		RequestHeader: hd,
		Username:      username,
	}

	var err error
	trace, ctx := tracer.StartTraceWithContext(ctx, "KafkaConsumer")
	defer func() {
		select {
		case <-ctx.Done():
			err = ctx.Err()
		default:
		}

		if r := recover(); r != nil {
			trace.SetTag("panic", true)
			err = fmt.Errorf("%v", r)
			fmt.Println(err)
			debug.PrintStack()
		}

		// whatever the response, we will mark the message
		session.MarkMessage(message, "")

		var sc = http.StatusOK
		if err != nil {
			sc = http.StatusInternalServerError
			ol.ErrorMessage = fmt.Sprintf("%s", err)
			// only if we implement the retry mechanism
			if c.retrier != nil {
				_ = c.retrier.retry(ctx, c, header, message, err)
			}
		} else {
			ol.ResponseBody = "success"
		}
		since := time.Since(start)
		ol.StatusCode = sc
		ol.ExecutionTime = since.Seconds()

		trace.SetTag("trace_id", tracer.GetTraceId(ctx))
		trace.Finish()
		monitoring.RecordPrometheus(sc, constants.Kafka.String(), ol.Endpoint, since)
		// when disable trace is false
		// means we will trace the logs
		if !c.disableTrace || err != nil {
			ol.Finalize(ctx)
		}
	}()

	trace.SetTag("brokers", c.option.brokerHosts)
	trace.SetTag("topic", message.Topic)
	trace.SetTag("key", message.Key)
	trace.SetTag("consumer_group", c.option.consumerGroup)
	trace.Log("header", header)
	trace.Log("message", message.Value)

	var lock = logger.NewLocker(ctx)
	// set to context with logger.LogKey as a context key
	ctx = context.WithValue(ctx, logger.LogKey, lock)
	logger.SetUsername(ctx, ol.Username) // set username to context
	log.Printf("\x1b[35;3mKafka Consumer: message consumed, timestamp = %v, topic = %s\x1b[0m", message.Timestamp, message.Topic)

	ec := c.messagePool.Get().(*types.EventContext)
	defer c.releaseMessagePool(ec)
	ec.SetContext(ctx)
	ec.SetWorkerType(constants.Kafka.String())
	ec.SetTopic(message.Topic)
	ec.SetHeader(header)
	ec.SetKey(string(message.Key))
	_, _ = ec.Write(message.Value)

	if err = handler.HandlerFunc(ec); err != nil {
		ec.SetError(err)
	}
	header = ec.Header()
}

func (c *consumerHandler) releaseMessagePool(ec *types.EventContext) {
	ec.Reset()
	c.messagePool.Put(ec)
}
