package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"runtime/debug"
	"time"

	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/errs"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/monitoring"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/zone"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const ands = "&&&"

type interceptor struct {
	opt         *option
	serviceName string
	host        string
}

func NewInterceptor(host string) *interceptor {
	return &interceptor{host: host}
}

func (i *interceptor) ChainUnaryClient(interceptors ...grpc.UnaryClientInterceptor) grpc.UnaryClientInterceptor {
	n := len(interceptors)

	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		chainer := func(currentIntercept grpc.UnaryClientInterceptor, currentInvoker grpc.UnaryInvoker) grpc.UnaryInvoker {
			return func(
				currentCtx context.Context,
				currentMethod string,
				currentReq, currentReply interface{},
				currentClientConn *grpc.ClientConn,
				currentOpts ...grpc.CallOption,
			) error {
				return currentIntercept(
					currentCtx,
					currentMethod,
					currentReq,
					currentReply,
					currentClientConn,
					currentInvoker,
					currentOpts...,
				)
			}
		}

		chainedInvoker := invoker
		for i := n - 1; i >= 0; i-- {
			chainedInvoker = chainer(interceptors[i], chainedInvoker)
		}

		return chainedInvoker(ctx, method, req, reply, cc, opts...)
	}
}

func (i *interceptor) UnaryClientTracerInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) (err error) {
	start := time.Now().In(zone.TzJakarta())

	ol := logger.OutgoingLog{
		StartTime:     start.Format(constants.LayoutDateTime),
		TargetService: i.host,
		URL:           method,
		RequestBody:   dumpData(req),
	}

	trace, ctx := tracer.StartTraceWithContext(ctx, ol.URL)
	defer func() {
		if r := recover(); r != nil {
			err = status.Errorf(codes.Aborted, "%s", r)
			fmt.Println(err)
			debug.PrintStack()
		}

		var sc = http.StatusOK
		if err != nil {
			sc = convertErrorToStatusCode(err)
			c := convertStatusCodeToCodes(sc)

			err = status.Error(c, err.Error())
			trace.SetError(err)
			ol.ResponseBody = fmt.Sprintf("%s", err)
		}

		// set tracer response
		trace.SetTag("response.grpc.body", dumpData(reply))
		trace.SetTag("trace_id", tracer.GetTraceId(ctx))
		// set logger
		ol.StatusCode = sc
		ol.ExecutionTime = time.Since(start).Seconds()
		if reflect.ValueOf(ol.ResponseBody).IsZero() {
			ol.ResponseBody = dumpData(reply)
		}

		trace.Finish()
		monitoring.RecordPrometheus(sc, constants.GRPC.String(), method, time.Duration(ol.ExecutionTime)*time.Second)
		ol.Store(ctx)
	}()

	// set metadata grpc (header)
	md, ok := metadata.FromOutgoingContext(ctx)
	// if metadata is not found, create new metadata
	if !ok {
		md = metadata.New(map[string]string{})
	}
	// set username to metadata if exists
	username := logger.GetUsername(ctx)
	if username != "" {
		md.Set("username", username)
	}
	ol.RequestHeader = dumpData(md)
	// set metadata to context
	ctx = metadata.NewOutgoingContext(ctx, md)

	trace.SetTag("request.grpc.endpoint", ol.URL)
	trace.SetTag("request.grpc.body", ol.RequestBody)
	trace.SetTag("request.grpc.metadata", ol.RequestHeader)

	err = invoker(ctx, method, req, reply, cc, opts...)
	return
}

// for unary a server
// chainUnaryserver creates a single interceptor out of a chain of many interceptors
func (i *interceptor) chainUnaryServer(interceptors ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	n := len(interceptors)

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		chainer := func(
			currentInterceptor grpc.UnaryServerInterceptor,
			currentHandler grpc.UnaryHandler,
		) grpc.UnaryHandler {
			return func(currentCtx context.Context, currentReq interface{}) (interface{}, error) {
				return currentInterceptor(currentCtx, currentReq, info, currentHandler)
			}
		}

		chainedHandler := handler
		for i := n - 1; i >= 0; i-- {
			chainedHandler = chainer(interceptors[i], chainedHandler)
		}

		return chainedHandler(ctx, req)
	}
}

func (i *interceptor) unaryServerTracerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	start := time.Now().In(zone.TzJakarta())

	var header string
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		header = dumpData(md)
	}

	requestId := uuid.NewString()
	if val, ok := md["request_id"]; ok {
		if len(val) > 0 {
			requestId = val[0]
		}
	}

	// set username and get username from context
	var username string
	if usernames, ok := md["username"]; ok {
		if len(usernames) > 0 {
			username = usernames[0]
		}
	}

	log := &logger.Logger{
		StartTime:     start.Format(time.RFC3339),
		RequestId:     requestId,
		HandlerType:   logger.GRPC,
		Service:       i.serviceName,
		Endpoint:      info.FullMethod,
		RequestBody:   dumpData(req),
		RequestHeader: header,
		Username:      username,
	}

	trace, ctx := tracer.StartTraceWithContext(ctx, fmt.Sprintf("GRPC: %s", log.Endpoint))
	defer func() {
		if r := recover(); r != nil {
			err = status.Errorf(codes.Aborted, "%s", r)
			fmt.Println(err)
			debug.PrintStack()
		}

		var sc = http.StatusOK
		if err != nil {
			trace.SetError(err)
			trace.SetTag("error", err)
			log.ErrorMessage = fmt.Sprintf("%s", err)

			// when error is context.DeadlineExceeded or context.Canceled
			// modify the error message and status code
			if errors.Is(err, context.DeadlineExceeded) { // when error is context.DeadlineExceeded
				sc = http.StatusInternalServerError
				err = fmt.Errorf("%d%s%s%s%s", errs.CONTEXT_DEADLINE_EXCEEDED.Code(), ands, errs.CONTEXT_DEADLINE_EXCEEDED.Message(), ands, errs.CONTEXT_DEADLINE_EXCEEDED.MoreInfo())
			} else if errors.Is(err, context.Canceled) { // when error is context.Canceled
				sc = http.StatusInternalServerError
				err = fmt.Errorf("%d%s%s%s%s", errs.CONTEXT_CANCELLED.Code(), ands, errs.CONTEXT_CANCELLED.Message(), ands, errs.CONTEXT_CANCELLED.MoreInfo())
			} else { // when error is not context.DeadlineExceeded or context.Canceled
				sc = convertErrorToStatusCode(err)
				// when error is *errs.Error or errs.CodeErr
				// modify the error message and status code
				switch er := err.(type) {
				case *errs.Error:
					sc = er.StatusCode()
					err = fmt.Errorf("%d%s%s%s%s", er.SystemCode(), ands, er.Message(), ands, er.MoreInfo())
				case errs.CodeErr:
					sc = er.StatusCode()
					err = fmt.Errorf("%d%s%s%s%s", er.Code(), ands, er.Message(), ands, er.MoreInfo())
				}
			}

			c := convertStatusCodeToCodes(sc)
			err = status.Error(c, err.Error())
		}

		// set tracer response
		trace.SetTag("response.grpc.body", dumpData(resp))
		trace.SetTag("request_id", log.RequestId)
		trace.SetTag("trace_id", tracer.GetTraceId(ctx))
		// set logger
		since := time.Since(start)
		log.ResponseBody = dumpData(resp)
		log.StatusCode = sc
		log.ExecutionTime = since.Seconds()

		trace.Finish()
		// store to prometheus
		log.Finalize(ctx)
		monitoring.RecordPrometheus(sc, constants.GRPC.String(), info.FullMethod, since)
	}()

	// implement locking logging stdout
	var lock = logger.NewLocker(ctx)
	// set to context with logger.LogKey as a context key
	ctx = context.WithValue(ctx, logger.LogKey, lock)
	logger.SetUsername(ctx, log.Username) // set username to context
	// set tracer request
	trace.SetTag("request.grpc.endpoint", log.Endpoint)
	trace.SetTag("request.grpc.body", log.RequestBody)
	resp, err = handler(ctx, req)
	return
}

func convertErrorToStatusCode(err error) (sc int) {
	switch r := err.(type) {
	case *errs.Error:
		sc = r.StatusCode()
	default:
		c := status.Code(r)

		switch c {
		case codes.FailedPrecondition, codes.InvalidArgument, codes.Unimplemented:
			sc = http.StatusBadRequest
		case codes.Unauthenticated:
			sc = http.StatusUnauthorized
		case codes.PermissionDenied:
			sc = http.StatusForbidden
		case codes.Unknown, codes.NotFound:
			sc = http.StatusNotFound
		case codes.AlreadyExists:
			sc = http.StatusConflict
		case codes.Aborted, codes.Canceled, codes.DeadlineExceeded, codes.Internal, codes.DataLoss:
			sc = http.StatusInternalServerError
		case codes.OutOfRange:
			sc = http.StatusBadGateway
		case codes.Unavailable:
			sc = http.StatusServiceUnavailable
		case codes.ResourceExhausted:
			sc = http.StatusGatewayTimeout
		default:
			sc = http.StatusOK
		}
	}

	if sc < 1 {
		sc = http.StatusInternalServerError
	}

	return
}

func convertStatusCodeToCodes(sc int) (c codes.Code) {
	c = codes.OK

	switch sc {
	case http.StatusBadRequest, http.StatusNotAcceptable, http.StatusGone, http.StatusUnprocessableEntity, http.StatusRequestEntityTooLarge:
		c = codes.InvalidArgument
	case http.StatusUnauthorized:
		c = codes.Unauthenticated
	case http.StatusPaymentRequired, http.StatusPreconditionRequired, http.StatusPreconditionFailed:
		c = codes.FailedPrecondition
	case http.StatusForbidden:
		c = codes.PermissionDenied
	case http.StatusNotFound:
		c = codes.NotFound
	case http.StatusConflict:
		c = codes.AlreadyExists
	case http.StatusTooManyRequests, http.StatusGatewayTimeout:
		c = codes.ResourceExhausted
	case http.StatusInternalServerError:
		c = codes.Aborted
	case http.StatusBadGateway:
		c = codes.OutOfRange
	case http.StatusServiceUnavailable:
		c = codes.Unavailable
	}

	return c
}

func dumpData(data interface{}) string {
	dByte, _ := json.Marshal(data)
	if len(dByte) > env.GetInt("MAX_BODY_SIZE", 1500) {
		return fmt.Sprintf("success response. body response length %d", len(dByte))
	}

	masked := logger.MaskedCredentials(dByte)
	return string(masked)
}
