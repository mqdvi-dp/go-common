package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/example/example-service/internal/usecase"
	"github.com/mqdvi-dp/go-common/validation"
)

type httpInstance struct {
	usecase    usecase.Usecase
	validator  validation.Validation
	middleware abstract.Middleware
}

func New(uc usecase.Usecase, mdl abstract.Middleware) abstract.RESTHandler {
	return &httpInstance{
		usecase:    uc,
		validator:  validation.New(),
		middleware: mdl,
	}
}

func (h *httpInstance) Router(r *gin.RouterGroup) {
	v1 := r.Group("/v1")

	v1.POST("/faker/:name", h.getFaker)
}
