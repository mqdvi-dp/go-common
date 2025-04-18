package http

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/response"
	"github.com/mqdvi-dp/go-common/tracer"
)

type requestFaker struct {
	Page   int           `json:"page" validate:"required,gt=10,lt=100"`
	Limit  int           `json:"limit_data" validate:"gt=100"`
	Name   string        `param:"name"`
	O      string        `query:"o" validate:"required,hexcolor|rgb|rgba"`
	Nested *requestFaker `json:"fakers"`
}

func (h *httpInstance) getFaker(c *gin.Context) {
	ctx := c.Request.Context()
	trace, ctx := tracer.StartTraceWithContext(ctx, "HttpHandler:GetFaker")
	defer trace.Finish()

	var req = new(requestFaker)
	err := h.validator.BindAndValidateWithContext(ctx, c, req)
	if err != nil {
		trace.SetError(err)

		response.Error(ctx, err).JSON(c)
		return
	}

	rr, _ := convert.InterfaceToString(req)
	tracer.Log(ctx, "requestFaker", fmt.Sprintf("requestFaker: %s", rr))

	resp, err := h.usecase.GetFaker(ctx)
	if err != nil {
		trace.SetError(err)
		response.Error(ctx, err).JSON(c)
		return
	}

	response.Success(ctx, http.StatusOK, resp).JSON(c)
}
