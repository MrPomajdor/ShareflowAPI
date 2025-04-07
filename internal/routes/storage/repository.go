package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/MrPomajdor/ShareFlowAPI/internal/entity"
	"github.com/MrPomajdor/ShareFlowAPI/internal/errors"
	"github.com/MrPomajdor/ShareFlowAPI/internal/routes/auth"
	"github.com/MrPomajdor/ShareFlowAPI/pkg/dbcontext"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/sirupsen/logrus"
)

type Repository interface {
	// Uploads provided file to server
	Upload(ctx context.Context, r UploadRequest) error
	// Get returns the file with the selected id
	Get(ctx context.Context, r GetRequest) (*entity.FileNode, error)
	// Removes a file with provided path from the server
	Remove(ctx context.Context, r RemoveRequest) error
	// Moves a file from provided path to a new destination
	Move(ctx context.Context, r MoveRequest) error
	// List provides a list of files and directories owned by the user
	GetRoot(ctx context.Context) (*entity.FileNode, error)
	// CreateURL creates and provides a link for sharing the selected file
	CreateURL(ctx context.Context, r CreateURLRequest) (string, error)
}

// repository persists files in database
type repository struct {
	db     *dbcontext.DB
	logger *logrus.Logger
}

// NewRepository creates a new file repository
func NewRepository(db *dbcontext.DB, logger *logrus.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Upload(ctx context.Context, req UploadRequest) error {
	return nil
}

func (r repository) Get(ctx context.Context, req GetRequest) (*entity.FileNode, error) {
	user := auth.CurrentUser(ctx)
	q := r.db.DB().Select("*").From("filesystem").Where(dbx.HashExp{"id": req.ID, "owner_id": user.GetID()})

	var node entity.FileNode
	row, err := q.Build().Rows()
	if err != nil {
		return nil, err
	}

	if err := row.Scan(&node.ID, &node.Name, &node.IsDir, &node.Parent_ID, &node.OwnerID, &node.CreatedAt); err != nil {
		return nil, err
	}

	node.Children, _ = r.getChildren(ctx, node.OwnerID)

	return &node, nil
}

func (r repository) getChildren(ctx context.Context, parent_id int) ([]*entity.FileNode, error) {
	q := r.db.DB().Select("*").From("filesystem").Where(dbx.HashExp{"parent_id": parent_id})

	rows, err := q.Build().Rows()
	if err != nil {
		return make([]*entity.FileNode, 0), err
	}

	nodes := make([]*entity.FileNode, 0)
	defer rows.Close()
	for rows.Next() {
		var node entity.FileNode
		if err := rows.Scan(&node.ID, &node.Name, &node.IsDir, &node.Parent_ID, &node.OwnerID, &node.CreatedAt); err != nil {
			return nil, err
		}

		node.Children, _ = r.getChildren(ctx, node.OwnerID)
		nodes = append(nodes, &node)
	}

	return nodes, nil
}

func (r repository) Remove(ctx context.Context, req RemoveRequest) error {
	q := r.db.DB().Delete("filesystem", dbx.HashExp{"id": req.ID, "owner_id": auth.CurrentUser(ctx).GetID()})
	_, err := q.Execute()

	if err == sql.ErrNoRows {
		return errors.NotFound("Requested file not found")
	}
	if err != nil {
		return errors.InternalServerError("")
	}

	return nil
}

func (r repository) Move(ctx context.Context, req MoveRequest) error {
	return nil
	// q := r.db.DB().Delete("filesystem",dbx.HashExp{"":req.})
}

func (r repository) GetRoot(ctx context.Context) (*entity.FileNode, error) {
	loger := r.logger.WithContext(ctx)
	user := auth.CurrentUser(ctx)
	q := r.db.DB().Select("*").From("filesystem").Where(dbx.HashExp{"parent_id": nil, "owner_id": user.GetID(), "name": "root"}).Build()

	var node entity.FileNode

	// Users root filesystem node should be created the moment he registeres,
	// but in a scenario where that doesn't happen, we create the root node now,
	// and throw an error 500
	err := q.Row(&node.ID, &node.Name, &node.IsDir, &node.Parent_ID, &node.OwnerID, &node.CreatedAt)

	if err != nil && err == sql.ErrNoRows {
		loger.WithError(err).Error("GetRoot error")
		r.CreateRoot(ctx)
		return nil, fmt.Errorf("no root")
	}
	if err != nil {
		loger.WithError(err).Error("GetRoot error")
		return nil, err
	}

	node.Children, _ = r.getChildren(ctx, node.OwnerID)

	return &node, nil
}

func (r repository) CreateURL(ctx context.Context, req CreateURLRequest) (string, error) {
	return "", nil
}

func (r repository) CreateRoot(ctx context.Context) error {
	user := auth.CurrentUser(ctx)

	q := r.db.DB().Insert("filesystem", dbx.Params{
		"name":      "root",
		"is_dir":    true,
		"parent_id": nil,
		"owner_id":  user.GetID(),
	})

	_, err := q.Execute()
	return err
}
