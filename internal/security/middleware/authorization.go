package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/labstack/echo/v4"
	"refactory/notes/internal/app"
	"refactory/notes/internal/security/token"
	"refactory/notes/internal/web"
)

type Authorization struct {
	enforcer *casbin.Enforcer
}

func NewAuthorization(enforcer *casbin.Enforcer) *Authorization {
	return &Authorization{enforcer: enforcer}
}

func (a *Authorization) Enforce() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, ok := c.Get("session").(*token.Token)
			if !ok {
				return web.ResponseError(c, app.UnauthenticateError)
			}

			authorized, err := a.enforcer.Enforce(token.Username, c.Path(), c.Request().Method)
			if nil != err || !authorized {
				return web.ResponseError(c, app.UnauthorizedError)
			}

			return next(c)
		}
	}
}
