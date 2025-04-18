package rest

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/factory"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/zone"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type rest struct {
	serverEngine *gin.Engine
	httpServer   *http.Server
	service      factory.ServiceFactory
	opt          option
	tz           *time.Location
}

// New create new rest server
func New(service factory.ServiceFactory, opts ...OptionFunc) factory.AppServerFactory {
	tz := zone.TzJakarta()
	gin.SetMode(gin.ReleaseMode)
	srv := &rest{
		serverEngine: gin.New(),
		service:      service,
		opt:          getDefaultOption(),
		tz:           tz,
	}

	for _, opt := range opts {
		opt(&srv.opt)
	}

	if srv.opt.engineOption != nil {
		srv.opt.engineOption(srv.serverEngine)
	}

	// add cors
	srv.serverEngine.Use(srv.opt.cors)

	// start for health-check
	h := healthCheck(service)
	pg := srv.serverEngine.Group("/ping")
	pg.GET("", gin.WrapH(h.Handler()))

	mg := srv.serverEngine.Group("/metrics")
	mg.GET("", gin.WrapH(promhttp.Handler()))

	// root path for HTTP handler
	rootPath := srv.serverEngine.Group(srv.opt.rootPath, srv.restTracerMiddleware, srv.validateAllowedHeaders)
	// register rest handler on each service
	if r := service.RESTHandler(); r != nil {
		// register service
		r.Router(rootPath)
	}

	// print all routes
	for _, route := range srv.serverEngine.Routes() {
		if strings.EqualFold(route.Method, http.MethodHead) {
			continue
		}

		logger.Cyan(fmt.Sprintf(`[REST-SERVER-ROUTE] (method): %-8s (route): %s (handler): %s`, `"`+route.Method+`"`, `"`+route.Path+`"`, `"`+route.Handler+`"`))
	}

	srv.httpServer = &http.Server{
		Addr:    env.GetString("SERVICE_HTTP_HOST", "0.0.0.0") + srv.opt.httpPort,
		Handler: srv.serverEngine,
	}
	return srv
}

func (r *rest) Serve() {
	logger.CyanBold(fmt.Sprintf("⇨ HTTP Server run at port [::]%s", r.opt.httpPort))

	err := r.httpServer.ListenAndServe()

	switch e := err.(type) {
	case *net.OpError:
		panic(fmt.Errorf("rest server: %s", e))
	}
}

func (r *rest) Shutdown(ctx context.Context) {
	defer fmt.Printf("\x1b[33;1m%s Stopping HTTP server:\x1b[0m \x1b[32;1mSUCCESS\x1b[0m\n", r.service.Name())

	_ = r.httpServer.Shutdown(ctx)
}

func (r *rest) Name() string {
	return string(constants.REST)
}
