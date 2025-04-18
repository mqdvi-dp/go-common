package response

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/logger"
)

type responseBi struct {
	statusCode int
	data       interface{}
}

func (rbi *responseBi) JSON(c *gin.Context) {
	// NOTES: trash way to set http response body
	body, _ := convert.InterfaceToString(rbi.data)
	// check response
	if len(body) > env.GetInt("MAX_BODY_SIZE", 1500) {
		body = "success request"
	}

	c.Set("HTTP_RESPONSE_BODY", body)
	if rbi.statusCode >= http.StatusBadRequest {
		c.AbortWithStatusJSON(rbi.statusCode, rbi.data)
		return
	}

	c.JSON(rbi.statusCode, rbi.data)
}

// Error returns an error response with the given status code and error message
func ErrorBI(ctx context.Context, statusCode int, data interface{}, err error) Response {
	if err != nil {
		// set error into context
		logger.SetErrorMessge(ctx, err)
	}

	return &responseBi{statusCode: statusCode, data: data}
}

// Success returns an success response
func SuccessBI(ctx context.Context, statusCode int, data interface{}) Response {
	return &responseBi{statusCode: statusCode, data: data}
}
