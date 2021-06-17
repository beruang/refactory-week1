package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"refactory/notes/internal/app"
)

type MediaRepository interface {
	InsertMedia(ctx context.Context, userId int, mime string, file []byte) (int, error)
	SelectMedia(ctx context.Context, id int) (string, []byte, error)
}

type mediaRepository struct {
	db *sqlx.DB
}

func NewMediaRepository(db *sqlx.DB) MediaRepository {
	return &mediaRepository{db: db}
}

func (m *mediaRepository) InsertMedia(ctx context.Context, userId int, mime string, file []byte) (int, error) {
	var mediaId int
	tx, err := m.db.Begin()
	if nil != err {
		return 0, errors.Wrap(err, "[db] InsertMedia - begin transaction")
	}

	rows := tx.QueryRow(`INSERT INTO notes.media(mime_type, file) VALUES ($1, $2) RETURNING id`, mime, file)
	if rows.Err() != nil {
		return 0, errors.Wrap(err, "[db] InsertMedia - insert media file")
	}
	if nil != rows.Scan(&mediaId) {
		tx.Rollback()
		return 0, errors.Wrap(err, "[db] InsertMedia - scan")
	}

	stmt, err := tx.PrepareContext(ctx, `UPDATE notes."user" SET media_id=$2 WHERE id=$1`)
	if nil != err {
		tx.Rollback()
		return 0, errors.Wrap(err, "[db] InsertMedia - prepare statement")
	}

	rs, err := stmt.Exec(userId, mediaId)
	if nil != err {
		tx.Rollback()
		if sql.ErrNoRows == err {
			return 0, app.NotFoundError
		}
		return 0, errors.Wrap(err, "[db] InsertMedia - update picture")
	}

	updated, _ := rs.RowsAffected()
	if updated == 0 {
		tx.Rollback()
		return 0, app.NotFoundError
	}

	if err := tx.Commit(); nil != err {
		tx.Rollback()
		return 0, errors.Wrap(err, "[db] InsertMedia - commit transaction")
	}

	return mediaId, nil
}

func (m *mediaRepository) SelectMedia(ctx context.Context, id int) (string, []byte, error) {
	var file []byte
	var mime string

	stmt, err := m.db.PrepareContext(ctx, `SELECT mime_type, file FROM notes.media WHERE id=$1`)
	if nil != err {
		return "", nil, errors.Wrap(err, "[db] SelectMedia - prepare statement")
	}

	row := stmt.QueryRow(id)
	if nil != row.Err() {
		if sql.ErrNoRows == row.Err() {
			return "", nil, app.NotFoundError
		}
		return "", nil, errors.Wrap(err, "[db] SelectMedia - query media")
	}

	if err := row.Scan(&mime, &file); nil != err {
		return "", nil, errors.Wrap(err, "[db] SelectMedia - scan")
	}

	return mime, file, nil
}
