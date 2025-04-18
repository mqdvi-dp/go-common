package errs

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc/status"
)

const equal = "="
const ands = "&&&"

// NewParseGrpcError parse error from GrpcServer to *errs.Error
func NewParseGrpcError(err error) error {
	s := status.Convert(err)

	msg := s.Message()
	// is there 'equal' in message
	if strings.Contains(msg, equal) {
		msge := strings.Split(msg, equal)
		switch len(msge) {
		case 1:
			msg = strings.TrimSpace(msge[0])
		case 2:
			msg = strings.TrimSpace(msge[1])
		case 3:
			msg = strings.TrimSpace(msge[2])
		}
	}

	messages := strings.Split(msg, ands)
	switch len(messages) {
	case 1:
		message := messages[0]
		if strings.Contains(strings.ToLower(message), context.DeadlineExceeded.Error()) {
			return context.DeadlineExceeded
		} else if strings.Contains(strings.ToLower(message), context.Canceled.Error()) {
			return context.Canceled
		}

		fmt.Println("len(messages) == 1, messages[0]: ", messages[0])
		return NewError(err, CONNECTION_RPC_ERROR.StatusCode(), CONNECTION_RPC_ERROR.Code(), messages[0])
	case 2, 3:
		sysCode, errParse := strconv.Atoi(messages[0])
		if errParse != nil {
			sysCode = CONNECTION_RPC_ERROR.Code()
		}

		codeErr := CodeErr(sysCode)
		errMsg := messages[1]
		moreInfos := []string{}
		if len(messages) > 2 {
			if strings.TrimSpace(messages[2]) != "" {
				moreInfos = []string{messages[2]}
			}
		}

		return NewError(err, codeErr.StatusCode(), codeErr.Code(), errMsg, moreInfos...)
	}

	return NewErrorWithCodeErr(err, CONNECTION_RPC_ERROR)
}
