package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/errs"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/monitoring"
	"github.com/mqdvi-dp/go-common/response"
	"github.com/mqdvi-dp/go-common/tracer"
)

func (r *rest) restTracerMiddleware(c *gin.Context) {
	req := c.Request
	start := time.Now().In(r.tz)
	ctx := req.Context()

	var statusCode int
	var err error

	// init logger data
	log := &logger.Logger{
		StartTime:   start.Format(time.RFC3339),
		HandlerType: logger.Http,
		Service:     r.service.Name(),
		Endpoint:    fmt.Sprintf("%s %s", req.Method, parseQueryParam(c)),
		Ip:          c.GetHeader(constants.ApplicationOriginalIp),
		DeviceId:    c.GetHeader(constants.ApplicationDevice),
		Timezone:    c.GetHeader(constants.ApplicationTimezone),
		Lat:         convert.StringToFloat(c.GetHeader(constants.ApplicationLatitude)),
		Lng:         convert.StringToFloat(c.GetHeader(constants.ApplicationLongitude)),
		Source:      strings.ToUpper(c.GetHeader(constants.ApplicationChannel)),
	}

	// create operationName
	shortUrl := parseUrl(c)
	operationName := fmt.Sprintf("%s %s", req.Method, shortUrl)

	trace, ctx := tracer.StartTraceWithContext(ctx, operationName)
	defer func() {
		if rec := recover(); rec != nil {
			err = fmt.Errorf("%s", rec)
			fmt.Println(err)
			debug.PrintStack()
		}

		if err != nil {
			trace.SetError(err)
			trace.SetTag("error", err)
			log.ErrorMessage = fmt.Sprintf("%v", err)

			if statusCode < http.StatusOK {
				statusCode = http.StatusInternalServerError
			}
		}
		log.StatusCode = statusCode

		// set execution time
		since := time.Since(start)
		log.ExecutionTime = since.Seconds()
		trace.SetTag("trace_id", tracer.GetTraceId(ctx))
		trace.Finish()

		// print log stdout
		log.Finalize(ctx)
		monitoring.RecordPrometheus(log.StatusCode, req.Method, shortUrl, since)
	}()
	// get traceId from activeTracer
	traceId := tracer.GetTraceId(ctx)
	// if activeTracer didn't have traceId, default is uuid.New v4
	if traceId == "" {
		traceId = uuid.NewString()
	}
	// set log.RequestId with tracerId
	log.RequestId = traceId
	logger.SetRequestId(ctx, log.RequestId)

	// implement locking logging stdout
	var lock = logger.NewLocker(ctx)
	// set to context with logger.LogKey as a context key
	ctx = context.WithValue(ctx, logger.LogKey, lock)

	maxBodySize := env.GetInt("MAX_BODY_SIZE", 1500)

	if req.Method == http.MethodPost || req.Method == http.MethodPut {
		reqBody, err := io.ReadAll(req.Body)
		reqBodyMasked := string(logger.MaskedCredentials(reqBody))
		defer req.Body.Close()
		if err != nil {
			return
		}

		if len(reqBody) > maxBodySize {
			trace.SetTag("http.request.body.size", len(reqBody))
			log.RequestBody = fmt.Sprintf("length body: %d", len(reqBody))
		} else {
			log.RequestBody = reqBodyMasked
			trace.SetTag("http.request.body", reqBodyMasked)
		}

		// put again into request body
		req.Body = io.NopCloser(bytes.NewReader(reqBody))
	}

	// dump request
	httpDump := dumpRequestHeader(c)
	httpDumpMasked := string(logger.MaskedCredentials([]byte(httpDump)))
	trace.SetTag("http.request.header", httpDumpMasked)
	trace.SetTag("http.url", log.Endpoint)
	// set request header and body to logger
	log.RequestHeader = httpDumpMasked

	// set gin context with tracer context
	c.Request = req.WithContext(ctx)
	// go to the next handler
	c.Next()

	// get status code of response
	statusCode = c.Writer.Status()
	trace.SetTag("http.status_code", statusCode)

	// get response body
	// NOTES: trash way to get http response body
	// TODO: need to fixed how to get http response body
	respBody := []byte(c.GetString("HTTP_RESPONSE_BODY"))
	respBodyMasked := logger.MaskedCredentials(respBody)

	if len(respBody) > maxBodySize {
		trace.SetTag("http.response.body.size", len(respBody))
	} else {
		trace.SetTag("http.response.body", respBodyMasked)
	}

	// set response header and body to logger
	log.ResponseBody = string(respBodyMasked)
}

// dumpRequestHeader Request Header return to string
func dumpRequestHeader(c *gin.Context) string {
	// var headers []string
	// for key, value := range c.Request.Header {
	// 	headers = append(headers, fmt.Sprintf("%s: %s", key, strings.Join(value, ", ")))
	// }

	// return strings.Join(headers, " || ")
	headers, err := convert.InterfaceToString(c.Request.Header)
	if err != nil {
		return ""
	}

	return headers
}

// parseUrl request url
func parseUrl(c *gin.Context) string {
	var url = c.Request.URL.Path

	for idx := range c.Params {
		key := c.Params[idx].Key
		val := c.Params[idx].Value
		// when url value contains with value from urlParams
		if strings.Contains(url, val) {
			// get index characters
			index := strings.Index(url, val)

			// replace value with key of params
			url = fmt.Sprintf("%s:%s%s", url[:index], key, url[len(val)+index:])
		}
	}

	return url
}

// parseQueryParam parse and masked the credentials data at endpoint
func parseQueryParam(c *gin.Context) string {
	var endpoint = c.Request.URL.Path
	var queryParams = c.Request.URL.Query()
	queryBytes, _ := convert.InterfaceToBytes(queryParams)
	queryMasking := logger.MaskedCredentials(queryBytes)

	var newQueryParams = make(map[string]interface{})
	err := json.Unmarshal(queryMasking, &newQueryParams)
	if err != nil {
		return c.Request.RequestURI
	}

	queryParams = url.Values{}
	for key, val := range newQueryParams {
		var value string
		switch v := val.(type) {
		case string:
			value = v
		case []interface{}:
			if vv, ok := v[0].(string); ok {
				value = vv
			}
		}

		queryParams.Set(key, value)
	}

	if queryParams.Encode() == "" {
		return c.Request.RequestURI
	}

	enc, err := url.QueryUnescape(fmt.Sprintf("%s?%s", endpoint, queryParams.Encode()))
	if err != nil {
		return c.Request.RequestURI
	}

	return enc
}

// validateAllowedHeaders validate allowed headers
func (r *rest) validateAllowedHeaders(c *gin.Context) {
	req := c.Request
	ctx := req.Context()

	channel := req.Header.Get(constants.ApplicationChannel)
	if channel == "" {
		return
	}

	allowedHeaders := env.GetListString(
		"HTTP_ALLOWED_HEADERS_FOR_NON_MOBILE_CHANNEL",
		// default values for allowed headers
		constants.ApplicationSignature,
		constants.ApplicationTimestamp,
		constants.ApplicationDevice,
		constants.ApplicationChannel,
	)

	if strings.EqualFold(channel, "m") {
		allowedHeaders = env.GetListString(
			"HTTP_ALLOWED_HEADERS_FOR_MOBILE_CHANNEL",
			// default values for allowed headers
			constants.ApplicationSignature,
			constants.ApplicationTimestamp,
			constants.ApplicationDevice,
			constants.ApplicationChannel,
			constants.ApplicationVersion,
		)
	}

	for _, ah := range allowedHeaders {
		// when header not exists
		if req.Header.Get(ah) == "" {
			response.Error(ctx, errs.NewErrorWithCodeErr(fmt.Errorf("header %s not exists", ah), errs.BAD_REQUEST_HEADER)).JSON(c)
			return
		}
	}

	c.Next()
}
