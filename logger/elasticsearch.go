package logger

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/mqdvi-dp/go-common/env"
)

// esClient is Elasticsearch client
var esClient *esLogger
var onceEs sync.Once

type esLogger struct {
	client   *elasticsearch.Client
	producer sarama.AsyncProducer
}

// ElasticsearchInit initializes the Elasticsearch client
func Elasticsearch(producer ...sarama.AsyncProducer) {
	if !env.GetBool("ELASTICSEARCH_ENABLED", false) {
		return
	}

	if strings.EqualFold(env.GetString("ELASTICSEARCH_FLOW_SEND_DATA"), "queue") {
		if len(producer) > 0 {
			esClient = &esLogger{producer: producer[0]}
		}
		return
	}

	hosts := env.GetListString("ELASTICSEARCH_HOST")
	if len(hosts) == 0 {
		panic(fmt.Errorf("env ELASTICSEARCH_HOST is empty"))
	}

	var username, password string
	if env.GetBool("ELASTICSEARCH_SECURE", false) {
		username = env.GetString("ELASTICSEARCH_USERNAME")
		if username == "" {
			panic(fmt.Errorf("env ELASTICSEARCH_USERNAME is empty"))
		}
		password = env.GetString("ELASTICSEARCH_PASSWORD")
		if password == "" {
			panic(fmt.Errorf("env ELASTICSEARCH_PASSWORD is empty"))
		}
	}

	cfg := elasticsearch.Config{
		Addresses: hosts,
		Username:  username,
		Password:  password,
		CloudID: env.GetString("ELASTICSEARCH_CLOUD_ID"),
		APIKey: env.GetString("ELASTICSEARCH_API_KEY"),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	esClient = &esLogger{client: client}
}
