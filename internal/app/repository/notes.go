package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"refactory/notes/internal/app"
	"refactory/notes/internal/app/model"
)

type NotesRepository interface {
	InsertNotes(ctx context.Context, notes *model.Notes) error
	GetNotes(ctx context.Context, userId int) ([]model.Notes, error)
	DetailNotes(ctx context.Context, userId int, id int) (model.Notes, error)
	UpdateNotes(ctx context.Context, notes *model.Notes) error
	DeleteNotes(ctx context.Context, userId int, id int) error
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

func (n notesRepository) GetNotes(ctx context.Context, userId int) ([]model.Notes, error) {
	var result []model.Notes
	stmt, err := n.db.Prepare(`SELECT id, type, title, body, secret from notes.notes where user_id=$1 and is_active`)
	if nil != err {
		return nil, errors.Wrap(err, "[db] GetNotes - prepared statement")
	}

	rows, err := stmt.QueryContext(ctx, userId)
	if nil != err {
		return nil, errors.Wrap(err, "[db] GetNotes - get notes")
	}

	for rows.Next() {
		r := model.Notes{}
		if err := rows.Scan(&r.Id, &r.Type, &r.Title, &r.Body, &r.Secret); nil != err {
			return nil, errors.Wrap(err, "[db] GetNotes - scanning result")
		}
		result = append(result, r)
	}

	return result, nil
}

func (n notesRepository) DetailNotes(ctx context.Context, userId int, id int) (model.Notes, error) {
	var result model.Notes
	stmt, err := n.db.Prepare(`SELECT id, type, title, body, secret from notes.notes where user_id=$1 and id=$2 and is_active`)
	if nil != err {
		return result, errors.Wrap(err, "[db] DetailNotes - prepare statement")
	}

	if err := stmt.QueryRowContext(ctx, userId, id).Scan(&result.Id, &result.Type, &result.Title,
		&result.Body, &result.Secret); nil != err {
		if err == sql.ErrNoRows {
			return result, app.NotFoundError
		}
		return result, errors.Wrap(err, "[db] DetailNotes - query context")
	}

	return result, nil
}

func (n notesRepository) UpdateNotes(ctx context.Context, notes *model.Notes) error {
	stmt, err := n.db.PrepareContext(ctx, `UPDATE notes.notes SET type=$1, title=$2, body=$3, secret=$4 where id=$5 and user_id=$6`)
	if nil != err {
		return errors.Wrap(err, "[db] UpdateNotes - prepare statement")
	}

	_, err = stmt.ExecContext(ctx, notes.Type, notes.Title, notes.Body, notes.Secret, notes.Id, notes.UserId)
	if nil != err {
		return errors.Wrap(err, "[db] UpdateNotes - update notes")
	}

	return nil
}

func (n notesRepository) DeleteNotes(ctx context.Context, userId int, id int) error {
	stmt, err := n.db.PrepareContext(ctx, `UPDATE notes.notes SET is_active=false where user_id=$1 and id=$2`)
	if nil != err {
		return errors.Wrap(err, "[db] DeleteNotes - prepare statement")
	}

	_, err = stmt.ExecContext(ctx, userId, id)
	if nil != err {
		return errors.Wrap(err, "[db] DeleteNotes - exec query delete")
	}

	return nil
}
