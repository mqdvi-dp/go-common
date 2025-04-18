package abstract

import "github.com/mqdvi-dp/go-common/constants"

type Broker interface {
	GetConfiguration() interface{}
	GetPublisher() Publisher
	GetName() constants.Worker
	Closer
}
