package request

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/monitoring"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/zone"
)

func (r *request) wrapper(ctx context.Context, payload []byte, method string) ([]byte, int, string, error) {
	trace := tracer.StartTrace(ctx, "RequestClient:Wrapper")
	start := time.Now().In(zone.TzJakarta())

	// start log for outgoing request
	log := logger.OutgoingLog{
		StartTime:     start.Format(constants.LayoutDateTime),
		URL:           fmt.Sprintf("%s %s", method, r.url),
		TargetService: r.target,
		RequestHeader: dumpHeader(r.header),
		RequestBody:   string(payload),
	}

	defer func() {
		// set execution time
		since := time.Since(start)
		log.ExecutionTime = since.Seconds()

		trace.Finish()
		log.Store(ctx)
		monitoring.RecordPrometheus(log.StatusCode, method, r.target, since)
	}()
	// set tracer
	trace.SetTag("request.http.url", log.URL)
	trace.SetTag("request.http.header", log.RequestHeader)
	trace.SetTag("request.http.body", log.RequestBody)

	res, respHeader, code, curlCommand, err := r.curl(payload, method)
	// set data response to log
	log.StatusCode = code
	// set data response to trace
	trace.SetTag("response.http.status_code", code)
	trace.SetTag("response.http.header", respHeader)

	if err != nil {
		log.ResponseBody = err.Error()
		trace.SetError(err)
	}

	if res != nil {
		log.ResponseBody = string(res)
		trace.SetTag("response.http.body", res)
	}

	return res, code, curlCommand, err
}

func dumpHeader(header http.Header) string {
	var headers []string
	for key, value := range header {
		headers = append(headers, fmt.Sprintf("%s: %s", key, value))
	}

	return strings.Join(headers, " | ")
}
