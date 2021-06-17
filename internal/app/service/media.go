package service

import (
	"context"
	"refactory/notes/internal/app/repository"
)

type MediaService interface {
	SaveMedia(ctx context.Context, userId int, mime string, file []byte) (int, error)
	GetMedia(ctx context.Context, id int) (string, []byte, error)
}

type mediaService struct {
	repo repository.MediaRepository
}

func NewMediaService(repo repository.MediaRepository) MediaService {
	return &mediaService{repo: repo}
}

func (m *mediaService) SaveMedia(ctx context.Context, userId int, mime string, file []byte) (int, error) {
	return m.repo.InsertMedia(ctx, userId, mime, file)
}

func (m *mediaService) GetMedia(ctx context.Context, id int) (string, []byte, error) {
	return m.repo.SelectMedia(ctx, id)
}
