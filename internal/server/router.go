package server

import (
	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
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
}

// @title Go Blog API
// @version 1.0
// @description Implementing back-end services for blog application
// @BasePath /v1
func NewRouter(validate *validator.Validate, db *sqlx.DB, cache redis.Client, enforcer *casbin.Enforcer) *echo.Echo {
	authenticationMiddleware := middleware.NewAuthorization(enforcer)
	e := echo.New()

	e.Validator = &CustomValidator{validate}

	module := getModule(db, cache, enforcer)
	api := e.Group("/api")
	api.POST("/registrasi", module.user.CreateUser)
	api.POST("/verification", module.user.VerifyCode, middleware.Claim(), middleware.Auth)
	api.POST("/login", module.user.Login)

	notes := api.Group("/notes", middleware.Claim(), middleware.Auth, authenticationMiddleware.Enforce())
	notes.POST("", module.notes.CreateNotes)
	notes.GET("", module.notes.ListNotes)
	notes.GET("/:id", module.notes.GetNotes)
	notes.PUT("/:id", module.notes.EditNotes)
	notes.DELETE("/:id", module.notes.DeleteNotes)

	admin := api.Group("/admin", middleware.Claim(), middleware.Auth, authenticationMiddleware.Enforce())
	admin.PUT("/notes/:id", module.notes.ReActiveNotes)

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

	return handlerModule{user: userHandler, notes: notesHandler}
}
