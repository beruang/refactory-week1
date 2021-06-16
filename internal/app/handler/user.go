package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"refactory/notes/internal/app"
	"refactory/notes/internal/app/model"
	"refactory/notes/internal/app/service"
	"refactory/notes/internal/security/token"
	"refactory/notes/internal/web"
)

type UserHandler interface {
	CreateUser(c echo.Context) error
	Login(c echo.Context) error
	VerifyCode(c echo.Context) error
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{userService: userService}
}

// @Router /registrasi [post]
// @Tags registrasi
// @Summary Create User
// @Description TODO
// @Accept json
// @Produce json
// @Param payload body model.UserRequest true "body request"
// @Success 200 {object} model.UserResponse
func (u *userHandler) CreateUser(c echo.Context) error {
	var req model.UserRequest

	if err := c.Bind(&req); nil != err {
		return echo.ErrBadRequest
	}

	if err := c.Validate(req); nil != err {
		return web.ResponseError(c, err)
	}

	result, err := u.userService.CreateUser(c.Request().Context(), req)
	if nil != err {
		return web.ResponseError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// @Router /verification [post]
// @Tags registrasi
// @Summary verified verification code
// @Description TODO
// @Accept json
// @Produce json
// @Param payload body model.LoginRequest true "body request"
// @Success 200
func (u *userHandler) VerifyCode(c echo.Context) error {
	var req model.VerifyRequest
	if err := c.Bind(&req); nil != err {
		return echo.ErrBadRequest
	}

	if err := c.Validate(req); nil != err {
		return web.ResponseError(c, err)
	}

	session, ok := c.Get("session").(*token.Token)
	if !ok {
		return web.ResponseError(c, app.InternalError)
	}

	if err := u.userService.VerifyCode(c.Request().Context(), *session, req.Code); nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, "")
}

// @Router /login [post]
// @Tags login
// @Summary Login User
// @Description TODO
// @Accept json
// @Produce json
// @Param payload body model.LoginRequest true "body request"
// @Success 200 {object} model.LoginResponse
func (u *userHandler) Login(c echo.Context) error {
	var req model.LoginRequest

	if err := c.Bind(&req); nil != err {
		return web.ResponseError(c, err)
	}

	if err := c.Validate(req); nil != err {
		return web.ResponseError(c, err)
	}

	resp, err := u.userService.Login(c.Request().Context(), req.Username, req.Password)
	if nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, resp)
}
