package rpc

import (
	"github.com/mqdvi-dp/go-common/abstract"
	"google.golang.org/grpc"
)

type rpcHandler struct{}

func New() abstract.GRPCHandler {
	return &rpcHandler{}
}

func (r *rpcHandler) RegisterGrpcServer(srv *grpc.Server) {
}
