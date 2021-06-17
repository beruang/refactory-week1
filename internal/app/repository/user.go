package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/go-redis/cache/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"refactory/notes/internal/app"
	"refactory/notes/internal/app/model"
	"refactory/notes/internal/db/redis"
	"refactory/notes/internal/security/token"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User, verifSession model.Session) (string, error)
	FindUser(ctx context.Context, username string) (*model.User, error)
	VerifyUser(ctx context.Context, username string) error
	UpdateSession(ctx context.Context, session model.Session) error
	FindSession(ctx context.Context, username string) (model.Session, error)
	ListUser(ctx context.Context) ([]*model.User, error)
	DetailUser(ctx context.Context, Id int) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id int) error
	ActiveUser(ctx context.Context, id int) error
}

const (
	roleUser int = iota + 1
	roleAdmin
)

type userRepository struct {
	db       *sqlx.DB
	cache    redis.Client
	enforcer *casbin.Enforcer
}

func NewUserRepository(db *sqlx.DB, cache redis.Client, enforcer *casbin.Enforcer) *userRepository {
	return &userRepository{db: db, cache: cache, enforcer: enforcer}
}

func (u *userRepository) Create(ctx context.Context, user *model.User, session model.Session) (string, error) {
	tx, err := u.db.Begin()
	if nil != err {
		return "", errors.Wrap(err, "[db] CreateUser - begin transaction")
	}
	defer tx.Commit()
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO notes."user" (first_name, last_name, email, password, username, role_id) VALUES ($1, $2, $3, $4, $5, $6) returning id`)
	if nil != err {
		return "", errors.Wrap(err, "[db] CreateUser - prepared statement")
	}
	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, user.FirstName, user.LastName, user.Email, user.Password, user.Username, user.Role).Scan(&user.Id); nil != err {
		tx.Rollback()
		return "", err
	}

	if err := u.UpdateSession(ctx, session); nil != err {
		tx.Rollback()
		return "", errors.Wrap(err, "[db] CreateUser - save to cache")
	}

	t, err := token.GenerateToken(session)
	if nil != err {
		tx.Rollback()
		return "", errors.Wrap(err, "[db] CreateUser - generate token")
	}

	return t, nil
}

func (u *userRepository) UpdateSession(ctx context.Context, session model.Session) error {
	key := fmt.Sprintf("session:%s", session.Username)

	if err := u.cache.Cache().Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: session,
		TTL:   time.Minute * 30,
	}); nil != err {
		return errors.Wrap(err, "[db] UpdateSession - save to cache")
	}

	return nil
}

func (u *userRepository) VerifyUser(ctx context.Context, username string) error {
	stmt, err := u.db.PrepareContext(ctx, `UPDATE notes."user" SET is_verified=true where username=$1`)
	if nil != err {
		return errors.Wrap(err, "[db] VerifyUser - prepared statement")
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, username); nil != err {
		return errors.Wrap(err, "[db] VerifyUser - update db")
	}

	var session model.Session
	key := fmt.Sprintf("session:%s", username)
	if err := u.cache.Cache().Get(ctx, key, &session); nil != err {
		return errors.Wrap(err, "[db] VerifyUser - get cache")
	}

	session.IsVerified = true
	if err := u.UpdateSession(ctx, session); nil != err {
		return errors.Wrap(err, "[db] VerifyUser - update cache")
	}

	u.enforcer.AddRoleForUser(username, "user")

	return nil
}

func (u *userRepository) FindSession(ctx context.Context, username string) (model.Session, error) {
	var session model.Session
	if err := u.cache.Cache().Get(ctx, fmt.Sprintf("session:%s", username), &session); nil != err && cache.ErrCacheMiss != err {
		return session, err
	}

	return session, nil
}

func (u *userRepository) FindUser(ctx context.Context, username string) (*model.User, error) {
	var result model.User
	stmt, err := u.db.PrepareContext(ctx, `SELECT id, first_name, last_name, email, password, username, is_verified, role_id, is_active from notes."user" where username=$1 and is_active`)
	if nil != err {
		return nil, errors.Wrap(err, "[db] FindUser - prepared statement")
	}
	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, username).Scan(&result.Id, &result.FirstName, &result.LastName,
		&result.Email, &result.Password, &result.Username, &result.IsVerified, &result.Role, &result.IsActive); nil != err {
		return nil, err
	}

	if err := u.cache.Cache().Set(&cache.Item{
		Ctx:   ctx,
		Key:   fmt.Sprintf("find:user:%s", username),
		Value: result,
		TTL:   time.Minute * 30,
	}); nil != err {
		log.Error(errors.Wrap(err, "[rdr] FindUser - update cache"))
	}

	return &result, nil
}

func (u *userRepository) ListUser(ctx context.Context) ([]*model.User, error) {
	var result []*model.User
	stmt, err := u.db.PrepareContext(ctx, `SELECT id, first_name, last_name, email, password, username, media_id, role_id FROM notes."user"`)
	if nil != err {
		return nil, errors.Wrap(err, "[db] ListUser - prepare statement")
	}

	rows, err := stmt.Query()
	if nil != err {
		if sql.ErrNoRows == err {
			return nil, app.NotFoundError
		}
		return nil, errors.Wrap(err, "[db] ListUser - query")
	}

	if nil != rows.Err() || !rows.Next() {
		return nil, app.NotFoundError
	}

	for rows.Next() {
		user := new(model.User)
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Username, &user.MediaId, &user.Role); nil != err {
			return nil, errors.Wrap(err, "[db] ListUser - scan")
		}
		result = append(result, user)
	}

	return result, nil
}

func (u *userRepository) DetailUser(ctx context.Context, Id int) (*model.User, error) {
	user := new(model.User)
	stmt, err := u.db.PrepareContext(ctx, `SELECT id, first_name, last_name, email, password, username, media_id, role_id FROM notes."user" WHERE id=$1`)
	if nil != err {
		return nil, errors.Wrap(err, "[db] DetailUser - prepare statement")
	}

	row := stmt.QueryRow(Id)
	if nil != row.Err() {
		if sql.ErrNoRows == row.Err() {
			return nil, app.NotFoundError
		}
		return nil, errors.Wrap(err, "[db] DetailUser - query")
	}

	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Username, &user.MediaId, &user.Role); nil != err {
		return nil, errors.Wrap(err, "[db] DetailUser - scan")
	}

	return user, nil
}

func (u *userRepository) UpdateUser(ctx context.Context, user *model.User) error {
	stmt, err := u.db.PrepareContext(ctx, `UPDATE notes."user" SET first_name=$2, last_name=$3, email=$4, username=$5 WHERE id=$1 AND is_active`)
	if nil != err {
		return errors.Wrap(err, "[db] UpdateUser - prepare statement")
	}
	rs, err := stmt.Exec(user.Id, user.FirstName, user.LastName, user.Email, user.Username)
	if nil != err {
		return errors.Wrap(err, "[db] UpdateUser- query")
	}
	updated, _ := rs.RowsAffected()
	if updated == 0 {
		return app.NotFoundError
	}

	return nil
}

func (u *userRepository) DeleteUser(ctx context.Context, id int) error {
	stmt, err := u.db.PrepareContext(ctx, `UPDATE notes."user" SET is_active=false where id=$1`)
	if nil != err {
		return errors.Wrap(err, "[db] DeleteUser - prepare statement")
	}

	rs, err := stmt.Exec(id)
	if nil != err {
		return errors.Wrap(err, "[db] DeleteUser - query")
	}

	affected, _ := rs.RowsAffected()
	if affected == 0 {
		return app.NotFoundError
	}

	return nil
}

func (u *userRepository) ActiveUser(ctx context.Context, id int) error {
	stmt, err := u.db.PrepareContext(ctx, `UPDATE notes."user" SET is_active=true where id=$1`)
	if nil != err {
		return errors.Wrap(err, "[db] DeleteUser - prepare statement")
	}

	rs, err := stmt.Exec(id)
	if nil != err {
		return errors.Wrap(err, "[db] DeleteUser - query")
	}

	affected, _ := rs.RowsAffected()
	if affected == 0 {
		return app.NotFoundError
	}

	return nil
}
