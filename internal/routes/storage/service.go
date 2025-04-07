package storage

import (
	"context"

	"github.com/MrPomajdor/ShareFlowAPI/internal/entity"
	"github.com/MrPomajdor/ShareFlowAPI/pkg/dbcontext"
	"github.com/sirupsen/logrus"
)

type Service interface {
	// Uploads provided file to server
	Upload(ctx context.Context, r UploadRequest) error
	// Removes a file with provided path from the server
	Remove(ctx context.Context, r RemoveRequest) error
	// Moves a file from provided path to a new destination
	Move(ctx context.Context, r MoveRequest) error
	// List provides a list of files and directories owned by the user
	GetRoot(ctx context.Context) (*entity.FileNode, error)
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

func (s service) Upload(ctx context.Context, r UploadRequest) error {
	return nil
}

func (s service) Remove(ctx context.Context, r RemoveRequest) error {
	return nil
}

func (s service) Move(ctx context.Context, r MoveRequest) error {
	return nil
}

func (s service) GetRoot(ctx context.Context) (*entity.FileNode, error) {
	// logger := s.logger.WithContext(ctx)

	root, err := s.repo.GetRoot(ctx)
	if err != nil {
		return nil, err
	}
	return root, nil
}

func (s service) CreateURL(ctx context.Context, r CreateURLRequest) error {
	return nil
}
