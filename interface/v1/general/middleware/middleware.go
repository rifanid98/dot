package middleware

import (
	"dot/pkg/helper"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"

	"dot/app/v1/deps"
	"dot/config"
	"dot/core"
	"dot/core/v1/port/auth"
	"dot/interface/v1/general/common"
	"dot/pkg/util"

	echoLog "github.com/labstack/gommon/log"
)

var log = util.NewLogger()

func RequiredCAHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c, err := ParseCAHeader(c)
		if err != nil {
			return err
		}

		return next(c)
	}
}

func InternalContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ic := ParseInternalContext(c)
		c.Set("ic", ic)
		return next(c)
	}
}

func ParseCAHeader(c echo.Context) (echo.Context, error) {
	c, err := ParseClientIdHeader(c)
	if err != nil {
		return c, err
	}

	c, err = ParseAppIdHeader(c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func ParseClientIdHeader(c echo.Context) (echo.Context, error) {
	ic := ParseInternalContext(c)

	clientId := c.Request().Header.Get(common.CLIENT_ID)
	if clientId == "" {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Message = "X-CLIENT-ID is required in header"
		return nil, c.JSON(
			http.StatusBadRequest,
			common.NewResponse(res, nil),
		)
	}

	c.Set("client_id", clientId)
	c.Set("ic", ic.AppendData(map[string]any{
		"client_id": clientId,
	}))

	return c, nil
}

func ParseInternalContext(c echo.Context) *core.InternalContext {
	ic := core.NewInternalContext(c.Get("tracker_id").(string))
	ic_ := c.Get("ic")
	if ic_ != nil {
		ic = ic_.(*core.InternalContext)
	}

	return ic
}

func ParseAppIdHeader(c echo.Context) (echo.Context, error) {
	ic := ParseInternalContext(c)

	appId := c.Request().Header.Get(common.APP_ID)
	if appId == "" {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Message = "X-APP-ID is required in header"
		return nil, c.JSON(
			http.StatusBadRequest,
			common.NewResponse(res, nil),
		)
	}

	c.Set("ic", ic.AppendData(map[string]any{
		"app_id": appId,
	}))

	return c, nil
}

//func RequiredCAHeader(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		clientId := c.Request().Header.Get(common.CLIENT_ID)
//		if clientId == "" {
//			res := common.ResultMap[core.BAD_REQUEST]
//			res.Message = "X-CLIENT-ID is required in header"
//			return c.JSON(
//				http.StatusBadRequest,
//				common.NewResponse(res, nil),
//			)
//		}
//
//		appId := c.Request().Header.Get(common.APP_ID)
//		if appId == "" {
//			res := common.ResultMap[core.BAD_REQUEST]
//			res.Message = "X-APP-ID is required in header"
//			return c.JSON(
//				http.StatusBadRequest,
//				common.NewResponse(res, nil),
//			)
//		}
//
//		return next(c)
//	}
//}

func RequiredCAAHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c, err := ParseCAHeader(c)
		if err != nil {
			return err
		}

		c, err = ParseAuthorizationHeader(c)
		if err != nil {
			return err
		}

		return next(c)
	}
}

func ParseAuthorizationHeader(c echo.Context) (echo.Context, error) {
	authorization := c.Request().Header.Get(common.AUTHORIZATION)
	if authorization == "" {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Message = "Authorization is required in header"
		return c, c.JSON(http.StatusBadRequest, common.NewResponse(res, nil))
	}

	return c, nil
}

//func RequiredCAAHeader(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		clientId := c.Request().Header.Get(common.CLIENT_ID)
//		if clientId == "" {
//			res := common.ResultMap[core.BAD_REQUEST]
//			res.Message = "X-CLIENT-ID is required in header"
//			return c.JSON(
//				http.StatusBadRequest,
//				common.NewResponse(res, nil),
//			)
//		}
//
//		appId := c.Request().Header.Get(common.APP_ID)
//		if appId == "" {
//			res := common.ResultMap[core.BAD_REQUEST]
//			res.Message = "X-APP-ID is required in header"
//			return c.JSON(
//				http.StatusBadRequest,
//				common.NewResponse(res, nil),
//			)
//		}
//
//		authorization := c.Request().Header.Get(common.AUTHORIZATION)
//		if authorization == "" {
//			res := common.ResultMap[core.BAD_REQUEST]
//			res.Message = "Authorization is required in header"
//			return c.JSON(
//				http.StatusBadRequest,
//				common.NewResponse(res, nil),
//			)
//		}
//
//		return next(c)
//	}
//}

func ParseApiKeyHeader(c echo.Context, cfg *config.AppConfig) (echo.Context, error) {
	skipped := apiSkipper(c)
	if skipped {
		return c, nil
	}

	apiKey := c.Request().Header.Get(common.API_KEY)
	if apiKey == "" {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Message = "X-API-KEY is required in header"
		return nil, c.JSON(
			http.StatusBadRequest,
			common.NewResponse(res, nil),
		)
	}

	if apiKey != cfg.ApiKey {
		res := common.ResultMap[core.UNAUTHORIZED]
		return nil, c.JSON(
			http.StatusUnauthorized,
			common.NewResponse(res, nil),
		)
	}

	c.Set("api_key", apiKey)

	return c, nil
}

func RequiredCAPHeader(deps deps.IDependency) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c, err := ParseCAHeader(c)
			if err != nil {
				return err
			}

			c, err = ParseApiKeyHeader(c, deps.GetBase().Cfg)
			if err != nil {
				return err
			}

			return next(c)
		}
	}
}

//func RequiredCAPHeader(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		clientId := c.Request().Header.Get(common.CLIENT_ID)
//		if clientId == "" {
//			res := common.ResultMap[core.BAD_REQUEST]
//			res.Message = "X-CLIENT-ID is required in header"
//			return c.JSON(
//				http.StatusBadRequest,
//				common.NewResponse(res, nil),
//			)
//		}
//
//		appId := c.Request().Header.Get(common.APP_ID)
//		if appId == "" {
//			res := common.ResultMap[core.BAD_REQUEST]
//			res.Message = "X-APP-ID is required in header"
//			return c.JSON(
//				http.StatusBadRequest,
//				common.NewResponse(res, nil),
//			)
//		}
//
//		apiKey := c.Request().Header.Get(common.API_KEY)
//		if apiKey == "" {
//			res := common.ResultMap[core.BAD_REQUEST]
//			res.Message = "X-API-KEY is required in header"
//			return c.JSON(
//				http.StatusBadRequest,
//				common.NewResponse(res, nil),
//			)
//		}
//
//		return next(c)
//	}
//}

func RequiredApiKeyHeader(deps deps.IDependency) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c, err := ParseApiKeyHeader(c, deps.GetBase().Cfg)
			if err != nil {
				return err
			}

			return next(c)
		}
	}
}

//func RequiredApiKeyHeader(next echo.HandlerFunc) echo.HandlerFunc {
//	return func(c echo.Context) error {
//		apiKey := c.Request().Header.Get(common.API_KEY)
//		if apiKey == "" {
//			res := common.ResultMap[core.BAD_REQUEST]
//			res.Message = "X-API-KEY is required in header"
//			return c.JSON(
//				http.StatusBadRequest,
//				common.NewResponse(res, nil),
//			)
//		}
//
//		return next(c)
//	}
//}

func JwtTokenMiddleware(uc auth.AuthUsecase, secrets ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			trackerID, _ := c.Get("tracker_id").(string)
			ic := core.NewInternalContext(trackerID)
			rawToken := c.Request().Header.Get("Authorization")

			tokenString, cerr := getTokenString(rawToken)
			if cerr != nil {
				return cerr
			}

			token, cerr := validateTokens(ic, tokenString, secrets)
			if cerr != nil {
				res := common.ResultMap[cerr.Code]
				log.Error(ic.ToContext(), "unauthorized invalid token "+res.Message, res.Errors)
				return c.JSON(
					res.StatusCode,
					common.NewResponse(res, nil),
				)
			}

			claim := token.Claims.(jwt.MapClaims)
			cerr = uc.IsActiveToken(ic, claim["id"].(string), tokenString)
			if cerr != nil {
				res := common.ResultMap[cerr.Code]
				return c.JSON(
					res.StatusCode,
					common.NewResponse(res, nil),
				)
			}

			c.Set("token", token)
			return next(c)
		}
	}
}

func JwtAccessTokenMiddleware(uc auth.AuthUsecase, secrets ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ic := ParseInternalContext(c)
			rawToken := c.Request().Header.Get("Authorization")

			tokenString, cerr := getTokenString(rawToken)
			if cerr != nil {
				res := common.ResultMap[cerr.Code]
				log.Error(ic.ToContext(), "jwt token is empty "+res.Message, res.Errors)
				return c.JSON(
					res.StatusCode,
					common.NewResponse(res, nil),
				)
			}

			token, cerr := validateTokens(ic, tokenString, secrets)
			if cerr != nil {
				res := common.ResultMap[cerr.Code]
				log.Error(ic.ToContext(), "unauthorized invalid token "+res.Message, res.Errors)
				return c.JSON(
					res.StatusCode,
					common.NewResponse(res, nil),
				)
			}

			claim := token.Claims.(jwt.MapClaims)
			if claim["is_refresh"] == true {
				res := common.ResultMap[core.INVALID_JWT_TOKEN]
				return c.JSON(
					res.StatusCode,
					common.NewResponse(res, nil),
				)
			}

			cerr = uc.IsActiveToken(ic, claim["id"].(string), tokenString)
			if cerr != nil {
				res := common.ResultMap[cerr.Code]
				return c.JSON(
					res.StatusCode,
					common.NewResponse(res, nil),
				)
			}

			c.Set("token", token)
			c.Set("token_string", tokenString)

			return next(c)
		}
	}
}

func JwtRefreshTokenMiddleware(uc auth.AuthUsecase, secrets ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			trackerID, _ := c.Get("tracker_id").(string)
			ic := core.NewInternalContext(trackerID)
			rawToken := c.Request().Header.Get("Authorization")

			tokenString, cerr := getTokenString(rawToken)
			if cerr != nil {
				return cerr
			}

			token, cerr := validateTokens(ic, tokenString, secrets)
			if cerr != nil {
				res := common.ResultMap[cerr.Code]
				log.Error(ic.ToContext(), "unauthorized invalid token "+res.Message, res.Errors)
				return c.JSON(
					res.StatusCode,
					common.NewResponse(res, nil),
				)
			}

			claim := token.Claims.(jwt.MapClaims)
			if claim["is_refresh"] == false {
				res := common.ResultMap[core.INVALID_JWT_REFRESH_TOKEN]
				return c.JSON(
					res.StatusCode,
					common.NewResponse(res, nil),
				)
			}

			cerr = uc.IsActiveToken(ic, claim["id"].(string), tokenString)
			if cerr != nil {
				res := common.ResultMap[cerr.Code]
				return c.JSON(
					res.StatusCode,
					common.NewResponse(res, nil),
				)
			}

			c.Set("token", token)
			c.Set("token_string", tokenString)
			return next(c)
		}
	}
}

func getTokenString(rawToken string) (string, *core.CustomError) {
	if len(strings.Split(rawToken, " ")) == 2 {
		return strings.Split(rawToken, " ")[1], nil
	}

	if len(rawToken) > 1 {
		return "", &core.CustomError{
			Code: core.INVALID_JWT_TOKEN,
		}
	}

	return "", &core.CustomError{
		Code: core.JWT_TOKEN_REQUIRED,
	}
}

func JwtAccessTokenParseMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	// note: parse jwt claims without validating the token
	return func(c echo.Context) error {
		ic := ParseInternalContext(c)

		tokenString := c.Request().Header.Get("Authorization")
		splitedToken := strings.Split(tokenString, " ")

		if len(splitedToken) < 2 {
			if len(tokenString) == 0 {
				res := common.ResultMap[core.JWT_TOKEN_REQUIRED]
				return c.JSON(res.StatusCode, common.NewResponse(res, nil))
			}
			res := common.ResultMap[core.INVALID_JWT_TOKEN]
			return c.JSON(res.StatusCode, common.NewResponse(res, nil))
		}

		token, _ := jwt.Parse(splitedToken[1], nil)
		if token == nil {
			res := common.ResultMap[core.INVALID_JWT_TOKEN]
			return c.JSON(res.StatusCode, common.NewResponse(res, nil))
		}

		claimsMap := token.Claims.(jwt.MapClaims)
		accountId, ok := claimsMap["id"]
		if !ok {
			res := common.ResultMap[core.INVALID_JWT_TOKEN]
			return c.JSON(res.StatusCode, common.NewResponse(res, nil))
		}
		verified, ok := claimsMap["verified"]
		if !ok {
			res := common.ResultMap[core.INVALID_JWT_TOKEN]
			return c.JSON(res.StatusCode, common.NewResponse(res, nil))
		}

		claims := common.JwtClaims{
			Id:       accountId.(string),
			Verified: int(helper.DataToInt(verified)),
		}

		c.Set("claims", claims)
		ic.AppendData(map[string]any{
			"claims": claims,
			"id":     accountId.(string),
		})
		return next(c)
	}
}

func validateToken(ic *core.InternalContext, tokenString string, jwtSecret string) (*jwt.Token, *core.CustomError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, valid := token.Method.(*jwt.SigningMethodHMAC); !valid || method != jwt.SigningMethodHS256 {
			return nil, &core.CustomError{
				Code: core.INVALID_JWT_TOKEN,
			}
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		log.Error(ic.ToContext(), "failed token.Method.(*jwt.SigningMethodHMAC)", err.Error())
		return nil, &core.CustomError{
			Code:    core.INVALID_JWT_TOKEN,
			Message: err.Error(),
		}
	}

	return token, nil
}

func validateTokens(ic *core.InternalContext, tokenString string, secrets []string) (*jwt.Token, *core.CustomError) {
	var token *jwt.Token
	var cerr *core.CustomError
	var valid = false

	for _, secret := range secrets {
		token, cerr = validateToken(ic, tokenString, secret)
		if cerr == nil {
			valid = true
			break

		}
	}

	if !valid {
		return nil, cerr
	}

	return token, nil
}

func ServiceTrackerID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("tracker_id", uuid.New().String())
		return next(c)
	}
}

func ServiceRequestTime(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Request().Header.Set("X-App-RequestTime", time.Now().Format(time.RFC3339))
		return next(c)
	}
}

func Recover() echo.MiddlewareFunc {
	return middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:           nil,
		StackSize:         1 << 10, // 1 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogLevel:          echoLog.ERROR,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			log.Error(c.Request().Context(), "[PANIC RECOVER]", err, string(stack))
			return err
		},
	})
}

func CORS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	})
}

func ApiKeyAuth(cfg *config.AppConfig) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper:   apiSkipper,
		KeyLookup: "header:X-API-KEY",
		Validator: apiKeyAuth(cfg),
	})
}

func apiKeyAuth(cfg *config.AppConfig) middleware.KeyAuthValidator {
	return func(auth string, c echo.Context) (bool, error) {
		valid := auth == cfg.ApiKey
		if valid {
			return true, nil
		}
		return false, errors.New("unauthorized")
	}
}

func apiSkipper(c echo.Context) bool {
	urls := []string{
		"login",
		"swagger",
		"logout",
		"password/forgot",
		"password/reset",
		"password/change",
	}
	for _, url := range urls {
		if strings.Contains(c.Request().URL.Path, url) {
			return true
		}
	}
	return false
}

//func apiSkipper(cfg *config.AppConfig) func(c echo.Context) bool {
//	return func(c echo.Context) bool {
//		urls := []string{
//			"login",
//			"swagger",
//			"logout",
//			"password/forgot",
//			"password/reset",
//			"password/change",
//		}
//		for _, url := range urls {
//			if strings.Contains(c.Request().URL.Path, url) {
//				return true
//			}
//		}
//		return false
//	}
//}
