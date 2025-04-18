package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/errs"
)

func checkIp(c *gin.Context) error {
	whitelists := env.GetListString("HTTP_WHITELIST_IP_SKIPPED_SIGNATURE")

	ip := c.GetHeader(constants.ApplicationOriginalIp)
	for _, w := range whitelists {
		if ip == w {
			return nil
		}
	}

	return errs.NewErrorWithCodeErr(ErrIpNotWhitelisted, errs.IP_BLOCKED)
}
