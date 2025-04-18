package logger

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/IBM/sarama"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/opensearch-project/opensearch-go"
)

// osc is opensearch client
var osc *osearch
var onceOs sync.Once

type osearch struct {
	client   *opensearch.Client
	producer sarama.AsyncProducer
}

// OpenSearch start the connection
func OpenSearch(producer ...sarama.AsyncProducer) {
	if !env.GetBool("OPENSEARCH_ENABLED", false) {
		return
	}

	if strings.EqualFold(env.GetString("OPENSEARCH_FLOW_SEND_DATA"), "queue") {
		if len(producer) > 0 {
			osc = &osearch{producer: producer[0]}
		}
		return
	}

	host := env.GetListString("OPENSEARCH_HOST")
	// when host is empty, we can skip that
	if len(host) < 1 {
		panic(fmt.Errorf("env OPENSEARCH_HOST is empty"))
	}

	var username, password string
	if env.GetBool("OPENSEARCH_SECURE", false) {
		// get the username and password from env
		// and check that value, if empty, we should send the panic
		osUsername := env.GetString("OPENSEARCH_USERNAME")
		if osUsername == "" {
			panic(fmt.Errorf("env OPENSEARCH_USERNAME is empty"))
		}
		osPassword := env.GetString("OPENSEARCH_PASSWORD")
		if osPassword == "" {
			panic(fmt.Errorf("env OPENSEARCH_PASSWORD is empty"))
		}

		info := url.UserPassword(osUsername, osPassword)
		username = info.Username()
		password, _ = info.Password()
	}

	// init client opensearch
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: host,
		Username:  username,
		Password:  password,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		panic(err)
	}

	osc = &osearch{client: client}
}
