package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"refactory/notes/internal/security/token"
)

func Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		session := user.Claims.(*token.Token)
		c.Set("session", session)

		return next(c)
	}
}

func Claim() echo.MiddlewareFunc {
	conf := middleware.JWTConfig{
		Claims:     new(token.Token),
		SigningKey: []byte("secret"),
	}
	return middleware.JWTWithConfig(conf)
}
