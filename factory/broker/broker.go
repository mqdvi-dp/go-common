package broker

import (
	"context"

	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/logger"
)

type Broker struct {
	brokers map[constants.Worker]abstract.Broker
}

// New register all broker for publisher and/or consumer
func New(brokers ...abstract.Broker) *Broker {
	broker := &Broker{
		brokers: make(map[constants.Worker]abstract.Broker),
	}

	for _, bk := range brokers {
		if _, ok := broker.brokers[bk.GetName()]; ok {
			logger.Log.Fatalf("Duplicate register broker: %s", bk.GetName())
		}

		broker.brokers[bk.GetName()] = bk
	}

	return broker
}

// GetBrokers return all registered brokers
func (b *Broker) GetBrokers() map[constants.Worker]abstract.Broker {
	return b.brokers
}

// RegisterBroker register enw broker
func (b *Broker) RegisterBroker(brokerName constants.Worker, broker abstract.Broker) {
	if b.brokers == nil {
		b.brokers = make(map[constants.Worker]abstract.Broker)
	}

	if _, ok := b.brokers[broker.GetName()]; ok {
		logger.Log.Fatalf("Duplicate register broker: %s", broker.GetName())
	}

	b.brokers[brokerName] = broker
}

// Disconnect disconnect all registered brokers
func (b *Broker) Disconnect(ctx context.Context) error {
	for _, broker := range b.brokers {
		broker.Disconnect(ctx)
	}

	return nil
}
