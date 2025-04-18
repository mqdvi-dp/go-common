package logger

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/mqdvi-dp/go-common/compare"
	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/zone"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/sirupsen/logrus"
)

// Finalize end of journey for logging application requests
func (l *Logger) Finalize(ctx context.Context) {
	value, ok := extract(ctx)
	if !ok {
		return
	}

	// get username from context
	val, ok := value.LoadAndDelete(_Username)
	if ok {
		l.Username = val.(string)
	}

	// get debugging information from context and store it in the struct
	val, ok = value.LoadAndDelete(_LogMessage)
	if ok {
		l.LogMessage = val.([]LogMessage)
	}

	// get outgoing log request from context
	val, ok = value.LoadAndDelete(_OutgoingLog)
	if ok {
		l.OutgoingLog = val.([]OutgoingLog)
	}

	// get database log from context
	val, ok = value.LoadAndDelete(_Database)
	if ok {
		l.Database = val.([]Database)
	}

	// get error message from context
	val, ok = value.LoadAndDelete(_ErrorMessage)
	if ok {
		l.ErrorMessage = val.(string)
	}

	l.logOutput()
}

func (l *Logger) logOutput() {
	var isDebug = env.GetBool("DEBUG", true)
	if compare.String(l.Endpoint, env.GetListString("OPENSEARCH_BLACKLIST_ENDPOINT")...) {
		return
	}

	// check, is outgoing_log more than the limits? if yes, get last index from limits
	maxLengthIndex := env.GetInt("MAX_LENGTH_ARRAY", 10)
	if len(l.Database) > maxLengthIndex {
		l.Database = l.Database[maxLengthIndex-1:]
	}

	if len(l.OutgoingLog) > maxLengthIndex {
		l.OutgoingLog = l.OutgoingLog[maxLengthIndex-1:]
	}

	if len(l.LogMessage) > maxLengthIndex {
		l.LogMessage = l.LogMessage[maxLengthIndex-1:]
	}

	// when debug is true, only sent to stdout, otherwise will publish the log into NSQ
	if isDebug {
		var level = logrus.ErrorLevel // default log level is error

		if l.StatusCode >= http.StatusOK && l.StatusCode < http.StatusBadRequest {
			level = logrus.InfoLevel
		} else if l.StatusCode >= http.StatusBadRequest && l.StatusCode < http.StatusInternalServerError {
			level = logrus.WarnLevel
		}

		logJ.WithField("log", l).Log(level)
		return
	}

	// running in backgrounds
	go func() {
		// send log to opensearch if flow send data directly
		// otherwise, we will send data into queue
		// default is queue
		flowSendData := env.GetString("OPENSEARCH_FLOW_SEND_DATA", "queue")
		if flowSendData == "queue" {
			sendLogToQueue(l)
		} else {
			sendLogDirectly(l)
		}
	}()
}

func sendLogDirectly(l *Logger) {
	if osc == nil {
		logT.Error(fmt.Errorf("osc variable is nil"))
		return
	}

	if osc.client == nil {
		logT.Error(fmt.Errorf("opensearch client is nil"))
		return
	}

	// get the index prefix from env variable
	// default is "service_ms_log"
	prefixIndexName := env.GetString("OPENSEARCH_PREFIX_INDEX", "service_ms_log")
	indexName := prefixIndexName
	// get current environment of service
	serverEnv := env.GetString("ENV")
	if serverEnv != "" && !strings.EqualFold(serverEnv, "production") {
		indexName += fmt.Sprintf("_%s", strings.ToLower(serverEnv))
	}
	// the format index will be:
	// 1. if there has serverEnv "service_ms_log_development-20231020"
	// 2. if not there "service_ms_log-20231020"
	indexName += fmt.Sprintf("-%s", time.Now().In(zone.TzJakarta()).Format("20060102"))
	body, _ := convert.InterfaceToBytes(l)
	reqBody := opensearchapi.IndexRequest{
		Index:      indexName,
		DocumentID: l.RequestId,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	// send to opensearch
	resp, err := reqBody.Do(context.Background(), osc.client)
	if err != nil {
		logT.Error(err)
		return
	}
	defer resp.Body.Close()

	// if there has error
	if resp.IsError() {
		err = fmt.Errorf("send logging into opensearch client is failed: actual response is %s", resp.String())

		logT.Error(err)
		return
	}
}

func sendLogToQueue(l *Logger) {
	if osc == nil {
		logT.Error(fmt.Errorf("osc variable is nil"))
		return
	}

	if osc.producer == nil {
		logT.Error(fmt.Errorf("nsq producer is nil"))
		return
	}

	key := fmt.Sprintf("svc_log: %s", l.Service)

	payload, _ := convert.InterfaceToBytes(l)
	topicLogging := env.GetString("OPENSEARCH_QUEUE", "os_queue")
	msg := &sarama.ProducerMessage{
		Topic:     topicLogging,
		Value:     sarama.ByteEncoder(payload),
		Key:       sarama.StringEncoder(key),
		Timestamp: time.Now().In(zone.TzJakarta()),
	}
	osc.producer.Input() <- msg
	select {
	case err := <-osc.producer.Errors():
		logT.Error(fmt.Errorf("failed to publish data logging: %s", err))
	case <-osc.producer.Successes():
	}
}
