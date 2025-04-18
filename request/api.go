package request

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/mqdvi-dp/go-common/env"
)

type request struct {
	url       string
	target    string
	header    http.Header
	client    *http.Client
	basicAuth *basicAuth
}

type basicAuth struct {
	username string
	password string
	set      bool
}

type ApiClient interface {
	Request(header http.Header, target, url string) MethodInterface
	RequestWithBasicAuth(header http.Header, username, password, target, url string) MethodInterface
}

func NewApiClient(clients ...*http.Client) ApiClient {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: env.GetDuration("HTTP_TIMEOUT", time.Duration(20)*time.Second),
	}

	if len(clients) > 0 {
		client = clients[0]
	}

	return &request{client: client}
}

func NewApiClientWithDeadline(clients ...*http.Client) ApiClient {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: env.GetDuration("HTTP_TIMEOUTS", time.Duration(1)*time.Millisecond),
	}

	if len(clients) > 0 {
		client = clients[0]
	}

	return &request{client: client}
}

func (r *request) Request(header http.Header, target, url string) MethodInterface {
	var req = &request{
		client:    r.client,
		url:       url,
		header:    header,
		target:    target,
		basicAuth: &basicAuth{},
	}

	return req
}

func (r *request) RequestWithBasicAuth(header http.Header, username, password, target, url string) MethodInterface {
	var req = &request{
		client: r.client,
		url:    url,
		header: header,
		target: target,
		basicAuth: &basicAuth{
			username: username,
			password: password,
			set:      true,
		},
	}

	return req
}
