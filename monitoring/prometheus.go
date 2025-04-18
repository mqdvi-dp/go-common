package monitoring

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/mqdvi-dp/go-common/logger"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	histogram *prometheus.HistogramVec
	counter   *prometheus.CounterVec
	once      sync.Once
)

func New(appName string) {
	once.Do(func() {
		// labels for histogram
		labelsHistogram := []string{"status_code", "method", "path"}
		// register histogram metrics
		histogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        "request_duration",
			Help:        "How long it took to process requests from the server and return them",
			ConstLabels: prometheus.Labels{"application": appName},
		}, labelsHistogram)

		if err := prometheus.Register(histogram); err != nil {
			logger.Log.Fatalf("failed to register histogram with an error %s", err)
		}

		// labels for counter vector
		labelsCounter := []string{"state", "method", "path"}
		// register counter metrics
		counter = prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "request_total",
			Help:        "How many requests were processed by state, method and path",
			ConstLabels: prometheus.Labels{"application": appName},
		}, labelsCounter)

		if err := prometheus.Register(counter); err != nil {
			logger.Log.Fatalf("failed to register counter with an error %s", err)
		}
	})
}

const (
	successState       = "ok"
	businessErrorState = "business_error"
	systemErrorState   = "system_error"
)

func RecordPrometheus(statusCode int, method, path string, duration time.Duration) {
	// make sure an instance of counter and histogram registered
	if counter == nil || histogram == nil {
		return
	}

	// define variable state
	var state string
	switch {
	case statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError:
		state = businessErrorState
	case statusCode >= http.StatusInternalServerError:
		state = systemErrorState
	default:
		state = successState
	}

	// record data to metrics
	// counter metrics
	counter.WithLabelValues(state, method, path).Inc()
	// histogram metrics
	histogram.WithLabelValues(strconv.Itoa(statusCode), method, path).Observe(duration.Seconds())
}
