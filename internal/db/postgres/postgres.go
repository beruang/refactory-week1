package postgres

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"net/url"
	"refactory/notes/internal/config"
)

func Open() (*sqlx.DB, error) {
	q := make(url.Values)
	q.Set("sslmode", "disable")
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(config.Cfg().PgUser, config.Cfg().PgPassword),
		Host:     fmt.Sprintf("%s:%s", config.Cfg().PgHost, config.Cfg().PgPort),
		Path:     config.Cfg().PgName,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("postgres", u.String())
	if nil != err {
		return nil, err
	}

	if err := status(context.Background(), db); nil != err {
		return nil, err
	}

	return db, nil
}

func status(ctx context.Context, db *sqlx.DB) error {
	var up bool
	return db.QueryRowContext(ctx, "SELECT true").Scan(&up)
}
