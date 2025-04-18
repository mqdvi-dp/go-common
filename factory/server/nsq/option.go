package nsq

import (
	"time"

	"github.com/mqdvi-dp/go-common/env"
)

type option struct {
	serviceName   string
	topic         string
	channel       string
	brokerHosts   []string
	queueDelay    time.Duration
	maxQueueDelay time.Duration
	msgTimeout    time.Duration
	maxAttempts   int
	maxInflight   int
	maxGoroutines int
}

type OptionFunc func(*option)

func getDefaultOption() option {
	return option{
		topic:         "default-topic",
		channel:       "default-channel",
		brokerHosts:   env.GetListString("NSQ_LOOKUPD"),
		maxInflight:   20,
		queueDelay:    time.Duration(1) * time.Minute,
		maxQueueDelay: time.Duration(5) * time.Minute,
		msgTimeout:    time.Duration(1) * time.Minute,
		maxAttempts:   100,
		maxGoroutines: 10,
	}
}

// SetChannel option
func SetChannel(channel string) OptionFunc {
	return func(o *option) {
		o.channel = channel
	}
}

// SetMaxGoroutines option
func SetMaxGoroutines(max int) OptionFunc {
	return func(o *option) {
		o.maxGoroutines = max
	}
}

// SetBrokerHost option
func SetBrokerHost(brokerHosts []string) OptionFunc {
	return func(o *option) {
		o.brokerHosts = brokerHosts
	}
}

// SetTopic option
func SetTopic(topic string) OptionFunc {
	return func(o *option) {
		o.topic = topic
	}
}

// SetQueueDelay option
func SetQueueDelay(delay time.Duration) OptionFunc {
	return func(o *option) {
		o.queueDelay = delay
	}
}

// SetMaxQueueDelay option
func SetMaxQueueDelay(maxDelay time.Duration) OptionFunc {
	return func(o *option) {
		o.maxQueueDelay = maxDelay
	}
}

// SetMaxAttempts option
func SetMaxAttempts(max int) OptionFunc {
	return func(o *option) {
		o.maxAttempts = max
	}
}

// SetMaxInFlight option
func SetMaxInFlight(max int) OptionFunc {
	return func(o *option) {
		o.maxInflight = max
	}
}

// SetMessageTimeout option
func SetMessageTimeout(msgTimeout time.Duration) OptionFunc {
	return func(o *option) {
		o.msgTimeout = msgTimeout
	}
}
