package abstract

import (
	"github.com/gin-gonic/gin"
	"github.com/mqdvi-dp/go-common/types"

	"google.golang.org/grpc"
)

// RESTHandler delivery factory for REST Handler
// Default: gin-gonic/gin rest framework
type RESTHandler interface {
	Router(r *gin.RouterGroup)
}

// GRPCHandler delivery factory for gRPC handler
type GRPCHandler interface {
	RegisterGrpcServer(srv *grpc.Server)
}

// WorkerHandler delivery factory for all worker handler
type WorkerHandler interface {
	Register(group *types.WorkerHandlerGroup)
}
