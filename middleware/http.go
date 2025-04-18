package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mqdvi-dp/go-common/compare"
	"github.com/mqdvi-dp/go-common/convert"
	"github.com/mqdvi-dp/go-common/helpers"
	"github.com/mqdvi-dp/go-common/zone"

	"github.com/mqdvi-dp/go-common/env"

	"github.com/mqdvi-dp/go-common/constants"
	"github.com/mqdvi-dp/go-common/errs"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/response"
	"github.com/mqdvi-dp/go-common/tracer"
	"github.com/mqdvi-dp/go-common/types"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (m *middleware) HTTPAuth(c *gin.Context) {
	ctx := c.Request.Context()

	trace, ctx := tracer.StartTraceWithContext(ctx, "HTTPMiddleware:HTTPAuth")
	defer trace.Finish()

	auth := c.GetHeader(constants.Authorization)
	// if auth don't exist
	if auth == "" {
		trace.SetError(ErrUnauthorized)
		logger.Log.Error(ctx, fmt.Errorf("header authorization is empty"))

		response.Error(ctx, errs.NewErrorWithCodeErr(fmt.Errorf("value from header Authorization is empty"), errs.UNAUTHORIZED)).JSON(c)
		return
	}

	tokenType, token, err := m.extractAuthType(auth)
	if err != nil {
		trace.SetError(err)
		logger.Log.Error(ctx, err)

		response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.INVALID_AUTH)).JSON(c)
		return
	}

	callFunc, ok := m.authTypeCheckerFunc[tokenType]
	if !ok {
		err = fmt.Errorf("token_type %s not yet implemented", tokenType)
		trace.SetError(err)
		logger.Log.Error(ctx, err)

		response.Error(ctx, errs.NewErrorWithCodeErr(ErrInvalidAuthorization, errs.INVALID_AUTH)).JSON(c)
		return
	}

	// callFunc is calling function Bearer or Basic based on tokenType
	tc, err := callFunc(ctx, token)
	if err != nil {
		trace.SetError(err)

		response.Error(ctx, err).JSON(c)
		return
	}
	// set username into context
	logger.SetUsername(ctx, tc.Username)

	// validate deviceId should be same with token
	if c.GetHeader(constants.ApplicationDevice) != tc.GetDeviceId() {
		err = fmt.Errorf("device id is not same with token. deviceId from request: %s deviceId from token: %s", c.GetHeader(constants.ApplicationDevice), tc.GetDeviceId())
		trace.SetError(err)
		logger.Log.Error(ctx, err)

		response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.ACCESS_DENIED)).JSON(c)
		return
	}

	// set to local context
	// sets user_id into context
	c.Set(constants.UserId, tc.GetUserUuid())
	// sets username into context
	c.Set(constants.Username, tc.GetUsername())
	// sets subscription_id into context
	c.Set(constants.SubscriptionId, tc.GetSubscriptionUuid())
	// set subscriptionStartAt into context
	c.Set(constants.SubscriptionStartAt, tc.GetSubscriptionStartAt())
	// set subscriptionEndAt into context
	c.Set(constants.SubscriptionEndAt, tc.GetSubscriptionEndAt())
	// set is user subscribed
	c.Set(constants.IsSubscribed, tc.GetIsSubscription())
	// set roleKey into context
	c.Set(constants.RoleKey, tc.GetRoleKey())
	// set roleId into context
	c.Set(constants.RoleId, tc.GetRoleId())
	// set statusVerify into context
	c.Set(constants.StatusVerify, tc.GetStatus())
	// set maxLimitedAccess into context
	c.Set(constants.MaxLimitedAccess, tc.GetMaxLimitedAccess())
	// set actual session token (key redis) to context
	c.Set(constants.Authorization, token)
	// set device id into context
	c.Set(constants.DeviceId, tc.GetDeviceId())
	// set app version into context
	c.Set(constants.ApplicationVersion, tc.GetAppVersion())
	// set business name into context
	c.Set(constants.BusinessName, tc.GetBusinessName())
	// set client id into context
	c.Set(constants.ClientId, tc.GetClientId())
	// set client name into context
	c.Set(constants.ClientName, tc.GetClientName())
	// set client key into context
	c.Set(constants.ClientKey, tc.GetClientKey())
	// set client secret into context
	c.Set(constants.ClientSecret, tc.GetClientSecret())
	// set public key into context
	c.Set(constants.PublicKey, tc.GetPublicKey())
	// for mobile identity user for can access money transafer brick bifast
	c.Set(constants.State, tc.GetState())
	c.Set(constants.City, tc.GetCity())
	c.Set(constants.Nik, tc.GetNik())
	c.Set(constants.FullName, tc.GetFullname())
	c.Set(constants.PlaceOfBirth, tc.GetPlaceOfBirth())
	c.Set(constants.DateOfBirth, tc.GetDateOfBirth())
	c.Set(constants.Gender, tc.GetGender())
	c.Set(constants.Address, tc.GetAddress())
	c.Set(constants.RtRw, tc.GetRtRw())
	c.Set(constants.AdministrativeVillage, tc.GetAdministrativeVillage())
	c.Set(constants.District, tc.GetDistrict())
	c.Set(constants.Religion, tc.GetReligion())
	c.Set(constants.MaritalStatus, tc.GetMaritalStatus())
	c.Set(constants.Occupation, tc.GetOccupation())
	c.Set(constants.Nationality, tc.GetNationality())
	c.Set(constants.BloodType, tc.GetBloodType())

	// go to the next handler
	c.Next()
}

func (m *middleware) HTTPApiKey(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		ctx := c.Request.Context()
		trace, ctx := tracer.StartTraceWithContext(ctx, "HTTPMiddleware:HTTPApiKey")
		defer trace.Finish()

		t := c.GetHeader(key)

		// when the key does not exist
		if t == "" {
			err = fmt.Errorf("value from header %s is empty", key)
			trace.SetError(err)
			logger.Log.Error(ctx, err)

			response.Error(ctx, errs.NewErrorWithCodeErr(ErrUnauthorized, errs.UNAUTHORIZED)).JSON(c)
			return
		}

		// when key from params is api_key, we should compare the value with env variables
		if strings.EqualFold(key, constants.AppliacationApiKey) {
			// if same, we can continue to next handler
			if t == env.GetString("API_KEY") {
				c.Next()
				return
			}

			// otherwise, we will send the error
			err = errs.UNAUTHORIZED
			trace.SetError(err)
			logger.Log.Errorf(ctx, "api_key from request is not match with env. from request: %s", t)

			response.Error(ctx, err).JSON(c)
			return
		}

		// callFunc is calling function Bearer or Basic based on tokenType
		tc, err := m.Bearer(ctx, "", t)
		if err != nil {
			trace.SetError(err)

			response.Error(ctx, err).JSON(c)
			return
		}

		// set username into context
		logger.SetUsername(ctx, tc.Username)

		// set to local context
		// sets user_id into context
		c.Set(constants.UserId, tc.GetUserUuid())
		// sets username into context
		c.Set(constants.Username, tc.GetUsername())
		// sets subscription_id into context
		c.Set(constants.SubscriptionId, tc.GetSubscriptionUuid())
		// set subscriptionStartAt into context
		c.Set(constants.SubscriptionStartAt, tc.GetSubscriptionStartAt())
		// set subscriptionEndAt into context
		c.Set(constants.SubscriptionEndAt, tc.GetSubscriptionEndAt())
		// set is user subscribed
		c.Set(constants.IsSubscribed, tc.GetIsSubscription())
		// set roleKey into context
		c.Set(constants.RoleKey, tc.GetRoleKey())
		// set roleId into context
		c.Set(constants.RoleId, tc.GetRoleId())
		// set statusVerify into context
		c.Set(constants.StatusVerify, tc.GetStatus())
		// set maxLimitedAccess into context
		c.Set(constants.MaxLimitedAccess, tc.GetMaxLimitedAccess())
		// set actual session token (key redis) to context
		c.Set(constants.Authorization, t)
		// set device id into context
		c.Set(constants.DeviceId, tc.GetDeviceId())
		// set app version into context
		c.Set(constants.ApplicationVersion, tc.GetAppVersion())
		// set business name into context
		c.Set(constants.BusinessName, tc.GetBusinessName())
		// set client id into context
		c.Set(constants.ClientId, tc.GetClientId())
		// set client name into context
		c.Set(constants.ClientName, tc.GetClientName())
		// set client key into context
		c.Set(constants.ClientKey, tc.GetClientKey())
		// set client secret into context
		c.Set(constants.ClientSecret, tc.GetClientSecret())
		// set public key into context
		c.Set(constants.PublicKey, tc.GetPublicKey())
		// for mobile identity user for can access money transafer brick bifast
		c.Set(constants.State, tc.GetState())
		c.Set(constants.City, tc.GetCity())
		c.Set(constants.Nik, tc.GetNik())
		c.Set(constants.FullName, tc.GetFullname())
		c.Set(constants.PlaceOfBirth, tc.GetPlaceOfBirth())
		c.Set(constants.DateOfBirth, tc.GetDateOfBirth())
		c.Set(constants.Gender, tc.GetGender())
		c.Set(constants.Address, tc.GetAddress())
		c.Set(constants.RtRw, tc.GetRtRw())
		c.Set(constants.AdministrativeVillage, tc.GetAdministrativeVillage())
		c.Set(constants.District, tc.GetDistrict())
		c.Set(constants.Religion, tc.GetReligion())
		c.Set(constants.MaritalStatus, tc.GetMaritalStatus())
		c.Set(constants.Occupation, tc.GetOccupation())
		c.Set(constants.Nationality, tc.GetNationality())
		c.Set(constants.BloodType, tc.GetBloodType())

		// go to the next handler
		c.Next()
	}
}

func (m *middleware) HTTPPermissionACL(resourcePermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var err error

		trace, ctx := tracer.StartTraceWithContext(ctx, "HTTPMiddleware:HTTPPermissionACL")
		defer trace.Finish()

		roleId := c.GetString(constants.RoleId)

		err = m.checkACLPermission(ctx, roleId, resourcePermission)
		if err != nil {
			trace.SetError(err)
			logger.Log.Error(ctx, err)

			response.Error(ctx, err).JSON(c)
			return
		}

		c.Next()
	}
}

func (m *middleware) HTTPRateLimit(limitRequest string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		trace, ctx := tracer.StartTraceWithContext(ctx, "HTTPMiddleware:HTTPRateLimit")
		defer trace.Finish()

		// split limitRequest into maxRequest and duration
		limits := strings.Split(limitRequest, "@")
		if len(limits) != 2 {
			err := ErrInvalidRateLimit
			trace.SetError(err)
			logger.Log.Errorf(ctx, "invalid parameters http_rate_limit: %s", err)

			response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.BAD_REQUEST)).JSON(c)
			return
		}

		// convert string (maxRequest) into integer
		maxRequest, err := strconv.Atoi(limits[0])
		if err != nil {
			trace.SetError(err)
			logger.Log.Error(ctx, err)

			response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.GENERAL_ERROR)).JSON(c)
			return
		}

		// convert string (duration time) into time.Duration
		duration, err := time.ParseDuration(limits[1])
		if err != nil {
			trace.SetError(err)
			logger.Log.Error(ctx, err)

			response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.GENERAL_ERROR)).JSON(c)
			return
		}

		var endpoint = c.Request.URL.Path
		var method = strings.ToUpper(c.Request.Method)

		// get user id from context
		userId := c.GetString(constants.UserId)
		if userId == "" {
			// get unique data from body request
			// for example data unique:
			// - username
			// - email
			// - phone_number
			// - value
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Log.Errorf(ctx, "failed to read body: %s", err)

				c.Next()
				return
			}
			defer c.Request.Body.Close()
			// get body request
			if len(body) > 0 {
				var bodyReq = make(map[string]interface{})
				err = json.Unmarshal(body, &bodyReq)
				if err != nil {
					logger.Log.Errorf(ctx, "failed to unmarshal body: %s", err)
				}

				var val interface{}
				var ok bool
				// hierarchy
				// - username
				// - email
				// - phone_number
				// - value
				if val, ok = bodyReq["username"]; !ok {
					if val, ok = bodyReq["email"]; !ok {
						if val, ok = bodyReq["phone_number"]; !ok {
							val = bodyReq["value"]
						}
					}
				}

				// avoid crash application
				// need to check data types and convert to string
				// if there's no implemented like struct, hashmap, e.t.c
				// then, we just ignore that, and let the userId is empty
				if val != nil {
					switch value := val.(type) {
					case string:
						userId = value
					case int:
						userId = convert.IntToString(int64(value))
					case int32:
						userId = convert.IntToString(int64(value))
					case int64:
						userId = convert.IntToString(value)
					case float32:
						userId = convert.FloatToString(float64(value))
					case float64:
						userId = convert.FloatToString(value)
					}
				}
			}

			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}

		// request to service rate-limit
		err = m.validator.RateLimit(
			ctx, types.RateLimit{
				MaxRequest: maxRequest,
				Duration:   duration,
				UserId:     userId,
				Method:     method,
				Endpoint:   endpoint,
			},
		)
		if err != nil {
			trace.SetError(err)
			logger.Log.Error(ctx, err)

			response.Error(ctx, err).JSON(c)
			return
		}

		c.Next()
	}
}

func (m *middleware) HTTPNotForPublic(c *gin.Context) {
	ctx := c.Request.Context()

	trace, ctx := tracer.StartTraceWithContext(ctx, "HTTPMiddleware:HTTPNotForPublic")
	defer trace.Finish()

	var err error
	roleKey := c.GetString(constants.RoleKey)

	if strings.EqualFold(roleKey, public) {
		err = fmt.Errorf("session is public. this route not for public")
		trace.SetError(err)

		response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.ACCESS_DENIED)).JSON(c)
		return
	}

	c.Next()
}

func (m *middleware) HTTPSignatureValidate(c *gin.Context) {
	ctx := c.Request.Context()

	trace, ctx := tracer.StartTraceWithContext(ctx, "HTTPMiddleware:HTTPSignatureValidate")
	defer trace.Finish()

	if !env.GetBool("FF_HTTP_SIGNATURE_VALIDATE_ENABLED", true) {
		c.Next()
		return
	}

	var err error
	var errValidation = errs.BAD_REQUEST_HEADER

	isMobile := strings.EqualFold(c.GetHeader(constants.ApplicationChannel), "m")
	// * before check the signature, we should check the ip is already block or not?
	// * if blocked, we send the error, otherwise, we can continue to next validation
	var ip = c.GetHeader(constants.ApplicationOriginalIp)
	if isMobile {
		err = m.isIpBlocked(ctx, ip)
		if err != nil {
			trace.SetError(err)

			response.Error(ctx, err).JSON(c)
			return
		}
	}

	xMockLocation := c.GetHeader(constants.ApplicationMockLocation)

	// condition when user using mock-location
	if strings.ToLower(c.GetHeader(constants.ApplicationChannel)) == "m" && env.GetBool("IS_VALIDATE_HEADER", true) {
		var isMockLocation bool
		isMockLocation, err = strconv.ParseBool(xMockLocation)
		if err != nil {
			err = fmt.Errorf("please allow to your location! %s", err)
			logger.Log.Error(ctx, err)
			trace.SetError(err)
			response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.BAD_REQUEST)).JSON(c)
			return
		}

		if isMockLocation {
			err = fmt.Errorf("please trun of mock location")
			logger.Log.Error(ctx, err)
			trace.SetError(err)
			response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.BAD_REQUEST)).JSON(c)
			return
		}
	}

	var message string
	var signatureRequest = c.Request.Header.Get(constants.ApplicationSignature)
	if signatureRequest == "" {
		// set the ip to list blocked
		if isMobile {
			err = m.blockedIp(ctx, ip)
			if err != nil {
				trace.SetError(err)

				response.Error(ctx, err).JSON(c)
				return
			}
		}

		err = fmt.Errorf("signature is empty")
		trace.SetError(err)
		logger.Log.Error(ctx, err)

		response.Error(ctx, errs.NewErrorWithCodeErr(err, errValidation)).JSON(c)
		return
	}

	// * add validation checking deviceId
	// * deviceId from request cannot be 'null'
	reqDeviceId := c.GetHeader(constants.ApplicationDevice) // deviceId from request
	if reqDeviceId == "" {
		// set the ip to list blocked
		if isMobile {
			err = m.blockedIp(ctx, ip)
			if err != nil {
				trace.SetError(err)

				response.Error(ctx, err).JSON(c)
				return
			}
		}

		err = fmt.Errorf("invalid device id")
		logger.Log.Printf(ctx, "device_id from request: %s", reqDeviceId)

		response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.INVALID_DEVICE_ID)).JSON(c)
		return
	}

	// * compare device id in whitelist
	username := c.GetString(constants.Username)
	hashedDeviceId := helpers.SHA1(reqDeviceId)
	isPassed, err := m.checkWhitelist(ctx, hashedDeviceId, username)
	if err != nil {
		trace.SetError(err)

		response.Error(ctx, err).JSON(c)
		return
	}

	if !isPassed {
		err := fmt.Errorf("device_id is not in whitelist")
		logger.Log.Printf(ctx, "new device_id detected: %s", hashedDeviceId)

		response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.NOT_MATCH_DEVICE)).JSON(c)
		return
	}

	// * handling non-auth checking
	auth := c.GetString(constants.Authorization)
	// * the deviceId from request must match with Token
	deviceId := c.GetString(constants.DeviceId) // deviceId from Token
	if reqDeviceId != deviceId && auth != "" {
		// set the ip to list blocked
		if isMobile {
			err = m.blockedIp(ctx, ip)
			if err != nil {
				trace.SetError(err)

				response.Error(ctx, err).JSON(c)
				return
			}
		}

		err = fmt.Errorf("device_id not match with token")
		logger.Log.Printf(ctx, "device_id from request: %s device_id from token: %s", reqDeviceId, deviceId)

		response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.TOKEN_EXPIRED)).JSON(c)
		return
	}

	// get the publicKey
	publicKeyEncoded := c.GetString(constants.PublicKey)
	publicKey, err := getPublicKey(ctx, publicKeyEncoded)
	if err != nil {
		trace.SetError(err)
		logger.Log.Error(ctx, err)

		response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.GENERAL_ERROR)).JSON(c)
		return
	}

	// generate message
	message, err = signatureFormat(ctx, c)
	if err != nil {
		trace.SetError(fmt.Errorf("failed to generate signature message: %s", err))
		logger.Log.Error(ctx, err)

		response.Error(ctx, errs.NewErrorWithCodeErr(err, errValidation)).JSON(c)
		return
	}

	// digest of signature
	digestSignature := hashSHA512([]byte(message))
	// compare
	if !compareSignatureRSA(ctx, publicKey, signatureRequest, digestSignature) {
		// set the ip to list blocked
		if isMobile {
			err = m.blockedIp(ctx, ip)
			if err != nil {
				trace.SetError(err)

				response.Error(ctx, err).JSON(c)
				return
			}
		}

		err = fmt.Errorf("compare is failed")

		trace.SetError(err)

		response.Error(ctx, errs.NewErrorWithCodeErr(err, errValidation)).JSON(c)
		return
	}

	// means compare is success, so request can to next handler
	c.Next()
}

func (m *middleware) HTTPCsrfTokenValidate(c *gin.Context) {
	ctx := c.Request.Context()

	trace, ctx := tracer.StartTraceWithContext(ctx, "HTTPCsrfTokenValidate")
	defer trace.Finish()

	// check the http method
	switch c.Request.Method {
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
	default:
		logger.Log.Print(ctx, "skipped because method is:", c.Request.Method)
		c.Next()
		return
	}

	csrf, err := c.Cookie(constants.CookieCsrfToken)
	if err != nil {
		// if error is ErrNoCookie, we need to check is there cookie with name `uid`?
		// if yes, we will generate cookie new one.
		// this is for handling cookie is expired, we configure expired cookie by env (etcd). the default is 20minutes
		if !errors.Is(err, http.ErrNoCookie) {
			response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.TOKEN_EXPIRED))
			return
		}

		err = m.refreshCsrf(ctx, c)
		if err != nil {
			trace.SetError(err)
			logger.Log.Error(ctx, err)

			response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.TOKEN_EXPIRED))
			return
		}

		// next handler
		c.Next()
		return
	}

	// validate csrf token
	err = m.verifyCsrf(ctx, c.GetString(constants.UserId), csrf)
	if err != nil {
		response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.TOKEN_EXPIRED))
	}

	c.Next()
}

func (m *middleware) HTTPCheckMaintenanceMode(c *gin.Context) {
	req := c.Request
	ctx := req.Context()

	trace, ctx := tracer.StartTraceWithContext(ctx, "HTTP:RestMaintenanceMiddleware")
	defer trace.Finish()

	// get channel
	channel := c.GetHeader(constants.ApplicationChannel)
	// when channels is A (Admin Panel), skip the maintenance mode
	if strings.EqualFold(channel, "a") {
		c.Next()
		return
	}

	// redis client
	rdc := m.redis
	var mt = new(types.Maintenance)
	err := rdc.GetStruct(ctx, mt, constants.PrefixMaintenanceWindow)
	// when error nil, means we've maintenance window
	if err == nil {
		logger.Log.Debugf(ctx, "channel: %s and channels: %s", channel, mt.Channels)
		// checking channel is in maintenance window whitelists
		// if yes, we will skip the maintenance window
		if compare.String(channel, mt.Channels...) {
			c.Next()
			return
		}

		// get user (phone_number) from context
		username := c.GetString(constants.Username)
		// when username is empty, check from body request
		if username == "" {
			// only post and put method
			if req.Method == http.MethodPost || req.Method == http.MethodPut {
				// get body request
				bodyPayload, err := io.ReadAll(req.Body)
				defer req.Body.Close()
				// when error is exists, just skip and go to next handler
				if err != nil {
					c.Next()
					return
				}
				// put back request body
				c.Request.Body = io.NopCloser(bytes.NewReader(bodyPayload))

				// unmarshal using map
				var mp = make(map[string]interface{})
				err = json.Unmarshal(bodyPayload, &mp)
				if err != nil {
					c.Next()
					return
				}
				var keyRequest = env.GetListString("HTTP_BODY_REQUEST_USER", []string{"username", "email", "client_id"}...)
				for _, key := range keyRequest {
					if v, ok := mp[key]; ok {
						val, ok := v.(string)
						if !ok {
							continue
						}
						username = helpers.PhoneNumber(val)
						break
					}
				}
			}
		}
		// skip the maintenance windows when user from context is whitelisted
		if compare.String(username, mt.WhitelistedUsers...) {
			c.Next()
			return
		}

		// get current time
		now := time.Now().In(zone.TzJakarta())
		var moreInfo string
		if !mt.StartTime.IsZero() && !mt.EndTime.IsZero() {
			// when current time is not reach the maintenance window, go to next handler
			if now.Before(mt.StartTime) || now.After(mt.EndTime) {
				c.Next()
				return
			}
			// when end time is not empty
			// give error information to clients
			moreInfo = fmt.Sprintf(env.GetString("MAINTENANCE_MODE_RESPONSE_MORE_INFO", maintenanceMoreInfo), mt.EndTime.Format(time.TimeOnly))
		}
		// return error response
		response.Error(ctx, errs.NewErrorWithCodeErr(fmt.Errorf("maintenance_window active"), errs.MAINTENANCE_MODE, moreInfo)).JSON(c)
		return
	}
	// if error is exists, just skip the middleware and go to next handler
	c.Next()
}

func (m *middleware) HTTPPublicKeyVerification(ctx *gin.Context) {
	const BEARER_SCHEMA = "Bearer "
	authToken := ctx.GetHeader("Authorization")
	authToken = authToken[len(BEARER_SCHEMA):]

	publidKey := []byte(env.GetString("PAYMENT_GATEWAY_PUBLIC_KEY"))

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publidKey)
	if err != nil {
		response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.INVALID_TOKEN))
	}

	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		response.Error(ctx, errs.NewErrorWithCodeErr(err, errs.UNAUTHORIZED))
	}

	_, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		ctx.Next()
	} else {
		response.Error(ctx, errs.INVALID_TOKEN)
	}

}
