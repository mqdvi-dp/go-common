package broker

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/monitoring"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/types"
	"github.com/mqdvi-dp/go-common/zone"
	"github.com/nsqio/go-nsq"
)

// nsqBroker configuration
type nsqBroker struct {
	producerHost string
	config       *nsq.Config
	client       *nsq.Producer
	publisher    abstract.Publisher
}

// NsqOptionFunc function type
type NsqOptionFunc func(*nsqBroker)

// SetNsqProducerHost set custom broker host
func SetNsqProducerHost(broker string) NsqOptionFunc {
	return func(nb *nsqBroker) {
		nb.producerHost = broker
	}
}

// SetNsqConfig set custom configuration nsq
func SetNsqConfig(cfg *nsq.Config) NsqOptionFunc {
	return func(nb *nsqBroker) {
		nb.config = cfg
	}
}

// SetNsqPublisher set customer publisher
func SetNsqPublisher(publisher abstract.Publisher) NsqOptionFunc {
	return func(nb *nsqBroker) {
		nb.publisher = publisher
	}
}

// defaultOptionNsq construct default nsq config
func defaultOptionNsq() *nsqBroker {
	return &nsqBroker{
		producerHost: env.GetString("NSQ_PRODUCER"),
		config:       nsq.NewConfig(),
	}
}

// NewNsqBroker setup nsq broker configuration for publisher, empty option param for default configuration
func NewNsqBroker(opts ...NsqOptionFunc) abstract.Broker {
	logger.PurpleItalic("Load nsq producer connection...")
	nb := defaultOptionNsq()

	for _, opt := range opts {
		opt(nb)
	}

	if nb.publisher == nil {
		if nb.producerHost == "" {
			panic("value env NSQ_PRODUCER not found")
		}

		producer, err := nsq.NewProducer(nb.producerHost, nb.config)
		if err != nil {
			panic(fmt.Errorf("%s. Brokers: %s", err, nb.producerHost))
		}

		nb.client = producer
		nb.publisher = NewNsqPublisher(nb.client)
		logger.GreenItalic("Load nsq producer connected!")
	}

	return nb
}

func (nb *nsqBroker) GetConfiguration() interface{} {
	return nb.config
}

func (nb *nsqBroker) GetPublisher() abstract.Publisher {
	return nb.publisher
}

func (nb *nsqBroker) GetName() constants.Worker {
	return constants.NSQ
}

func (nb *nsqBroker) Health() map[string]error {
	mErr := make(map[string]error)
	mErr[string(constants.NSQ)] = nil

	return mErr
}

func (nb *nsqBroker) Disconnect(ctx context.Context) error {
	nb.client.Stop()
	return nil
}

type nsqPublisher struct {
	producer *nsq.Producer
}

func NewNsqPublisher(producer *nsq.Producer) abstract.Publisher {
	np := &nsqPublisher{
		producer: producer,
	}

	return np
}

func (np *nsqPublisher) PublishMessage(ctx context.Context, arg *types.PublisherArgument) (err error) {
	if reflect.ValueOf(arg).IsZero() {
		err = fmt.Errorf("arguments cannot be empty")
		return
	}

	reqMessage := arg.Message
	if len(reqMessage) > env.GetInt("MAX_BODY_SIZE", 1500) {
		reqMessage = []byte(fmt.Sprintf("request body too long %d", len(reqMessage)))
	}

	now := time.Now().In(zone.TzJakarta())
	headers, _ := convert.InterfaceToString(arg.Header)
	trace := tracer.StartTrace(ctx, fmt.Sprintf("NSQ:PublishMessage:%s", arg.Topic))
	ol := logger.OutgoingLog{
		StartTime:     now.Format(constants.LayoutDateTime),
		TargetService: constants.NSQ.String(),
		URL:           fmt.Sprintf("topic: %s", arg.Topic),
		RequestBody:   string(reqMessage),
		RequestHeader: headers,
		StatusCode:    http.StatusOK,
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", err)
		}

		if err != nil {
			trace.SetError(err)
			ol.StatusCode = http.StatusInternalServerError
			ol.ResponseBody = fmt.Sprintf("%s", err)
		}
		since := time.Since(now)
		ol.ExecutionTime = since.Seconds()

		ol.Store(ctx)
		trace.Finish()
		monitoring.RecordPrometheus(ol.StatusCode, constants.NSQ.String(), ol.URL, since)
	}()

	trace.SetTag("topic", arg.Topic)
	trace.SetTag("message", arg.Message)

	err = np.producer.Publish(arg.Topic, arg.Message)
	if err != nil {
		trace.SetError(err)
		return
	}

	return nil
}

func (np *nsqPublisher) PublishMessages(ctx context.Context, req []*types.PublisherArgument) error {
	return nil
}
