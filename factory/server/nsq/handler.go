package nsq

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/monitoring"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/types"
	"github.com/nsqio/go-nsq"
)

type consumerHandler struct {
	opt          *option
	topic        string
	channel      string
	handlerFuncs types.WorkerHandler
	isAutoAck    bool
	ready        chan struct{}
}

func (c *consumerHandler) HandleMessage(msg *nsq.Message) (err error) {
	msg.Touch()
	ctx := context.Background()
	start := time.Now()
	ts := time.Unix(msg.Timestamp, 0)
	handler := c.handlerFuncs

	log.Printf("NSQ Consumer: message touch, timestamp = %s, topic = %s, channel = %s", ts, c.topic, c.channel)

	header := map[string]interface{}{
		"timestamp": ts.Format(time.RFC3339),
		"attempts":  strconv.Itoa(int(msg.Attempts)),
	}

	reqBody := msg.Body
	if len(reqBody) > env.GetInt("MAX_BODY_SIZE", 1500) {
		reqBody = []byte(fmt.Sprintf("request body too long %d", len(reqBody)))
	}

	// init logger data
	ol := &logger.Logger{
		StartTime:     start.Format(time.RFC3339),
		RequestId:     uuid.NewString(),
		HandlerType:   logger.NSQ,
		Service:       c.opt.serviceName,
		Endpoint:      fmt.Sprintf("topic: %s", c.topic),
		RequestBody:   string(reqBody),
		RequestHeader: fmt.Sprintf("Topic: %s | Channel: %s | Header: %v", c.topic, c.channel, header),
	}

	trace, ctx := tracer.StartTraceWithContext(ctx, "NSQConsumer")
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}

		sc := http.StatusOK

		ack := false
		if c.isAutoAck {
			ack = true
		}

		if err != nil {
			// send to queue again when has error
			msg.Requeue(c.opt.queueDelay)
			trace.SetError(err)
			sc = http.StatusInternalServerError
			ol.ErrorMessage = fmt.Sprintf("%s", err)
		} else {
			ol.ResponseBody = "success"

			// when we set the auto ack is true
			if ack {
				msg.Finish()
			}
		}
		since := time.Since(start)
		ol.StatusCode = sc
		ol.ExecutionTime = since.Seconds()

		trace.SetTag("trace_id", tracer.GetTraceId(ctx))
		trace.Finish()

		ol.Finalize(ctx)
		monitoring.RecordPrometheus(sc, constants.NSQ.String(), ol.Endpoint, since)
	}()

	trace.SetTag("message_id", msg.ID)
	trace.SetTag("topic", c.topic)
	trace.SetTag("channel", c.channel)
	trace.Log("header", header)
	trace.Log("message", msg.Body)

	var lock = logger.NewLocker(ctx)
	// set to context with logger.LogKey as a context key
	ctx = context.WithValue(ctx, logger.LogKey, lock)

	var ec types.EventContext
	ec.SetContext(ctx)
	ec.SetWorkerType(string(constants.NSQ))
	ec.SetTopic(c.topic)
	ec.SetHeader(header)
	ec.SetKey(c.channel)
	_, _ = ec.Write(msg.Body)

	if err = handler.HandlerFunc(&ec); err != nil {
		ec.SetError(err)
	}

	return
}
