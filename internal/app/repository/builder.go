package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type builder struct {
	db     *sqlx.DB
	base   string
	params map[string]interface{}
}

func newBuilder(db *sqlx.DB) *builder {
	return &builder{db: db, params: make(map[string]interface{})}
}

func (b *builder) baseQuery(base string) *builder {
	b.base = fmt.Sprintf("%s WHERE", base)
	return b
}

func (b *builder) addParam(params string, value interface{}) *builder {
	b.params[params] = value
	return b
}

func (b *builder) build() *builder {
	if len(b.params) > 0 {
		for k, v := range b.params {
			if nil == v {
				b.base = fmt.Sprintf("%s AND %s", b.base, k)
			} else {
				b.base = fmt.Sprintf("%s AND %s=%v", b.base, k, v)
			}
		}
	}

	t := strings.SplitAfter(b.base, "WHERE")

	if t[1] != "" {
		b.base = fmt.Sprintf("%s %s", t[0], t[1][5:])
	} else {
		b.base = strings.Replace(b.base, "WHERE", "", -1)
	}

	return b
}

func (b *builder) query(ctx context.Context) (*sql.Rows, error) {
	return b.db.QueryContext(ctx, b.base)
}
