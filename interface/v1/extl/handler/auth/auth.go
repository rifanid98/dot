package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"dot/core"
	"dot/core/v1/port/auth"
	"dot/interface/v1/general/common"
	"dot/pkg/util"
)

var log = util.NewLogger()

type Handler struct {
	AuthUsecase auth.AuthUsecase
}

func New(service auth.AuthUsecase) *Handler {
	return &Handler{
		AuthUsecase: service,
	}
}

// Register 	godoc
// @Summary		Register Account.
// @Description	Register Account.
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param  		Body 			body		Register 	true	"body"
// @Failure		400
// @Failure		401
// @Failure		422
// @Failure		500
// @Router	/api/v1/auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	ic := c.Get("ic").(*core.InternalContext)

	request := new(Register)
	if err := c.Bind(request); err != nil {
		res := common.HandleBind(err)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	if err := c.Validate(request); err != nil {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Errors = util.GetValidatorMessage(err)
		return c.JSON(
			http.StatusBadRequest,
			common.NewResponse(res, nil),
		)
	}

	register, cerr := h.AuthUsecase.Register(ic, request.Account())
	if cerr != nil {
		res := common.HandleError(cerr)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	return c.JSON(
		http.StatusOK,
		common.NewResponse(common.ResultMap[core.CREATED], register),
	)
}

// Login godoc
// @Summary 	Get Token.
// @Description	Get Token.
// @Tags 		Auth
// @Accept	 	json
// @Produce 	json
// @Param  		Body 			body		Login 	true	"body"
// @Success		200				{object}	TokenResponse
// @Failure 	400
// @Failure 	401
// @Failure 	500
// @Router /api/v1/auth/login [post]
func (handler *Handler) Login(c echo.Context) error {
	ic := c.Get("ic").(*core.InternalContext)

	request := new(Login)
	if err := c.Bind(request); err != nil {
		res := common.ResultMap[core.BAD_REQUEST]
		errors := common.ErrorBindBodyRequest(err)
		res.Errors = []string{errors.Field + " " + errors.Format}
		return c.JSON(
			http.StatusBadRequest,
			common.NewResponse(res, nil),
		)
	}

	if err := c.Validate(request); err != nil {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Errors = util.GetValidatorMessage(err)
		return c.JSON(
			http.StatusBadRequest,
			common.NewResponse(res, nil),
		)
	}

	data, cerr := handler.AuthUsecase.Login(ic, request.Email, request.Password)
	if cerr != nil {
		res := common.HandleError(cerr)
		return c.JSON(
			res.StatusCode,
			common.NewResponse(res, nil),
		)
	}

	return c.JSON(
		http.StatusCreated,
		common.NewResponse(
			common.ResultMap[core.OK],
			TokenResponse{
				AccessToken:         data.AccessToken,
				AccessTokenExpired:  time.Unix(data.AccessTokenExpired, 0),
				RefreshToken:        data.RefreshToken,
				RefreshTokenExpired: time.Unix(data.RefreshTokenExpired, 0),
			}),
	)
}

// Relogin godoc
// @Summary 	Get Refresh Token.
// @Description	Get Refresh Token.
// @Tags 		Auth
// @Accept 		json
// @Produce 	json
// @Param		Authorization	header		string				true	"Bearer token"
// @Success		200				{object}	TokenResponse
// @Failure 	400
// @Failure 	401
// @Failure 	500
// @Router /api/v1/auth/relogin [post]
func (handler *Handler) Relogin(c echo.Context) error {
	ic := c.Get("ic").(*core.InternalContext)

	token := c.Get("token").(*jwt.Token)
	claim := token.Claims.(jwt.MapClaims)
	data, cerr := handler.AuthUsecase.RefreshToken(ic, claim["id"].(string))
	if cerr != nil {
		res := common.HandleError(cerr)
		return c.JSON(
			res.StatusCode,
			common.NewResponse(res, nil),
		)
	}

	return c.JSON(
		http.StatusCreated,
		common.NewResponse(
			common.ResultMap[core.OK],
			TokenResponse{
				AccessToken:         data.AccessToken,
				AccessTokenExpired:  time.Unix(data.AccessTokenExpired, 0),
				RefreshToken:        data.RefreshToken,
				RefreshTokenExpired: time.Unix(data.RefreshTokenExpired, 0),
			}),
	)
}

// Logout godoc
// @Summary 	Revoke Token.
// @Description	Revoke Token.
// @Tags 		Auth
// @Accept 		json
// @Produce 	json
// @Param		Authorization	header		string				true	"Bearer token"
// @Success		200
// @Failure 	400
// @Failure 	401
// @Failure 	500
// @Router /api/v1/auth/logout [post]
func (handler *Handler) Logout(c echo.Context) error {
	ic := c.Get("ic").(*core.InternalContext)

	token := c.Get("token").(*jwt.Token)
	claim := token.Claims.(jwt.MapClaims)
	cerr := handler.AuthUsecase.RevokeToken(ic, claim["id"].(string))
	if cerr != nil {
		res := common.HandleError(cerr)
		return c.JSON(
			res.StatusCode,
			common.NewResponse(res, nil),
		)
	}

	res := common.ResultMap[core.OK]
	res.Message = "logout successfully"
	return c.JSON(
		res.StatusCode,
		common.NewResponse(res, nil),
	)
}

// Validate 	godoc
// @Summary 	Validate Token.
// @Description	Validate Token.
// @Tags 		Auth
// @Accept 		json
// @Produce 	json
// @Param		Authorization	header		string	true	"Bearer token"
// @Success		200
// @Failure 	400
// @Failure 	401
// @Failure 	500
// @Router /api/v1/auth/validate [post]
func (handler *Handler) Validate(c echo.Context) error {
	res := common.ResultMap[core.OK]
	res.Message = "jwt token is valid"
	return c.JSON(
		res.StatusCode,
		common.NewResponse(res, nil),
	)
}

// ChangePassword	godoc
// @Summary			Change Password.
// @Description		Change Password.
// @Tags			Auth
// @Accept			json
// @Produce			json
// @Param			Authorization	header		string			true 	"Bearer token"
// @Param  			Body 			body		ChangePassword 	true	"body"
// @Failure			400
// @Failure			401
// @Failure			422
// @Failure			500
// @Router	/api/v1/auth/password/change [post]
func (handler *Handler) ChangePassword(c echo.Context) error {
	ic := c.Get("ic").(*core.InternalContext)
	claims := c.Get("claims").(common.JwtClaims)

	ic.AppendData(map[string]any{
		"claims": claims,
	})

	request := new(ChangePassword)
	if err := c.Bind(request); err != nil {
		res := common.HandleBind(err)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	if err := c.Validate(request); err != nil {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Errors = util.GetValidatorMessage(err)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	if request.Password != request.PasswordConfirm {
		res := common.ResultMap[core.BAD_REQUEST]
		res.Message = "password does not match"
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	account := request.Account()
	account.Id = claims.Id
	cerr := handler.AuthUsecase.ChangePassword(ic, request.OldPassword, account)
	if cerr != nil {
		res := common.HandleError(cerr)
		return c.JSON(res.StatusCode, common.NewResponse(res, nil))
	}

	res := common.ResultMap[core.OK]
	res.Message = "new password was set successfully"
	return c.JSON(
		http.StatusOK,
		common.NewResponse(res, nil),
	)
}
