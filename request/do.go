package request

import (
	"context"
	"fmt"

	"github.com/mqdvi-dp/go-common/factory/server/rpc"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func (r *rpcRequest) Do(ctx context.Context, host, fullMethodName string, req, reply interface{}, mds ...metadata.MD) error {
	trace, ctx := tracer.StartTraceWithContext(ctx, "RpcRequest:Do")
	defer trace.Finish()

	// validate value of request body
	err := validatePtr(req)
	if err != nil {
		trace.SetError(err)
		return err
	}

	// validate value of reply (destination of response data types)
	err = validatePtr(reply)
	if err != nil {
		trace.SetError(err)
		return err
	}

	// check the host
	if host == "" {
		return fmt.Errorf("host is empty")
	}

	// check the method rpc name
	if fullMethodName == "" {
		return fmt.Errorf("rpc fullMethodName is empty")
	}

	// we create the headers and the default headers is request_id for tracing grpc request from each services
	md := metadata.Pairs("request_id", logger.GetRequestId(ctx))
	if len(mds) > 0 {
		for _, mtds := range mds {
			for key, vals := range mtds {
				md[key] = vals
			}
		}
	}
	// add into metadata grpc
	ctx = metadata.NewOutgoingContext(ctx, md)
	// creates a interceptor grpc
	intercept := rpc.NewInterceptor(host)
	// init grpc connection
	conn, err := grpc.Dial(
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: backoff.DefaultConfig}),
		grpc.WithChainUnaryInterceptor(intercept.ChainUnaryClient(intercept.UnaryClientTracerInterceptor)),
	)
	if err != nil {
		trace.SetError(err)
		return err
	}
	// close the connection
	defer conn.Close()

	// request to client
	err = conn.Invoke(ctx, fullMethodName, req, reply)
	if err != nil {
		trace.SetError(err)
		return err
	}

	// request is done
	return nil
}
