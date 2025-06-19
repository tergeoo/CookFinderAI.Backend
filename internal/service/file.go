package service

import (
	"CookFinder.Backend/internal/model"
	repository "CookFinder.Backend/internal/repo"
	"context"
)

type FileService struct {
	repo *repository.FileRepository
}

func NewFileService(repo *repository.FileRepository) *FileService {
	return &FileService{repo: repo}
}

func (s *FileService) CreateFile(ctx context.Context, file *model.File) error {
	return s.repo.Create(ctx, file)
}

func (s *FileService) GetAllFiles(ctx context.Context) ([]model.File, error) {
	return s.repo.GetAll(ctx)
}

func (s *FileService) GetFileByID(ctx context.Context, id string) (*model.File, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FileService) DeleteFile(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
