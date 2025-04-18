package rest

import (
	"fmt"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mqdvi-dp/go-common/env"
)

type option struct {
	cors         gin.HandlerFunc
	httpPort     string
	rootPath     string
	debugMode    bool
	engineOption func(app *gin.Engine)
}

type OptionFunc func(*option)

func getDefaultOption() option {
	allowOrigins := strings.Split(env.GetString("CORS_ALLOW_ORIGINS", "*"), ",")
	allowMethods := strings.Split(env.GetString("CORS_ALLOW_METHODS", "GET,PUT,POST,DELETE,OPTION"), ",")
	allowHeaders := strings.Split(env.GetString("CORS_ALLOW_HEADERS", "Authorization,Content-Type"), ",")
	allowCredentials := env.GetBool("CORS_ALLOW_CREDENTIALS", true)

	return option{
		httpPort:  fmt.Sprintf(":%d", env.GetInt("SERVICE_HTTP_PORT", 8080)),
		rootPath:  "",
		debugMode: env.GetBool("DEBUG_MODE"),
		cors: cors.New(
			cors.Config{
				AllowCredentials: allowCredentials,
				AllowHeaders:     allowHeaders,
				AllowMethods:     allowMethods,
				AllowOrigins:     allowOrigins,
			},
		),
	}
}

// SetHTTPPort option func
func SetHTTPPort(port int) OptionFunc {
	return func(o *option) {
		o.httpPort = fmt.Sprintf(":%d", port)
	}
}

// SetRootPath option func
func SetRootPath(rootPath string) OptionFunc {
	return func(o *option) {
		o.rootPath = rootPath
	}
}

// SetDebugMode option func
func SetDebugMode(debugMode bool) OptionFunc {
	return func(o *option) {
		o.debugMode = debugMode
	}
}

// SetCorsHandler option func
func SetCorsHandler(cors gin.HandlerFunc) OptionFunc {
	return func(o *option) {
		o.cors = cors
	}
}

// SetEngineOption option func
func SetEngineOption(engine func(app *gin.Engine)) OptionFunc {
	return func(o *option) {
		o.engineOption = engine
	}
}
