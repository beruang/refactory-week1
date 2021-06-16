package web

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"refactory/notes/internal/app"
	"refactory/notes/internal/translator"
)

type GeneralResponse struct {
	Meta   Meta        `json:"meta"`
	Data   interface{} `json:"data,omitempty"`
	Errors []Error     `json:"errors,omitempty"`
}

func (g GeneralResponse) AddErrors(err Error) GeneralResponse {
	g.Errors = append(g.Errors, err)
	return g
}

func (g GeneralResponse) SetErrors(errs []Error) GeneralResponse {
	g.Errors = errs
	return g
}

func (g GeneralResponse) SetData(data interface{}) GeneralResponse {
	g.Data = data
	return g
}

type Meta struct {
	Version   string   `json:"version"`
	Authors   []string `json:"authors"`
	Copyright string   `json:"copyright"`
}

type Error struct {
	Code    errCode `json:"code"`
	Message string  `json:"message"`
}

type errCode int

const (
	internalServerCode errCode = iota + 1
	badRequestCode
	unauthorized
	unauthenticated
	duplicateCode
	notfoundCode
)

var defaultResponse GeneralResponse = GeneralResponse{Meta: Meta{Version: "1.0.0", Copyright: "Copyright 2021 Refactory.id", Authors: []string{"Zacky Mughni Mubarok"}}}

func ResponseError(c echo.Context, err error) error {
	if vErrs, ok := errors.Cause(err).(validator.ValidationErrors); ok {
		var errs []Error
		uni := translator.GetTranslator()
		trans, _ := uni.GetTranslator("en")
		for _, vErr := range vErrs {
			errs = append(errs, Error{Code: badRequestCode, Message: vErr.Translate(trans)})
		}
		return c.JSON(http.StatusBadRequest, defaultResponse.SetErrors(errs))
	}

	if vErrs, ok := errors.Cause(err).(app.Error); ok {
		switch vErrs {
		case app.UnauthorizedError:
			return c.JSON(http.StatusUnauthorized, defaultResponse.AddErrors(Error{Code: unauthorized, Message: "You are Unauthorized to use this resource"}))
		case app.UnauthenticateError:
			return c.JSON(http.StatusUnauthorized, defaultResponse.AddErrors(Error{Code: unauthenticated, Message: "You are Unauthenticated"}))
		case app.DuplicateError:
			return c.JSON(http.StatusConflict, defaultResponse.AddErrors(Error{Code: duplicateCode, Message: "Data Already Exists"}))
		case app.NotFoundError:
			return c.JSON(http.StatusNotFound, defaultResponse.AddErrors(Error{Code: notfoundCode, Message: "Data Not Found"}))
		}
	}

	return c.JSON(http.StatusInternalServerError, defaultResponse.AddErrors(Error{Code: internalServerCode, Message: "Internal Server Error"}))
}

func Response(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, defaultResponse.SetData(data))
}
