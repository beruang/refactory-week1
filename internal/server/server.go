package server

import (
	sqlxadapter "github.com/Blank-Xu/sqlx-adapter"
	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"net/http"
	"refactory/notes/internal/config"
	"refactory/notes/internal/db/postgres"
	"refactory/notes/internal/db/redis"
	"refactory/notes/internal/translator"
)

func Start() error {
	validate, err := translator.GetValidator()
	if nil != err {
		return errors.Wrap(err, "registering translator")
	}

	httpServer := &http.Server{
		Addr: config.Cfg().WebAddress,
	}

	db, err := postgres.Open()
	if nil != err {
		return errors.Wrap(err, "initialize postgres driver")
	}

	cache, err := redis.NewClient()
	if nil != err {
		return errors.Wrap(err, "initialize caching")
	}

	enforcer, err := createEnforcer(db)
	if nil != err {
		return err
	}

	router := NewRouter(validate, db, cache, enforcer)

	err = router.StartServer(httpServer)
	if nil != err {
		return err
	}

	return nil
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func createEnforcer(db *sqlx.DB) (*casbin.Enforcer, error) {
	adapter, err := sqlxadapter.NewAdapter(db, "rules")
	if nil != err {
		return nil, err
	}

	enforcer, err := casbin.NewEnforcer("./casbin/model.conf", adapter)
	if nil != err {
		return nil, err
	}

	if err := enforcer.LoadModel(); nil != err {
		return nil, err
	}
	if err := enforcer.LoadPolicy(); nil != err {
		return nil, err
	}
	if err := enforcer.SavePolicy(); nil != err {
		return nil, err
	}

	enforcer.EnableAutoSave(true)
	return enforcer, nil
}
