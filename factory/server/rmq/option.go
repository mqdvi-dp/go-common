package rmq

import "github.com/mqdvi-dp/go-common/env"

type option struct {
	consumerGroup string
	exchangeName  string
	queue         string
	broker        string
	maxGoroutines int
	debugMode     bool
	serviceName   string
}

type OptionFunc func(*option)

func getDefaultOption() option {
	return option{
		maxGoroutines: env.GetInt("BROKER_MAX_GOROUTINES", 20),
		debugMode:     env.GetBool("DEBUG_MODE"),
	}
}

// SetMaxGoroutines option func
func SetMaxGoroutines(maxGoroutines int) OptionFunc {
	return func(o *option) {
		o.maxGoroutines = maxGoroutines
	}
}

// SetDebugMode option func
func SetDebugMode(debugMode bool) OptionFunc {
	return func(o *option) {
		o.debugMode = debugMode
	}
}

// SetBrokerHost option func
func SetBrokerHost(broker string) OptionFunc {
	return func(o *option) {
		o.broker = broker
	}
}

// SetConsumerGroup option func, for log when startup
func SetConsumerGroup(consumerGroup string) OptionFunc {
	return func(o *option) {
		o.consumerGroup = consumerGroup
	}
}

// SetExchangeName option func
func SetExchangeName(exchangeName string) OptionFunc {
	return func(o *option) {
		o.exchangeName = exchangeName
	}
}

// SetServiceName option func
func SetServiceName(serviceName string) OptionFunc {
	return func(o *option) {
		o.serviceName = serviceName
	}
}
