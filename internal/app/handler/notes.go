package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"refactory/notes/internal/app"
	"refactory/notes/internal/app/model"
	"refactory/notes/internal/app/service"
	"refactory/notes/internal/security/token"
	"refactory/notes/internal/web"
	"strconv"
)

type NotesHandler interface {
	CreateNotes(c echo.Context) error
	ListNotes(c echo.Context) error
	GetNotes(c echo.Context) error
	EditNotes(c echo.Context) error
	DeleteNotes(c echo.Context) error
	ReActiveNotes(c echo.Context) error
}

type notesHandler struct {
	s service.NotesService
}

func NewNotesHandler(s service.NotesService) *notesHandler {
	return &notesHandler{s: s}
}

func (n *notesHandler) CreateNotes(c echo.Context) error {
	var req model.NotesRequest
	if err := c.Bind(&req); nil != err {
		return echo.ErrBadRequest
	}
	if err := c.Validate(&req); nil != err {
		return web.ResponseError(c, err)
	}

	session, ok := c.Get("session").(*token.Token)
	if !ok {
		return web.ResponseError(c, app.InternalError)
	}

	notes := model.NewNotes(0, session.UserId, req.Type, req.Title, req.Body, req.Secret)
	response, err := n.s.CreateNotes(c.Request().Context(), notes)
	if nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, response)
}

func (n *notesHandler) ListNotes(c echo.Context) error {
	session, ok := c.Get("session").(*token.Token)
	if !ok {
		web.ResponseError(c, app.InternalError)
	}

	result, err := n.s.GetNotes(c.Request().Context(), session.UserId)
	if nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, result)
}

func (n *notesHandler) GetNotes(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if nil != err {
		return echo.ErrBadRequest
	}

	session, ok := c.Get("session").(*token.Token)
	if !ok {
		web.ResponseError(c, app.InternalError)
	}

	result, err := n.s.DetailNotes(c.Request().Context(), session.UserId, id)
	if nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, result)
}

func (n *notesHandler) EditNotes(c echo.Context) error {
	var req model.NotesRequest

	id, err := strconv.Atoi(c.Param("id"))
	if nil != err {
		return echo.ErrBadRequest
	}
	if err := c.Bind(&req); nil != err {
		return echo.ErrBadRequest
	}
	if err := c.Validate(&req); nil != err {
		return web.ResponseError(c, err)
	}
	session, ok := c.Get("session").(*token.Token)
	if !ok {
		web.ResponseError(c, app.InternalError)
	}

	response, err := n.s.EditNotes(c.Request().Context(), model.NewNotes(id, session.UserId, req.Type, req.Title, req.Body, req.Secret))
	if nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, response)
}

func (n *notesHandler) DeleteNotes(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if nil != err {
		return echo.ErrBadRequest
	}

	session, ok := c.Get("session").(*token.Token)
	if !ok {
		web.ResponseError(c, app.InternalError)
	}

	if err := n.s.DeleteNotes(c.Request().Context(), session.UserId, id); nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, "Note has deleted")
}

func (n *notesHandler) ReActiveNotes(c echo.Context) error {
	return c.JSON(http.StatusOK, "OK")
}
