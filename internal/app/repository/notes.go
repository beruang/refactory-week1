package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"refactory/notes/internal/app"
	"refactory/notes/internal/app/model"
)

type NotesRepository interface {
	InsertNotes(ctx context.Context, notes *model.Notes) error
	GetNotes(ctx context.Context, userId int, roleId int) ([]model.Notes, error)
	DetailNotes(ctx context.Context, userId int, id int, roleId int) (model.Notes, error)
	GetSecret(ctx context.Context, id int) (string, error)
	UpdateNotes(ctx context.Context, notes *model.Notes) error
	DeleteNotes(ctx context.Context, userId int, id int) error
	ReActiveNotes(ctx context.Context, id int) error
}

type notesRepository struct {
	db *sqlx.DB
}

func NewNotesRepository(db *sqlx.DB) NotesRepository {
	return &notesRepository{db: db}
}

func (n notesRepository) InsertNotes(ctx context.Context, notes *model.Notes) error {
	stmt, err := n.db.Prepare(`INSERT INTO notes."notes" (user_id, type, title, body, secret)
								VALUES ($1, $2, $3, $4, $5) RETURNING id`)
	if nil != err {
		return errors.Wrap(err, "[db] InsertNotes - prepare statement")
	}

	if err := stmt.QueryRowContext(ctx, notes.UserId, notes.Type, notes.Title, notes.Body, notes.Secret).Scan(&notes.Id); nil != err {
		return errors.Wrap(err, "[db] InsertNotes - insert data")
	}

	return nil
}

func (n notesRepository) GetNotes(ctx context.Context, userId int, roleId int) ([]model.Notes, error) {
	var result []model.Notes

	b := newBuilder(n.db).
		baseQuery(`SELECT id, type, title, body, secret from notes.notes`)
	if roleId == roleUser {
		b.addParam("user_id", userId).addParam("is_active", nil)
	}

	rows, err := b.build().query(ctx)

	if nil != err {
		if err == sql.ErrNoRows {
			return nil, app.NotFoundError
		}
		return nil, errors.Wrap(err, "[db] GetNotes - query")
	}

	for rows.Next() {
		var notes model.Notes
		if err := rows.Scan(&notes.Id, &notes.Type, &notes.Title, &notes.Body, &notes.Secret); nil != err {
			return nil, errors.Wrap(err, "[db] GetNotes - scan struct")
		}
		result = append(result, notes)
	}

	return result, nil
}

func (n notesRepository) DetailNotes(ctx context.Context, userId int, id int, roleId int) (model.Notes, error) {
	var result model.Notes

	query := newBuilder(n.db).baseQuery(`SELECT id, type, title, body, secret from notes.notes`).addParam("id", id)
	if roleId == roleUser {
		query.addParam("user_id", userId).addParam("is_active", nil)
	}

	rows, err := query.build().query(ctx)
	if nil != err {
		if sql.ErrNoRows == err {
			return result, app.NotFoundError
		}
		return result, errors.Wrap(err, "[db] DetailNotes - queries")
	}

	if nil != rows.Err() || !rows.Next() {
		return result, app.NotFoundError
	}
	for rows.Next() {
		if err := rows.Scan(&result.Id, &result.Type, &result.Title, &result.Body, &result.Secret); nil != err {
			return result, errors.Wrap(err, "[db] DetailNotes - scan rows")
		}
	}

	return result, nil
}

func (n notesRepository) GetSecret(ctx context.Context, id int) (string, error) {
	var result string

	rows, err := newBuilder(n.db).
		baseQuery(`SELECT secret FROM notes.notes`).
		addParam("id", id).
		build().query(ctx)

	if nil != err {
		if sql.ErrNoRows == err {
			return "", app.NotFoundError
		}
		return "", app.InternalError
	}

	if nil != rows.Err() || !rows.Next() {
		return "", app.NotFoundError
	}

	if err := rows.Scan(&result); nil != err {
		return "", errors.Wrap(err, "[db] GetSecret - scan rows")
	}

	return result, nil
}

func (n notesRepository) UpdateNotes(ctx context.Context, notes *model.Notes) error {
	query := `UPDATE notes.notes SET type=$1, title=$2, body=$3, secret=$4 where id=$5`
	if roleUser == 0 {
		query = fmt.Sprintf("%s AND user_id=%d", query, notes.UserId)
	}

	stmt, err := n.db.PrepareContext(ctx, query)
	if nil != err {
		return errors.Wrap(err, "[db] UpdateNotes - prepare statement")
	}

	rs, err := stmt.ExecContext(ctx, notes.Type, notes.Title, notes.Body, notes.Secret, notes.Id)
	if nil != err {
		return errors.Wrap(err, "[db] UpdateNotes - update notes")
	}

	updated, _ := rs.RowsAffected()
	if updated == 0 {
		return app.NotFoundError
	}

	return nil
}

func (n notesRepository) DeleteNotes(ctx context.Context, userId int, id int) error {
	query := `UPDATE notes.notes SET is_active=false where id=$1`
	if roleUser == 0 {
		query = fmt.Sprintf("%s AND user_id=%d", query, userId)
	}
	stmt, err := n.db.PrepareContext(ctx, query)
	if nil != err {
		return errors.Wrap(err, "[db] DeleteNotes - prepare statement")
	}

	rs, err := stmt.ExecContext(ctx, id)
	if nil != err {
		return errors.Wrap(err, "[db] DeleteNotes - exec query delete")
	}

	inserted, _ := rs.RowsAffected()
	if inserted == 0 {
		return app.NotFoundError
	}

	return nil
}

func (n notesRepository) ReActiveNotes(ctx context.Context, id int) error {
	stmt, err := n.db.PrepareContext(ctx, `UPDATE notes.notes SET is_active=true where id=$1`)
	if nil != err {
		return errors.Wrap(err, "[db] ReActiveNotes - prepare statement")
	}

	rs, err := stmt.ExecContext(ctx, id)
	if nil != err {
		return errors.Wrap(err, "[db] ReActiveNotes - exec query delete")
	}
	inserted, _ := rs.RowsAffected()
	if inserted == 0 {
		return app.NotFoundError
	}

	return nil
}
