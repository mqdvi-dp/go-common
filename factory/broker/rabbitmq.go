package broker

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/types"
	"github.com/mqdvi-dp/go-common/zone"

	"github.com/streadway/amqp"
)

type rabbitMqBroker struct {
	host         string
	exchangeName string
	conn         *amqp.Connection
	ch           *amqp.Channel
	publisher    abstract.Publisher
}

func rabbitMqDefaultBrokerOption() *rabbitMqBroker {
	return &rabbitMqBroker{
		host:         env.GetString("RABBIT_MQ_HOST", "amqp://guest:guest@127.0.0.1:5672"),
		exchangeName: env.GetString("RABBIT_MQ_EXCHANGE_NAME", "default-exchange"),
	}
}

type RabbitMQOptionFunc func(*rabbitMqBroker)

// RabbitMQSetBrokerHost set host
func RabbitMQSetBrokerHost(host string) RabbitMQOptionFunc {
	return func(broker *rabbitMqBroker) {
		broker.host = host
	}
}

// RabbitMQSetExchangeName set broker exchange
func RabbitMQSetExchangeName(exchangeName string) RabbitMQOptionFunc {
	return func(broker *rabbitMqBroker) {
		broker.exchangeName = exchangeName
	}
}

// RabbitMQSetChannel set rabbit-mq channel
func RabbitMQSetChannel(ch *amqp.Channel) RabbitMQOptionFunc {
	return func(broker *rabbitMqBroker) {
		broker.ch = ch
	}
}

// RabbitMQSetPublisher set publisher rabbit-mq
func RabbitMQSetPublisher(publisher abstract.Publisher) RabbitMQOptionFunc {
	return func(broker *rabbitMqBroker) {
		broker.publisher = publisher
	}
}

// NewRabbitMQBroker setup rabbit-mq configuration for publisher and/or consumer, default connection from RABBIT_MQ_HOST environment
func NewRabbitMQBroker(opts ...RabbitMQOptionFunc) abstract.Broker {
	logger.PurpleItalic("Load rabbitmq client configuration...")
	var err error

	rmq := rabbitMqDefaultBrokerOption()
	for _, opt := range opts {
		opt(rmq)
	}

	rmq.conn, err = amqp.Dial(rmq.host)
	if err != nil {
		logger.Log.Fatalf("RabbitMQ: cannot connect to server broker: %s", err)
	}

	if rmq.ch == nil {
		rmq.ch, err = rmq.conn.Channel()
		if err != nil {
			logger.Log.Fatalf("RabbitMQ Channel: %s", err)
		}

		if err = rmq.ch.ExchangeDeclare(
			rmq.exchangeName,
			"topic",
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			logger.Log.Fatalf("RabbitMQ exchange declare delayed: %s", err)
		}

		if err = rmq.ch.Qos(2, 0, false); err != nil {
			logger.Log.Fatalf("RabbitMQ Qos: %s", err)
		}
	}

	if rmq.publisher == nil {
		rmq.publisher = NewRabbitMQPublisher(rmq.conn)
	}

	logger.GreenItalic("rabbitmq client connected!")
	return rmq
}

// GetConfiguration will return a rabbit-mq broker channel
func (r *rabbitMqBroker) GetConfiguration() interface{} {
	return r.ch
}

// GetPublisher return publisher rabbit-mq
func (r *rabbitMqBroker) GetPublisher() abstract.Publisher {
	return r.publisher
}

// GetName return worker name
func (r *rabbitMqBroker) GetName() constants.Worker {
	return constants.RabbitMQ
}

// Disconnect close connection rabbit-mq
func (r *rabbitMqBroker) Disconnect(ctx context.Context) error {
	logger.RedBold("rabbitmq: disconnecting...")
	defer fmt.Printf("\x1b[31;1mRabbitMQ Disconnecting:\x1b[0m \x1b[32;1mSUCCESS\x1b[0m\n")

	return r.conn.Close()
}

type rabbitMqPublisher struct {
	conn *amqp.Connection
}

// NewRabbitMQPublisher setup only rabbit-mq publisher with client connection
func NewRabbitMQPublisher(conn *amqp.Connection) abstract.Publisher {
	return &rabbitMqPublisher{conn: conn}
}

// PublishMessage publish a message to the topic with exchange name
func (r *rabbitMqPublisher) PublishMessage(ctx context.Context, req *types.PublisherArgument) (err error) {
	var ch *amqp.Channel

	trace := tracer.StartTrace(ctx, "rabbitmq:publish_message")
	defer func() {
		if re := recover(); re != nil {
			err = fmt.Errorf("%s", re)
		}

		if err != nil {
			trace.SetError(err)
		}

		if ch != nil {
			_ = ch.Close()
		}

		trace.Finish()
	}()

	ch, err = r.conn.Channel()
	if err != nil {
		return err
	}

	if reflect.ValueOf(req.ContentType).IsZero() {
		req.ContentType = constants.ApplicationJson
	}

	trace.SetTag("topic", req.Topic)
	trace.SetTag("key", req.Key)

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now().In(zone.TzJakarta()),
		ContentType:  req.ContentType,
		Headers:      amqp.Table(req.Header),
		Body:         req.Message,
	}

	trace.SetTag("headers", msg.Headers)
	trace.SetTag("body", msg.Body)
	trace.SetTag("timestamp", msg.Timestamp)

	if reflect.ValueOf(req.ExchangeName).IsZero() {
		req.ExchangeName = env.GetString("RABBIT_MQ_EXCHANGE_NAME", "default-exchange")
	}

	return ch.Publish(req.ExchangeName, req.Topic, false, false, msg)
}

func (r *rabbitMqPublisher) PublishMessages(ctx context.Context, req []*types.PublisherArgument) error {
	return nil
}
