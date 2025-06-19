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

func (it *FileService) CreateFile(ctx context.Context, file *model.File) error {
	return it.repo.Create(ctx, file)
}

func (it *FileService) GetAllFiles(ctx context.Context) ([]model.File, error) {
	return it.repo.GetAll(ctx)
}

func (it *FileService) GetFileByID(ctx context.Context, id string) (*model.File, error) {
	return it.repo.GetByID(ctx, id)
}

func (it *FileService) DeleteFile(ctx context.Context, id string) error {
	return it.repo.Delete(ctx, id)
}
