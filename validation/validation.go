package validation

import (
	"context"
	"reflect"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mqdvi-dp/go-common/logger"
)

type CustomValidationFunc func(*validation)

type validation struct {
	validate         *validator.Validate
	customValidation map[string]customValidation
}

type customValidation struct {
	callValidationEvenIfNull bool
	fn                       validator.Func
}

// SetCustomValidation to sets the custom validation
func SetCustomValidation(tag string, fn validator.Func, isCallFuncValidation ...bool) CustomValidationFunc {
	return func(v *validation) {
		if len(v.customValidation) < 1 || v.customValidation == nil {
			v.customValidation = map[string]customValidation{}
		}

		// is validation already registered?
		if cv, ok := v.customValidation[tag]; ok {
			funcName := runtime.FuncForPC(reflect.ValueOf(cv.fn).Pointer()).Name()
			// if already registered, will return an fatal error
			logger.Log.Fatalf("validation tag %s already registered with funcName %s", tag, funcName)
			return
		}

		isCallFunc := false
		if len(isCallFuncValidation) > 0 {
			isCallFunc = isCallFuncValidation[0]
		}

		v.customValidation[tag] = customValidation{
			callValidationEvenIfNull: isCallFunc,
			fn:                       fn,
		}
	}
}

// New initiates validation
func New(opts ...CustomValidationFunc) Validation {
	v := &validation{
		validate:         validator.New(),
		customValidation: make(map[string]customValidation),
	}

	for _, opt := range opts {
		opt(v)
	}

	for tag, val := range v.customValidation {
		v.validate.RegisterValidation(tag, val.fn, val.callValidationEvenIfNull)
	}

	v.validate.RegisterTagNameFunc(validateJsonTag)
	v.validate.RegisterTagNameFunc(validateFormDataTag)
	v.validate.RegisterTagNameFunc(validateParamTag)
	v.validate.RegisterTagNameFunc(validateQueryTag)
	v.validate.RegisterTagNameFunc(validateHeaderTag)

	return v
}

type Validation interface {
	// BindJSONAndValidate bind JSON data and validate the struct
	BindJSONAndValidate(ctx context.Context, bytes []byte, dest interface{}) error
	// BindAndValidate bind data and validate the struct
	BindAndValidate(c *gin.Context, dest interface{}) error
	// BindAndValidateWithContext bind data and validate the struct with context
	BindAndValidateWithContext(ctx context.Context, c *gin.Context, dest interface{}) error
	// Validate the struct
	Validate(ctx context.Context, dest interface{}) error
}

func validateJsonTag(field reflect.StructField) string {
	name := strings.SplitN(field.Tag.Get(jsonTag), ",", 2)[0]
	if name == "-" {
		return name
	}

	return name
}

func validateFormDataTag(field reflect.StructField) string {
	name := field.Tag.Get(form)
	if name == "-" {
		return ""
	}

	return name
}

func validateParamTag(field reflect.StructField) string {
	name := field.Tag.Get(param)
	if name == "-" {
		return ""
	}

	return name
}

func validateQueryTag(field reflect.StructField) string {
	name := field.Tag.Get(query)
	if name == "-" {
		return ""
	}

	return name
}

func validateHeaderTag(field reflect.StructField) string {
	name := field.Tag.Get(header)
	if name == "-" {
		return ""
	}

	return name
}
