package rmq

import (
	"encoding/json"

	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/logger"
)

type queueKey struct {
	Exchange  string `json:"exchange"`
	QueueName string `json:"queue_name"`
	Channel   string `json:"channel"`
}

func (q *queueKey) String() string {
	b, _ := json.Marshal(q)

	return string(b)
}

func QueueKey(exchange, queueName, channel string) string {
	qk := &queueKey{
		Exchange:  exchange,
		QueueName: queueName,
		Channel:   channel,
	}

	return qk.String()
}

func ParseQueueKey(val string) (exchange, queueName, channel string) {
	var qk = new(queueKey)

	err := json.Unmarshal([]byte(val), &qk)
	if err != nil {
		logger.Log.Fatalf("error parse queue key: %s", err)
		return
	}

	exchange = qk.Exchange
	queueName = qk.QueueName
	channel = convert.StringWithUnderscored(qk.Channel)
	return
}
