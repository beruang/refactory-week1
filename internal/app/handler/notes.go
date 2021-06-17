package handler

import (
	"github.com/labstack/echo/v4"
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

// @Router /notes [post]
// @Tags notes
// @Summary Create Notes
// @Description TODO
// @Accept json
// @Produce json
// @Param payload body model.NotesRequest true "body request"
// @Success 200 {object} model.NotesResponse
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

// @Router /notes [get]
// @Tags notes
// @Summary Get List Notes
// @Description TODO
// @Accept json
// @Produce json
// @Success 200 {array} model.NotesResponse
func (n *notesHandler) ListNotes(c echo.Context) error {
	session, ok := c.Get("session").(*token.Token)
	if !ok {
		return web.ResponseError(c, app.InternalError)
	}

	result, err := n.s.GetNotes(c.Request().Context(), session.UserId, session.RoleId)
	if nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, result)
}

// @Router /notes/{id} [get]
// @Tags notes
// @Summary Get Notes Detail
// @Description TODO
// @Accept json
// @Produce json
// @Param id path int true "id notes"
// @Success 200 {object} model.NotesResponse
func (n *notesHandler) GetNotes(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if nil != err {
		return echo.ErrBadRequest
	}

	session, ok := c.Get("session").(*token.Token)
	if !ok {
		return web.ResponseError(c, app.InternalError)
	}

	result, err := n.s.DetailNotes(c.Request().Context(), session.UserId, id, session.RoleId)
	if nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, result)
}

// @Router /notes/{id} [put]
// @Tags notes
// @Summary Update Notes
// @Description TODO
// @Accept json
// @Produce json
// @Param id path int true "id notes"
// @Success 200 {object} model.NotesResponse
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
		return web.ResponseError(c, app.InternalError)
	}

	response, err := n.s.EditNotes(c.Request().Context(), model.NewNotes(id, session.UserId, req.Type, req.Title, req.Body, req.Secret))
	if nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, response)
}

// @Router /notes/{id} [delete]
// @Tags notes
// @Summary Update Notes
// @Description TODO
// @Accept json
// @Produce json
// @Param id path int true "id notes"
// @Success 200 {string} result
func (n *notesHandler) DeleteNotes(c echo.Context) error {
	var req model.SecretRequest
	if err := c.Bind(&req); nil != err {
		return web.ResponseError(c, app.InternalError)
	}
	id, err := strconv.Atoi(c.Param("id"))
	if nil != err {
		return echo.ErrBadRequest
	}

	session, ok := c.Get("session").(*token.Token)
	if !ok {
		web.ResponseError(c, app.InternalError)
	}

	if err := n.s.DeleteNotes(c.Request().Context(), session.UserId, id, req.Secret); nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, "Note has deleted")
}

// @Router /admin/notes/{id} [put]
// @Tags admin
// @Summary ReActive Notes
// @Description TODO
// @Accept json
// @Produce json
// @Param id path int true "id notes"
// @Success 200 {string} result
func (n *notesHandler) ReActiveNotes(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if nil != err {
		return echo.ErrBadRequest
	}

	if err := n.s.ReActiveNotes(c.Request().Context(), id); nil != err {
		return web.ResponseError(c, err)
	}

	return web.Response(c, "Notes has active")
}
