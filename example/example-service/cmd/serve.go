package cmd

import (
	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/config"
	"github.com/mqdvi-dp/go-common/example/example-service/internal/handler/http"
	"github.com/mqdvi-dp/go-common/example/example-service/internal/handler/rpc"
	"github.com/mqdvi-dp/go-common/example/example-service/internal/repository/api"
	"github.com/mqdvi-dp/go-common/example/example-service/internal/usecase"
	"github.com/mqdvi-dp/go-common/factory/server"
	"github.com/mqdvi-dp/go-common/factory/server/rest"
)

func Serve(cfg config.Config) *server.Server {
	deps := dependencies(cfg)

	apiRepo := api.New()
	// cacheRepo := cache.New(deps.GetRedisDatabase().Client())
	uc := usecase.New(apiRepo, nil)

	svc := server.NewApplicationService(
		server.SetConfiguration(cfg),
		runHttp(uc, deps.GetMiddleware()),
		runGrpc(uc),
		// runNsq(uc),
	)
	return server.New(svc)
}

func runHttp(uc usecase.Usecase, mdl abstract.Middleware) server.ServiceFunc {
	return server.SetRestHandler(
		http.New(uc, mdl),
		rest.SetHTTPPort(9999),
	)
}

func runGrpc(uc usecase.Usecase) server.ServiceFunc {
	return server.SetGrpcHandler(rpc.New())
}
