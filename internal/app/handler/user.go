package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/common/log"
	"refactory/notes/internal/app"
	"refactory/notes/internal/app/model"
	"refactory/notes/internal/app/service"
	"refactory/notes/internal/security/token"
	"refactory/notes/internal/web"
	"strconv"
)

type UserHandler interface {
	CreateUser(c echo.Context) error
	Login(c echo.Context) error
	VerifyCode(c echo.Context) error
	ListUser(c echo.Context) error
	DetailUser(c echo.Context) error
	EditUser(c echo.Context) error
	DeleteUser(c echo.Context) error
	ActiveUser(c echo.Context) error
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

	return web.Response(c, result)
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

// @Router /users [get]
// @Tags users
// @Summary List User
// @Description TODO
// @Accept json
// @Produce json
// @Success 200 {array} model.UserResponse
func (u *userHandler) ListUser(c echo.Context) error {
	response, err := u.userService.ListUser(c.Request().Context())
	if nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, response)
}

// @Router /users/{id} [get]
// @Tags users
// @Summary Detail User
// @Description TODO
// @Accept json
// @Produce json
// Param id path int true "user id"
// @Success 200 {object} model.UserResponse
func (u *userHandler) DetailUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if nil != err {
		return web.ResponseError(c, app.BadRequestError)
	}

	response, err := u.userService.DetailUser(c.Request().Context(), id)
	if nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, response)
}

// @Router /users/{id} [put]
// @Tags users
// @Summary Update User
// @Description TODO
// @Accept json
// @Produce json
// Param id path int true "user id"
// @Param payload body model.UserRequest true "body request"
// @Success 200 {array} model.UserResponse
func (u *userHandler) EditUser(c echo.Context) error {
	var req model.UserRequest
	id, err := strconv.Atoi(c.Param("id"))
	if nil != err {
		return web.ResponseError(c, app.BadRequestError)
	}
	c.Bind(&req)
	if err := c.Validate(req); nil != err {
		return web.ResponseError(c, err)
	}

	response, err := u.userService.UpdateUser(c.Request().Context(), model.NewUser(id, req.FirstName, req.LastName, req.Email, req.Username, req.Password, req.Photo, 0))
	if nil != err {
		log.Error(err)
		return web.ResponseError(c, err)
	}

	return web.Response(c, response)
}

// @Router /users/{id} [delete]
// @Tags users
// @Summary Delete User
// @Description TODO
// @Accept json
// @Produce json
// Param id path int true "user id"
// @Success 200 {string} result
func (u *userHandler) DeleteUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if nil != err {
		return web.ResponseError(c, app.BadRequestError)
	}

	if err := u.userService.DeleteUser(c.Request().Context(), id); nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, "User has deleted")
}

// @Router /admin/users/{id} [put]
// @Tags admin
// @Summary Active User
// @Description TODO
// @Accept json
// @Produce json
// Param id path int true "user id"
// @Success 200 {string} result
func (u *userHandler) ActiveUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if nil != err {
		return web.ResponseError(c, app.BadRequestError)
	}

	if err := u.userService.ActiveUser(c.Request().Context(), id); nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, "User has active")
}
