package errs

import (
	"fmt"
	"strings"
)

// Error instance
type Error struct {
	err        error
	message    string
	moreInfo   string
	statusCode int
	systemCode int
}

// NewError returns a new strandard error
func NewError(err error, statusCode, systemCode int, message string, moreInfos ...string) error {
	return &Error{
		err:        err,
		statusCode: statusCode,
		systemCode: systemCode,
		message:    message,
		moreInfo:   strings.Join(moreInfos, ","),
	}
}

// NewErrorWithCodeErr returns a new error with CodeErr type
func NewErrorWithCodeErr(err error, codeErr CodeErr, moreInfos ...string) error {
	return &Error{
		err:        err,
		statusCode: codeErr.StatusCode(),
		message:    codeErr.Message(),
		systemCode: codeErr.Code(),
		moreInfo:   codeErr.MoreInfo(moreInfos...),
	}
}

func NewErrorWithAdditionalMessage(err error, codeErr CodeErr, message string, moreInfos ...string) error {
	return &Error{
		err:        err,
		statusCode: codeErr.StatusCode(),
		message:    fmt.Sprintf(codeErr.Message(), message),
		systemCode: codeErr.Code(),
		moreInfo:   codeErr.MoreInfo(moreInfos...),
	}
}

// NewErrorWithCodeErrFormatted returns a new error with CodeErr type and formatted more info
func NewErrorWithCodeErrFormatted(err error, codeErr CodeErr, moreInfoFormatted ...interface{}) error {
	return &Error{
		err:        err,
		statusCode: codeErr.StatusCode(),
		message:    codeErr.Message(),
		systemCode: codeErr.Code(),
		moreInfo:   codeErr.MoreInfoFormatted(moreInfoFormatted...),
	}
}

// ParseError returns an instance of Error
func ParseError(err error) *Error {
	switch r := err.(type) {
	case *Error:
		return r
	default:
		return nil
	}
}

func (e *Error) Error() string {
	return e.err.Error()
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) MoreInfo() string {
	return e.moreInfo
}

func (e *Error) SystemCode() int {
	return e.systemCode
}

func (e *Error) StatusCode() int {
	return e.statusCode
}
