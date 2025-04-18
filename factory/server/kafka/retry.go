package kafka

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/types"
)

const (
	keyHeaderAttempt     = "attempt"
	keyHeaderErrorReason = "error_reason"
	retryWithUnderscored = "retry_"
)

type Dlq struct {
	Topic  string `json:"topic"`
	Key    string `json:"key"`
	Value  string `json:"value"`
	Header string `json:"header"`
}

func defaultMaxRetry() int               { return 1 }
func defaultRetryBackOff() time.Duration { return 2 * time.Second }

type retryableAction interface {
	execute(ctx context.Context, req *types.PublisherArgument) error
	getRetryBackOff(topic string) (int, time.Duration)
	publishToDlq(ctx context.Context, topic string, message *sarama.ConsumerMessage, headers map[string]interface{}) error
}

// retrier is struct responsibility for handling retry logic
type retrier struct {
}

// newRetrier create new instance from retrier with some configurations
func newRetrier() *retrier {
	return &retrier{}
}

// retry performs the retryable action and returns the result
func (r *retrier) retry(ctx context.Context, action retryableAction, headers map[string]interface{}, message *sarama.ConsumerMessage, errs error) error {
	var err error
	attempt := 1 // default

	// check header attempts
	if val, ok := headers[keyHeaderAttempt]; ok {
		attemptString, ok := val.(string)
		if ok {
			attempt, err = strconv.Atoi(attemptString)
			if err != nil {
				attempt = 1
			}
		}
	}
	// set headers attempt
	headers[keyHeaderAttempt] = fmt.Sprintf("%d", attempt+1)
	// set headers error_reason if exists
	if errs != nil {
		headers[keyHeaderErrorReason] = fmt.Sprintf("%s", errs)
	}
	// get retry backoff from handler
	maxRetry, retryBackOff := action.getRetryBackOff(message.Topic)
	// if attempt > maxRetry, then we need to send to dlq (dead letter queue)
	if attempt > maxRetry {
		_ = action.publishToDlq(ctx, message.Topic, message, headers)
		return nil
	}

	// set the exponential backoff
	delay := math.Pow(retryBackOff.Seconds(), float64(attempt))
	go func() {
		select {
		case <-time.After(time.Duration(delay) * time.Second):
			key := string(message.Key)
			if !strings.Contains(key, retryWithUnderscored) {
				key = fmt.Sprintf("%s%s", retryWithUnderscored, key)
			}

			_ = action.execute(ctx, &types.PublisherArgument{Topic: message.Topic, Message: message.Value, Header: headers, Key: key})
		case <-ctx.Done():
			_ = action.publishToDlq(context.Background(), message.Topic, message, headers)
		}
	}()

	return nil
}

// execute publishes a retryable message into the topic again
func (c *consumerHandler) execute(ctx context.Context, req *types.PublisherArgument) error {
	return c.publisher.PublishMessage(ctx, req)
}

// getRetryBackOff returns the retry backoff for a given topic
func (c *consumerHandler) getRetryBackOff(topic string) (int, time.Duration) {
	handler, ok := c.handlerFuncs[topic]
	if !ok {
		return defaultMaxRetry(), defaultRetryBackOff()
	}

	return handler.MaxRetry, handler.RetryBackoff
}

// publishToDlq publishes a message to the global dlq topic
func (c *consumerHandler) publishToDlq(ctx context.Context, topic string, message *sarama.ConsumerMessage, headers map[string]interface{}) error {
	header, _ := convert.InterfaceToString(headers)

	key := string(message.Key)
	if strings.Contains(key, retryWithUnderscored) {
		key = strings.ReplaceAll(key, retryWithUnderscored, "")
	}

	reqDlq, _ := convert.InterfaceToBytes(&Dlq{
		Topic:  message.Topic,
		Key:    key,
		Value:  string(message.Value),
		Header: header,
	})

	return c.execute(ctx, &types.PublisherArgument{Topic: env.GetString("TOPIC_DLQ", "klikoo-dlq"), Message: reqDlq, Key: fmt.Sprintf("topic: %s with key: %s", message.Topic, key)})
}
