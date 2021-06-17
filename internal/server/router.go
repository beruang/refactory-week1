package server

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	middleware2 "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "refactory/notes/docs"
	"refactory/notes/internal/app/handler"
	"refactory/notes/internal/app/repository"
	"refactory/notes/internal/app/service"
	"refactory/notes/internal/db/redis"
	"refactory/notes/internal/security/middleware"
)

type handlerModule struct {
	user  handler.UserHandler
	notes handler.NotesHandler
	media handler.MediaHandler
}

// @title RSP Notes API
// @version 1.0
// @description Implementing back-end services for RSP Notes application
// @BasePath /api
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
func NewRouter(validate *validator.Validate, db *sqlx.DB, cache redis.Client, enforcer *casbin.Enforcer) *echo.Echo {
	authenticationMiddleware := middleware.NewAuthorization(enforcer)
	e := echo.New()

	e.Validator = &CustomValidator{validate}

	module := getModule(db, cache, enforcer)
	api := e.Group("/api")
	api.POST("/registrasi", module.user.CreateUser)
	api.POST("/verification", module.user.VerifyCode, middleware.Claim(), middleware.Auth)
	api.POST("/login", module.user.Login)

	user := api.Group("/users", middleware.Claim(), middleware.Auth, authenticationMiddleware.Enforce())
	user.GET("", module.user.ListUser)
	user.GET("/:id", module.user.DetailUser)
	user.PUT("/:id", module.user.EditUser)
	user.DELETE("/:id", module.user.DeleteUser)

	notes := api.Group("/notes", middleware.Claim(), middleware.Auth, authenticationMiddleware.Enforce())
	notes.POST("", module.notes.CreateNotes)
	notes.GET("", module.notes.ListNotes)
	notes.GET("/:id", module.notes.GetNotes)
	notes.PUT("/:id", module.notes.EditNotes)
	notes.DELETE("/:id", module.notes.DeleteNotes)

	media := api.Group("/media")
	media.POST("", module.media.UploadMedia, middleware2.BodyLimit("10M"), middleware.Claim(), middleware.Auth)
	media.GET("/:id", module.media.DownloadMedia)

	admin := api.Group("/admin", middleware.Claim(), middleware.Auth, authenticationMiddleware.Enforce())
	admin.PUT("/notes/:id", module.notes.ReActiveNotes)
	admin.PUT("/users/:id", module.user.ActiveUser)

	api.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}

func getModule(db *sqlx.DB, cache redis.Client, enforcer *casbin.Enforcer) handlerModule {
	// user module
	userRepo := repository.NewUserRepository(db, cache, enforcer)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// notes module
	notesRepo := repository.NewNotesRepository(db)
	notesService := service.NewNotesService(notesRepo)
	notesHandler := handler.NewNotesHandler(notesService)

	// media module
	mediaRepo := repository.NewMediaRepository(db)
	mediaService := service.NewMediaService(mediaRepo)
	mediaHandler := handler.NewMediaHandler(mediaService)

	return handlerModule{user: userHandler, notes: notesHandler, media: mediaHandler}
}
