package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/helpers/paseto"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/types"
)

func (m *middleware) verifyCsrf(ctx context.Context, userUuid string, value string) error {
	csrfToken, err := m.decryptCsrfCookie(ctx, value)
	if err != nil {
		tracer.SetError(ctx, err)

		return err
	}

	// validate user_uuid from auth is match with csrf
	if csrfToken.UserUuid == userUuid {
		return nil
	}

	// sending the error
	err = fmt.Errorf("user_uuid from csrf not match with auth. csrf_userid: %s user_uuid token: %s", csrfToken.UserUuid, userUuid)

	tracer.SetError(ctx, err)
	return err
}

func (m *middleware) refreshCsrf(ctx context.Context, c *gin.Context) error {
	// get the cookie with name `uid` first, after that we generate new one
	uid, err := c.Cookie(constants.CookieUid)
	if err != nil {
		// when cookie uid not found
		// users should generate from login
		err = fmt.Errorf("failed to get cookie name: %s and the error is: %s", constants.CookieUid, err)
		return err
	}

	csrfToken, err := m.decryptCsrfCookie(ctx, uid)
	if err != nil {
		return err
	}

	// generate types with paseto
	ct, err := paseto.GeneratePaseto(ctx, env.GetString("KLIKOO_AUTH_PRIVATE_KEY"), csrfToken)
	if err != nil {
		return err
	}

	// set cookie
	maxAge := env.GetInt("COOKIE_CSRF_DURATION", 1200) // default is 20minutes
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(constants.CookieCsrfToken, ct, maxAge, "/", env.GetString("COOKIE_DOMAIN"), true, true)
	return nil
}

func (m *middleware) decryptCsrfCookie(ctx context.Context, data string) (*types.CsrfToken, error) {
	payload, _, err := paseto.VerifyPaseto(ctx, env.GetString("KLIKOO_AUTH_PUBLIC_KEY"), data)
	if err != nil {
		err = fmt.Errorf("cookie value is invalid: %s", err)
		return nil, err
	}

	// type assertion, type must be *types.CsrfToken
	csrfToken, ok := payload.(*types.CsrfToken)
	if !ok {
		err = fmt.Errorf("types of cookie uid decrpyted is not *types.CsrfToken")
		return nil, err
	}

	return csrfToken, nil
}
