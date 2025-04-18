package kafka

import (
	"time"

	"github.com/IBM/sarama"
	"github.com/mqdvi-dp/go-common/env"
)

type OptionFunc func(*option)

type option struct {
	serviceName     string
	brokerHosts     []string
	saslEnabled     bool
	saslUser        string
	saslPassword    string
	consumerGroup   string
	maxGoroutines   int
	offsetInitial   int64
	balanceStrategy sarama.BalanceStrategy
	retryBackoff    time.Duration
	maxRetry        int
	isNeedProducer  bool
}

func getDefaultOption() option {
	return option{
		brokerHosts:     env.GetListString("KAFKA_HOSTS"),
		saslEnabled:     env.GetBool("KAFKA_SASL_ENABLED", false),
		saslUser:        env.GetString("KAFKA_SASL_USER"),
		saslPassword:    env.GetString("KAFKA_SASL_PASSWORD"),
		maxGoroutines:   env.GetInt("BROKER_MAX_GOROUTINES", 20),
		offsetInitial:   sarama.OffsetOldest,
		balanceStrategy: sarama.NewBalanceStrategyRoundRobin(),
		retryBackoff:    env.GetDuration("KAFKA_CONSUMER_RETRY_BACKOFF", 2*time.Second),
		maxRetry:        5,
		isNeedProducer:  true,
	}
}

// SetBrokerHosts set kafka hosts
func SetBrokerHosts(hosts []string) OptionFunc {
	return func(o *option) {
		o.brokerHosts = hosts
	}
}

// SetConsumerGroup set consumer group of kafka
func SetConsumerGroup(consumerGroup string) OptionFunc {
	return func(o *option) {
		o.consumerGroup = consumerGroup
	}
}

// SetMaxGoroutines set maximum of concurrency kafka
func SetMaxGoroutines(maxGoroutines int) OptionFunc {
	return func(o *option) {
		o.maxGoroutines = maxGoroutines
	}
}

// SetOffsetInitial set kafka-consumer to set offset when initial consumer
func SetOffsetInitial(offset int64) OptionFunc {
	return func(o *option) {
		o.offsetInitial = offset
	}
}

// SetKafkaBalanceStrategy set consumer kafka to set balance strategy
func SetKafkaBalanceStrategy(strategy sarama.BalanceStrategy) OptionFunc {
	return func(o *option) {
		o.balanceStrategy = strategy
	}
}

// SetRetryBackOff set publisher of kafka
func SetRetryBackOff(retryBackOff time.Duration) OptionFunc {
	return func(o *option) {
		o.retryBackoff = retryBackOff
	}
}

// SetMaxRetry set max retry re-queue message
func SetMaxRetry(maxRetry int) OptionFunc {
	return func(o *option) {
		o.maxRetry = maxRetry
	}
}

// SetIsNeedProducer set consumer need producer?
// used for re-queue
func SetIsNeedProducer(needProducer bool) OptionFunc {
	return func(o *option) {
		o.isNeedProducer = needProducer
	}
}

// SetSASLUser set authenticated user
func SetSASLUser(user string) OptionFunc {
	return func(o *option) {
		o.saslUser = user
	}
}

// SetSASLPassword set authenticated Password
func SetSASLPassword(password string) OptionFunc {
	return func(o *option) {
		o.saslPassword = password
	}
}

// SetSASLEnabled set authenticated Password
func SetSASLEnabled(enabled bool) OptionFunc {
	return func(o *option) {
		o.saslEnabled = enabled
	}
}
