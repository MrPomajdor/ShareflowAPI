package storage

import (
	"context"
	"mime/multipart"

	"github.com/MrPomajdor/ShareFlowAPI/internal/entity"
	"github.com/MrPomajdor/ShareFlowAPI/pkg/dbcontext"
	"github.com/sirupsen/logrus"
)

type Service interface {
	// Uploads provided file to server
	Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, category_name string) (*entity.FileNode, error)
	// Removes a file with provided path from the server
	Remove(ctx context.Context, fileID int) error
	// Moves a file from provided path to a new destination
	Move(ctx context.Context, r MoveRequest) error
	// GetRoot returns a list of names of categories that user owns
	GetCategories(ctx context.Context) ([]string, error)
	// GetRoot returns a list of names of categories that user owns
	GetCategoryContent(ctx context.Context, category_name string) ([]*entity.FileNode, error)
	// CreateURL creates and provides a link for sharing the selected file
	CreateURL(ctx context.Context, r CreateURLRequest) error
	// Returns service logger
	GetLogger() *logrus.Logger
	// Returns storage path for user files. Set in config
	GetStoragePath() string
}

type service struct {
	repo       Repository
	db         *dbcontext.DB
	logger     *logrus.Logger
	uploadPath string
}

func (s service) GetLogger() *logrus.Logger {
	return s.logger
}

func (s service) GetStoragePath() string {
	return s.uploadPath
}

func NewService(repo Repository, logger *logrus.Logger, db *dbcontext.DB, uploadPath string) Service {
	return service{repo, db, logger, uploadPath}
}

func (s service) StoreFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, category_name string) (*entity.FileNode, error) {

	file_s, err := s.repo.InitializeNewFile(ctx, header.Filename, category_name)
	if err != nil {
		return nil, err
	}
	if err := s.repo.WriteFile(ctx, file, file_s); err != nil {
		return nil, err
	}
	return file_s, nil
}

func (s service) Upload(ctx context.Context, file multipart.File, header *multipart.FileHeader, category_name string) (*entity.FileNode, error) {
	f, err := s.StoreFile(ctx, file, header, category_name)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (s service) Remove(ctx context.Context, fileID int) error {
	return nil
}

func (s service) Move(ctx context.Context, r MoveRequest) error {
	return nil
}

func (s service) GetCategories(ctx context.Context) ([]string, error) {
	return s.repo.GetCategories(ctx)
}

func (s service) GetCategoryContent(ctx context.Context, category_name string) ([]*entity.FileNode, error) {
	nodes, err := s.repo.GetCategoryContent(ctx, category_name)
	if nodes == nil {
		nodes = []*entity.FileNode{}
	}

	return nodes, err

}

func (s service) CreateURL(ctx context.Context, r CreateURLRequest) error {
	return nil
}
