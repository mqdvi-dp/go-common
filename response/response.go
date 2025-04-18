package response

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/errs"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/success"
	"github.com/mqdvi-dp/go-common/types"
)

type response struct {
	statusCode       int                `json:"-"`
	Code             int                `json:"code"`
	Message          string             `json:"message"`
	MoreInfo         string             `json:"more_info"`
	Data             interface{}        `json:"data"`
	Meta             *types.Meta        `json:"meta,omitempty"`
	ErrorValidations []errorValidations `json:"error_validations,omitempty"`
}

type errorValidations struct {
	Field      string `json:"field"`
	Validation string `json:"validation"`
	Message    string `json:"message"`
}

type Response interface {
	JSON(c *gin.Context)
}

func (r *response) JSON(c *gin.Context) {
	// NOTES: trash way to set http response body
	body, _ := convert.InterfaceToString(r)
	// check response
	if len(body) > env.GetInt("MAX_BODY_SIZE", 1500) {
		body = "success request"
	}

	c.Set("HTTP_RESPONSE_BODY", body)
	if r.statusCode >= http.StatusBadRequest {
		c.AbortWithStatusJSON(r.statusCode, r)
		return
	}

	c.JSON(r.statusCode, r)
}

// Error returns an error response with the given status code and error message
func Error(ctx context.Context, err error) Response {
	var resp = &response{
		statusCode: http.StatusInternalServerError,
		Code:       errs.GENERAL_ERROR.Code(),
		Message:    errs.InternalServerError,
	}

	// set error into context
	logger.SetErrorMessge(ctx, err)

	if err != nil {
		switch er := err.(type) {
		case errs.CodeErr:
			resp = &response{
				statusCode: er.StatusCode(),
				Code:       er.Code(),
				Message:    er.Message(),
			}
		case *errs.Error:
			resp = &response{
				statusCode: er.StatusCode(),
				Code:       er.SystemCode(),
				Message:    er.Message(),
				MoreInfo:   er.MoreInfo(),
			}
		case validator.ValidationErrors:
			resp = &response{
				statusCode:       errs.VALIDATION_ERROR.StatusCode(),
				Code:             errs.VALIDATION_ERROR.Code(),
				Message:          errs.VALIDATION_ERROR.Message(),
				ErrorValidations: make([]errorValidations, 0),
			}

			for _, fe := range er {
				v := fe.Tag()
				if fe.Param() != "" {
					v += fmt.Sprintf("=%s", fe.Param())
				}

				resp.ErrorValidations = append(resp.ErrorValidations, errorValidations{
					Field:      fe.Namespace(),
					Validation: v,
					Message:    fe.Error(),
				})
			}
		}
	}

	return resp
}

func ErrorUpsertWithCustomData(ctx context.Context, err error, data interface{}) Response {
	var resp = &response{
		statusCode: http.StatusBadRequest,
		Code:       errs.FAIL_TO_UPSERT_BULK.Code(),
		Message:    errs.FAIL_TO_UPSERT_BULK.Message(),
		Data:       data,
	}
	return resp
}

// Success returns an success response
func Success(ctx context.Context, statusCode int, data interface{}, metas ...*types.Meta) Response {
	var successCode success.CodeSuccess
	switch statusCode {
	case http.StatusOK:
		successCode = success.SUCCESS_GET
	case http.StatusCreated:
		successCode = success.SUCCESS_CREATED
	default:
		successCode = success.SUCCESS_GET
	}
	var meta *types.Meta
	if len(metas) > 0 {
		meta = metas[0]
	}
	return &response{statusCode: statusCode, Code: successCode.Code(), Message: successCode.Message(), Data: data, Meta: meta}
}
