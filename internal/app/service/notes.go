package service

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"refactory/notes/internal/app"
	"refactory/notes/internal/app/model"
	"refactory/notes/internal/app/repository"
)

type NotesService interface {
	CreateNotes(ctx context.Context, notes *model.Notes) (*model.NotesResponse, error)
	GetNotes(ctx context.Context, userId int) ([]*model.NotesResponse, error)
	DetailNotes(ctx context.Context, userId int, id int) (*model.NotesResponse, error)
	EditNotes(ctx context.Context, notes *model.Notes) (*model.NotesResponse, error)
	DeleteNotes(ctx context.Context, userId int, id int) error
}

type notesService struct {
	repo repository.NotesRepository
}

func NewNotesService(repo repository.NotesRepository) NotesService {
	return &notesService{repo: repo}
}

func (n *notesService) CreateNotes(ctx context.Context, notes *model.Notes) (*model.NotesResponse, error) {
	if err := n.repo.InsertNotes(ctx, notes); nil != err {
		if vErr, ok := err.(*pq.Error); ok && vErr.Code == "23505" {
			return nil, app.Error{Code: app.DuplicateCode.Int(), Message: fmt.Sprintf("duplicate value for field %s", vErr.Column)}
		} else {
			return nil, err
		}
	}
	return model.NewNotesResponse(notes.Id, notes.Type, notes.Title, notes.Body, notes.Secret), nil
}

func (n *notesService) GetNotes(ctx context.Context, userId int) ([]*model.NotesResponse, error) {
	var responses []*model.NotesResponse
	result, err := n.repo.GetNotes(ctx, userId)
	if nil != err {
		return nil, err
	}

	if len(result) < 1 {
		return nil, app.NotFoundError
	}

	for _, r := range result {
		responses = append(responses, model.NewNotesResponse(r.Id, r.Type, r.Title, r.Body, r.Secret))
	}

	return responses, nil
}

func (n *notesService) DetailNotes(ctx context.Context, userId int, id int) (*model.NotesResponse, error) {
	result, err := n.repo.DetailNotes(ctx, userId, id)
	if nil != err {
		return nil, err
	}

	return model.NewNotesResponse(result.Id, result.Type, result.Title, result.Body, result.Secret), nil
}

func (n *notesService) EditNotes(ctx context.Context, notes *model.Notes) (*model.NotesResponse, error) {
	if err := n.repo.UpdateNotes(ctx, notes); nil != err {
		return nil, err
	}

	return model.NewNotesResponse(notes.Id, notes.Type, notes.Title, notes.Body, notes.Secret), nil
}

func (n *notesService) DeleteNotes(ctx context.Context, userId int, id int) error {
	if err := n.repo.DeleteNotes(ctx, userId, id); nil != err {
		return err
	}

	return nil
}
