package handler

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/common/log"
	"net/http"
	"refactory/notes/internal/app"
	"refactory/notes/internal/app/model"
	"refactory/notes/internal/app/service"
	"refactory/notes/internal/security/token"
	"refactory/notes/internal/web"
	"strconv"
)

type MediaHandler interface {
	UploadMedia(ctx echo.Context) error
	DownloadMedia(ctx echo.Context) error
}

type mediaHandler struct {
	s service.MediaService
}

func NewMediaHandler(s service.MediaService) MediaHandler {
	return &mediaHandler{s: s}
}

var (
	supportedMedia = map[string]bool{
		"image/png":  true,
		"image/jpeg": true,
	}
)

// @Router /media [post]
// @Tags media
// @Summary upload media
// @Description TODO
// @Accept jpeg
// @Produce json
// @Success 200 {object} model.MediaResponse
func (m mediaHandler) UploadMedia(ctx echo.Context) error {
	session, ok := ctx.Get("session").(*token.Token)
	if !ok {
		return web.ResponseError(ctx, app.InternalError)
	}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(ctx.Request().Body)
	if nil != err {
		log.Error(err)
		return echo.ErrInternalServerError
	}

	if !supportedMedia[ctx.Request().Header.Get(echo.HeaderContentType)] {
		return web.ResponseError(ctx, app.BadRequestError)
	}

	id, err := m.s.SaveMedia(ctx.Request().Context(), session.UserId, ctx.Request().Header.Get(echo.HeaderContentType), buf.Bytes())
	if nil != err {
		log.Error(err)
		return web.ResponseError(ctx, err)
	}

	return web.Response(ctx, model.MediaResponse{Id: id})
}

// @Router /media/{id} [get]
// @Tags media
// @Summary download media
// @Description TODO
// @Accept json
// @Produce jpeg
// @Param id path int true "id media"
// @Success 200 {string} binary
func (m mediaHandler) DownloadMedia(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if nil != err {
		log.Error(err)
		return web.ResponseError(ctx, app.BadRequestError)
	}

	mime, picture, err := m.s.GetMedia(ctx.Request().Context(), id)
	if nil != err {
		log.Error(err)
		return web.ResponseError(ctx, err)
	}

	return ctx.Blob(http.StatusOK, mime, picture)
}
