package abstract

import (
	"context"

	"github.com/mqdvi-dp/go-common/types"

	"github.com/gin-gonic/gin"
)

// Middleware abstractions
type Middleware interface {
	// Basic authorization
	Basic(ctx context.Context, auth string) error

	// Bearer authorization
	Bearer(ctx context.Context, tokenType, token string) (types.TokenClaim, error)

	// HTTPApiKey header authentication with declared keys
	HTTPApiKey(key string) gin.HandlerFunc

	// HTTPAuth is for checking authorization
	HTTPAuth(c *gin.Context)

	// HTTPPermissionACL checking role acl permission
	HTTPPermissionACL(resourcePermission string) gin.HandlerFunc

	// HTTPRateLimit checking limitation user to access api
	HTTPRateLimit(limitRequest string) gin.HandlerFunc

	// HTTPNotForPublic checking roleKey is public or not
	// if public, then given an error access denied
	HTTPNotForPublic(c *gin.Context)

	// HTTPSignatureValidate to validate the application request
	// this means only application can access to our API
	HTTPSignatureValidate(c *gin.Context)

	// HTTPCsrfTokenValidate to validate the csrf_token
	HTTPCsrfTokenValidate(c *gin.Context)

	// HTTPCheckMaintenanceMode checking maintenance windows at our systems
	HTTPCheckMaintenanceMode(c *gin.Context)

	HTTPPublicKeyVerification(ctx *gin.Context)
}

// AuthenticationValidator abstraction interface for validating token authorization
type AuthenticationValidator interface {
	// ValidateToken validation token
	ValidateToken(ctx context.Context, tokenType, token string) (types.TokenClaim, error)

	// CheckPermission checking permission with level Role, Resource and Permission
	CheckPermission(ctx context.Context, roleId, resource string) error

	// RateLimit checking is user has limitation to access this resource or not
	// if user has limited, it'll send response error status_code: 429 (To Many Request)
	RateLimit(ctx context.Context, rl types.RateLimit) error
}
