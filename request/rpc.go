package request

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type rpcRequest struct {
}

// RpcClient implements rpc client request
type RpcClient interface {
	Do(ctx context.Context, host, fullMethodName string, req, reply interface{}, mds ...metadata.MD) error
}

// NewRpcClient will creates the client of rpc with specific host connection
func NewRpcClient() RpcClient {
	return &rpcRequest{}
}
