package broker

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/monitoring"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/types"
	"github.com/mqdvi-dp/go-common/zone"
)

type KafkaBrokerConfigFunc func(*kafkaBrokerConfig)

type kafkaBrokerConfig struct {
	brokerHosts           []string
	kafkaVersion          string
	saslEnabled           bool
	saslUser              string
	saslPassword          string
	clientId              string
	producerMaxRetry      int
	producerBackoffRetry  time.Duration
	producerAsync         bool
	producerAck           sarama.RequiredAcks
	producerReturnSuccess bool
	publisher             abstract.Publisher
}

// SetKafkaBrokerHosts set kafka hosts
func SetKafkaBrokerHosts(hosts []string) KafkaBrokerConfigFunc {
	return func(kbc *kafkaBrokerConfig) {
		kbc.brokerHosts = hosts
	}
}

// SetKafkaVersion set kafka versions
func SetKafkaVersion(version string) KafkaBrokerConfigFunc {
	return func(kbc *kafkaBrokerConfig) {
		kbc.kafkaVersion = version
	}
}

// SetKafkaClientId sets kafka client id
func SetKafkaClientId(clientId string) KafkaBrokerConfigFunc {
	return func(kbc *kafkaBrokerConfig) {
		kbc.clientId = clientId
	}
}

// SetKafkaProducerMaxRetry sets kafka-producer retry publish message
func SetKafkaProducerMaxRetry(maxRetry int) KafkaBrokerConfigFunc {
	return func(kbc *kafkaBrokerConfig) {
		kbc.producerMaxRetry = maxRetry
	}
}

// SetKafkaProducerBackoffRetry sets duration to backoff retry publish the message
func SetKafkaProducerBackoffRetry(backOffRetry time.Duration) KafkaBrokerConfigFunc {
	return func(kbc *kafkaBrokerConfig) {
		kbc.producerBackoffRetry = backOffRetry
	}
}

// SetKafkaProducerACK sets kafka prodcuer to required acks
func SetKafkaProducerACK(requiredAcks sarama.RequiredAcks) KafkaBrokerConfigFunc {
	return func(kbc *kafkaBrokerConfig) {
		kbc.producerAck = requiredAcks
	}
}

// SetKafkaProducerReturnSuccess set kafka-producer sent the success return
func SetKafkaProducerReturnSuccess(isReturnSuccess bool) KafkaBrokerConfigFunc {
	return func(kbc *kafkaBrokerConfig) {
		kbc.producerReturnSuccess = isReturnSuccess
	}
}

// SetKafkaPublisher set publisher of kafka
func SetKafkaPublisher(publisher abstract.Publisher) KafkaBrokerConfigFunc {
	return func(kbc *kafkaBrokerConfig) {
		kbc.publisher = publisher
	}
}

// SetKafkaSASLUser set authenticated user
func SetKafkaSASLUser(user string) KafkaBrokerConfigFunc {
	return func(kbc *kafkaBrokerConfig) {
		kbc.saslUser = user
	}
}

// SetKafkaSASLPassword set authenticated Password
func SetKafkaSASLPassword(password string) KafkaBrokerConfigFunc {
	return func(kbc *kafkaBrokerConfig) {
		kbc.saslPassword = password
	}
}

// SetKafkaSASLEnabled set authenticated enabled
func SetKafkaSASLEnabled(enabled bool) KafkaBrokerConfigFunc {
	return func(kbc *kafkaBrokerConfig) {
		kbc.saslEnabled = enabled
	}
}

// SetKafkaProducerAsync set kafka producer to async
func SetKafkaProducerAsync(async bool) KafkaBrokerConfigFunc {
	return func(kbc *kafkaBrokerConfig) {
		kbc.producerAsync = async
	}
}

// defaultOptionKafka connection
func defaultOptionKafka() *kafkaBrokerConfig {
	return &kafkaBrokerConfig{
		brokerHosts:           env.GetListString("KAFKA_HOSTS"),
		kafkaVersion:          env.GetString("KAFKA_VERSION", "2.0.0"),
		saslEnabled:           env.GetBool("KAFKA_SASL_ENABLED", false),
		saslUser:              env.GetString("KAFKA_SASL_USER"),
		saslPassword:          env.GetString("KAFKA_SASL_PASSWORD"),
		clientId:              env.GetString("KAFKA_CLIENT_ID"),
		producerMaxRetry:      env.GetInt("KAFKA_PRODUCER_MAX_RETRY", 20),
		producerBackoffRetry:  env.GetDuration("KAFKA_PRODUCER_BACKOFF_RETRY", time.Duration(1)*time.Second),
		producerAck:           sarama.NoResponse,
		producerReturnSuccess: true,
		producerAsync:         true,
	}
}

func NewKafkaBroker(opts ...KafkaBrokerConfigFunc) abstract.Broker {
	logger.PurpleItalic("Load kafka client...")
	opt := defaultOptionKafka()
	for _, o := range opts {
		o(opt)
	}

	cfg := sarama.NewConfig()
	cfg.ClientID = opt.clientId
	cfg.Version, _ = sarama.ParseKafkaVersion(opt.kafkaVersion)
	// producer configuration
	cfg.Producer.Retry.Max = opt.producerMaxRetry
	cfg.Producer.Retry.Backoff = opt.producerBackoffRetry
	cfg.Producer.RequiredAcks = opt.producerAck
	cfg.Producer.Return.Successes = opt.producerReturnSuccess
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner

	// if sasl is enabled
	if opt.saslEnabled {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = opt.saslUser
		cfg.Net.SASL.Password = opt.saslPassword
	}

	kb := new(kafkaBroker)
	kb.brokerHosts = opt.brokerHosts
	kb.config = cfg
	if opt.publisher != nil {
		kb.publisher = opt.publisher
	}

	// connect to broker clients
	client, err := sarama.NewClient(kb.brokerHosts, cfg)
	if err != nil {
		logger.Log.Fatalf("Kafka: cannot connect to server brokers: %s", err)
	}
	kb.client = client

	// is publisher nil? if yes, we will create new publisher
	// otherwise will use existing publisher
	if kb.publisher == nil {
		if opt.producerAsync {
			kb.publisher = NewKafkaAsyncPublisher(kb.client)
		} else {
			kb.publisher = NewKafkaPublisher(kb.client)
		}
	}

	logger.GreenItalic("kafka client connected!")
	return kb
}

func (kb *kafkaBroker) GetConfiguration() interface{} {
	return kb.client
}

func (kb *kafkaBroker) GetPublisher() abstract.Publisher {
	return kb.publisher
}

func (kb *kafkaBroker) GetName() constants.Worker {
	return constants.Kafka
}

func (kb *kafkaBroker) Disconnect(ctx context.Context) error {
	logger.RedBold("kafka: disconnecting...")
	defer fmt.Printf("\x1b[31;1mKafka Disconnecting:\x1b[0m \x1b[32;1mSUCCESS\x1b[0m\n")

	return kb.client.Close()
}

type kafkaBroker struct {
	brokerHosts []string
	config      *sarama.Config
	client      sarama.Client
	publisher   abstract.Publisher
}

type kafkaPublisher struct {
	producer sarama.SyncProducer
}

func NewKafkaPublisher(client sarama.Client) abstract.Publisher {
	logger.PurpleItalic("Load kafka producer connection...")
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		logger.Log.Fatalf("Kafka Publisher: cannot connect with the exist client: %s", err)
	}

	logger.GreenItalic("kafka producer connected!")
	return &kafkaPublisher{producer: producer}
}

func (kp *kafkaPublisher) PublishMessage(ctx context.Context, req *types.PublisherArgument) error {
	trace := tracer.StartTrace(ctx, "Kafka:PublishMessage")

	trace.SetTag("topic", req.Topic)
	trace.SetTag("body", req.Message)
	trace.SetTag("headers", req.Header)

	msg := &sarama.ProducerMessage{
		Topic:     req.Topic,
		Key:       sarama.StringEncoder(req.Key),
		Value:     sarama.ByteEncoder(req.Message),
		Timestamp: time.Now().In(zone.TzJakarta()),
	}

	// if header is empty, create new header
	if len(req.Header) < 1 || req.Header == nil {
		req.Header = make(map[string]interface{})
	}
	// set username to header
	req.Header["username"] = logger.GetUsername(ctx)

	for key, val := range req.Header {
		value, _ := convert.InterfaceToBytes(val)
		msg.Headers = append(msg.Headers, sarama.RecordHeader{
			Key:   []byte(key),
			Value: value,
		})
	}

	reqMessage := req.Message
	if len(reqMessage) > env.GetInt("MAX_BODY_SIZE", 1500) {
		reqMessage = []byte(fmt.Sprintf("request body too long %d", len(reqMessage)))
	}

	headers, _ := convert.InterfaceToString(msg.Headers)
	ol := logger.OutgoingLog{
		StartTime:     msg.Timestamp.Format(constants.LayoutDateTime),
		TargetService: constants.Kafka.String(),
		URL:           fmt.Sprintf("topic: %s", req.Topic),
		RequestBody:   string(reqMessage),
		RequestHeader: headers,
		StatusCode:    http.StatusOK,
	}

	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", err)
		}

		if err != nil {
			ol.StatusCode = http.StatusInternalServerError
			ol.ResponseBody = fmt.Sprintf("%s", err)
		}
		since := time.Since(msg.Timestamp)
		ol.ExecutionTime = since.Seconds()

		ol.Store(ctx)
		trace.Finish()
		monitoring.RecordPrometheus(ol.StatusCode, constants.Kafka.String(), ol.URL, since)
	}()

	_, _, err = kp.producer.SendMessage(msg)
	if err != nil {
		trace.SetError(err)

		return err
	}

	return nil
}

func (kp *kafkaPublisher) PublishMessages(ctx context.Context, reqs []*types.PublisherArgument) error {
	trace := tracer.StartTrace(ctx, "Kafka:PublishMessages")
	defer trace.Finish()

	var messages []*sarama.ProducerMessage
	for _, req := range reqs {
		trace.SetTag("topic", req.Topic)
		trace.SetTag("body", req.Message)
		trace.SetTag("headers", req.Header)

		msg := &sarama.ProducerMessage{
			Topic:     req.Topic,
			Key:       sarama.StringEncoder(req.Key),
			Value:     sarama.ByteEncoder(req.Message),
			Timestamp: time.Now().In(zone.TzJakarta()),
		}
		for key, val := range req.Header {
			value, _ := convert.InterfaceToBytes(val)
			msg.Headers = append(msg.Headers, sarama.RecordHeader{
				Key:   []byte(key),
				Value: value,
			})
		}

		reqMessage := req.Message
		if len(reqMessage) > env.GetInt("MAX_BODY_SIZE", 1500) {
			reqMessage = []byte(fmt.Sprintf("request body too long %d", len(reqMessage)))
		}

		headers, _ := convert.InterfaceToString(msg.Headers)

		ol := logger.OutgoingLog{
			StartTime:     msg.Timestamp.Format(constants.LayoutDateTime),
			TargetService: constants.Kafka.String(),
			URL:           fmt.Sprintf("topic: %s", req.Topic),
			RequestBody:   string(reqMessage),
			RequestHeader: headers,
			StatusCode:    http.StatusOK,
		}

		since := time.Since(msg.Timestamp)
		ol.ExecutionTime = since.Seconds()

		ol.Store(ctx)
		monitoring.RecordPrometheus(ol.StatusCode, constants.Kafka.String(), ol.URL, since)

		messages = append(messages, msg)
	}

	err := kp.producer.SendMessages(messages)
	if err != nil {
		trace.SetError(err)
		return err
	}

	return nil
}

// kafkaAsyncPublisher is a kafka publisher that uses async producer
type kafkaAsyncPublisher struct {
	producer sarama.AsyncProducer
}

// NewKafkaAsyncPublisher creates a new kafka async publisher
func NewKafkaAsyncPublisher(client sarama.Client) abstract.Publisher {
	logger.PurpleItalic("Load kafka async producer connection...")

	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		logger.Log.Fatalf("Kafka Publisher: cannot connect with the exist client: %s", err)
	}

	if strings.EqualFold(env.GetString("OPENSEARCH_FLOW_SEND_DATA"), "queue") {
		// logger.OpenSearch(producer)
		logger.Elasticsearch(producer)
	} else {
		logger.Elasticsearch(producer)
	}

	logger.GreenItalic("kafka async producer connected!")
	return &kafkaAsyncPublisher{producer: producer}
}

// PublishMessage publishes a single message to kafka
func (kap *kafkaAsyncPublisher) PublishMessage(ctx context.Context, req *types.PublisherArgument) error {
	trace := tracer.StartTrace(ctx, "KafkaAsync:PublishMessage")

	trace.SetTag("topic", req.Topic)
	trace.SetTag("body", req.Message)
	trace.SetTag("headers", req.Header)

	msg := &sarama.ProducerMessage{
		Topic:     req.Topic,
		Key:       sarama.StringEncoder(req.Key),
		Value:     sarama.ByteEncoder(req.Message),
		Timestamp: time.Now().In(zone.TzJakarta()),
	}

	// if header is empty, create new header
	if len(req.Header) < 1 || req.Header == nil {
		req.Header = make(map[string]interface{})
	}
	// set username to header
	req.Header["username"] = logger.GetUsername(ctx)

	for key, val := range req.Header {
		value, _ := convert.InterfaceToBytes(val)
		msg.Headers = append(msg.Headers, sarama.RecordHeader{
			Key:   []byte(key),
			Value: value,
		})
	}

	reqMessage := req.Message
	if len(reqMessage) > env.GetInt("MAX_BODY_SIZE", 1500) {
		reqMessage = []byte(fmt.Sprintf("request body too long %d", len(reqMessage)))
	}

	headers, _ := convert.InterfaceToString(msg.Headers)
	ol := &logger.OutgoingLog{
		StartTime:     msg.Timestamp.Format(constants.LayoutDateTime),
		TargetService: constants.Kafka.String(),
		URL:           fmt.Sprintf("topic: %s", req.Topic),
		RequestBody:   string(reqMessage),
		RequestHeader: headers,
		StatusCode:    http.StatusOK,
	}

	var err error
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", err)
		}

		if err != nil {
			ol.StatusCode = http.StatusInternalServerError
			ol.ResponseBody = fmt.Sprintf("%s", err)
		}
		since := time.Since(msg.Timestamp)
		ol.ExecutionTime = since.Seconds()

		ol.Store(ctx)
		trace.Finish()
		monitoring.RecordPrometheus(ol.StatusCode, constants.Kafka.String(), ol.URL, since)
	}()

	kap.producer.Input() <- msg
	select {
	case err := <-kap.producer.Errors():
		if err != nil {
			return err
		}

		return nil
	case <-kap.producer.Successes():
		return nil
	}
}

// PublishMessages publishes multiple messages to kafka
func (kap *kafkaAsyncPublisher) PublishMessages(ctx context.Context, reqs []*types.PublisherArgument) error {
	trace := tracer.StartTrace(ctx, "Kafka:PublishMessages")
	defer trace.Finish()

	for _, req := range reqs {
		trace.SetTag("topic", req.Topic)
		trace.SetTag("body", req.Message)
		trace.SetTag("headers", req.Header)

		msg := &sarama.ProducerMessage{
			Topic:     req.Topic,
			Key:       sarama.StringEncoder(req.Key),
			Value:     sarama.ByteEncoder(req.Message),
			Timestamp: time.Now().In(zone.TzJakarta()),
		}
		for key, val := range req.Header {
			value, _ := convert.InterfaceToBytes(val)
			msg.Headers = append(msg.Headers, sarama.RecordHeader{
				Key:   []byte(key),
				Value: value,
			})
		}

		reqMessage := req.Message
		if len(reqMessage) > env.GetInt("MAX_BODY_SIZE", 1500) {
			reqMessage = []byte(string(reqMessage[:env.GetInt("MAX_BODY_SIZE", 1500)-3]) + "...")
		}

		headers, _ := convert.InterfaceToString(msg.Headers)

		ol := logger.OutgoingLog{
			StartTime:     msg.Timestamp.Format(constants.LayoutDateTime),
			TargetService: constants.Kafka.String(),
			URL:           fmt.Sprintf("topic: %s", req.Topic),
			RequestBody:   string(reqMessage),
			RequestHeader: headers,
			StatusCode:    http.StatusOK,
		}

		since := time.Since(msg.Timestamp)
		ol.ExecutionTime = since.Seconds()

		ol.Store(ctx)
		monitoring.RecordPrometheus(ol.StatusCode, constants.Kafka.String(), ol.URL, since)

		kap.producer.Input() <- msg
		select {
		case err := <-kap.producer.Errors():
			if err != nil {
				return err
			}
			continue
		case <-kap.producer.Successes():
			continue
		}
	}

	return nil
}
