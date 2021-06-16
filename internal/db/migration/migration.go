package migration

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"net/url"
	"refactory/notes/internal/config"
)

func getUrl() (sourceUrl, databaseUrl string) {
	sourceUrl = "file://./migrations"

	q := make(url.Values)
	q.Set("sslmode", "disable")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(config.Cfg().PgUser, config.Cfg().PgPassword),
		Host:     fmt.Sprintf("%s:%s", config.Cfg().PgHost, config.Cfg().PgPort),
		Path:     config.Cfg().PgName,
		RawQuery: q.Encode(),
	}
	databaseUrl = u.String()
	return
}

func Up() error {
	m, err := migrate.New(getUrl())
	if nil != err {
		return err
	}

	err = m.Up()
	return ignoreErrNoChange(err)
}

func Down() error {
	m, err := migrate.New(getUrl())
	if nil != err {
		return err
	}
	err = m.Down()
	return ignoreErrNoChange(err)
}

func Steps(n int) error {
	m, err := migrate.New(getUrl())
	if nil != err {
		return err
	}
	err = m.Steps(n)
	return ignoreErrNoChange(err)
}

func Drop() error {
	m, err := migrate.New(getUrl())
	if nil != err {
		return err
	}
	err = m.Drop()
	return ignoreErrNoChange(err)
}

func ignoreErrNoChange(err error) error {
	if nil != err && migrate.ErrNoChange != err {
		return err
	}
	return nil
}
