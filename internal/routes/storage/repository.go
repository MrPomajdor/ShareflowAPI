package storage

import (
	"context"
	"database/sql"
	"mime/multipart"
	"path"
	"slices"
	"strconv"

	"github.com/MrPomajdor/ShareFlowAPI/internal/entity"
	"github.com/MrPomajdor/ShareFlowAPI/internal/errors"
	"github.com/MrPomajdor/ShareFlowAPI/internal/filesystem"
	"github.com/MrPomajdor/ShareFlowAPI/internal/routes/auth"
	"github.com/MrPomajdor/ShareFlowAPI/pkg/dbcontext"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context/ctxhttp"
)

type Repository interface {
	// Uploads provided file to server
	Upload(ctx context.Context, r UploadRequest) error
	// Get returns the file with the selected id
	Get(ctx context.Context, id int) (*entity.FileNode, error)
	// Removes a file with provided path from the server
	Remove(ctx context.Context, r RemoveRequest) error
	// Moves a file from provided path to a new destination
	Move(ctx context.Context, r MoveRequest) error
	// GetCategories provides a list of directories created by the user.
	// Includes the default category
	GetCategories(ctx context.Context) ([]string, error)
	// GetCategoryContet provides a list of files and directories owned by the user
	GetCategoryContent(ctx context.Context, category string) ([]*entity.FileNode, error)
	// CreateURL creates and provides a link for sharing the selected file
	CreateURL(ctx context.Context, id int) (string, error)
	// Initializes new file in the database with a name and category
	// returns ID of the file in the database
	InitializeNewFile(ctx context.Context, name, category_name string) (*entity.FileNode, error)
	// WriteFile writes provided file onto the disk
	WriteFile(ctx context.Context, file multipart.File, file_struct *entity.FileNode) error
}

// repository persists files in database
type repository struct {
	db          *dbcontext.DB
	logger      *logrus.Logger
	storagePath string
}

// NewRepository creates a new file repository
func NewRepository(db *dbcontext.DB, logger *logrus.Logger, storagePath string) Repository {
	return repository{db, logger, storagePath}
}

func (r repository) Upload(ctx context.Context, req UploadRequest) error {
	return nil
}

func (r repository) Get(ctx context.Context, id int) (*entity.FileNode, error) {
	user := auth.CurrentUser(ctx)
	q := r.db.DB().Select("*").From("filesystem").Where(dbx.HashExp{"id": id, "owner_id": user.GetID()})

	var node entity.FileNode
	row, err := q.Build().Rows()
	if err != nil {
		return nil, err
	}

	if err := row.Scan(&node.ID, &node.Name, &node.OwnerID, &node.CreatedAt, &node.Size); err != nil {
		return nil, err
	}

	return &node, nil
}

func (r repository) InitializeNewFile(ctx context.Context, name, category string) (*entity.FileNode, error) {
	user := auth.CurrentUser(ctx)
	logger := r.logger.WithContext(ctx)
	q := r.db.DB().Insert("filesystem", dbx.Params{
		"name":          name,
		"owner_id":      user.GetID(),
		"category_name": category,
	})

	res, err := q.Execute()
	if err != nil {
		logger.WithError(err).Error("File initialization error")
		return nil, errors.InternalServerError("")
	}
	// TODO : LastInsertId is not supported by all databases.
	// Here we also implement a check if the selected database even has that function,
	// but for now we assume that it has.
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	new_file := entity.FileNode{}
	new_file.ID = (int)(id)
	new_file.Name = name
	new_file.OwnerID = user.GetID()
	new_file.CategoryName = category
	return &new_file, nil
}

func (r repository) WriteFile(ctx context.Context, file multipart.File, file_struct *entity.FileNode) error {
	user := auth.CurrentUser(ctx)

	filePath := path.Join(
		r.storagePath,
		strconv.FormatInt((int64)(user.GetID()), 10),
		strconv.FormatInt((int64)(file_struct.ID), 10),
	)
	return filesystem.WriteFile(file, file_struct, filePath)
}

func (r repository) Remove(ctx context.Context, req RemoveRequest) error {
	user := auth.CurrentUser(ctx)

	q := r.db.DB().Delete("filesystem", dbx.HashExp{"id": req.ID, "owner_id": auth.CurrentUser(ctx).GetID()})
	_, err := q.Execute()

	if err == sql.ErrNoRows {
		return errors.NotFound("Requested file not found")
	}
	if err != nil {
		return errors.InternalServerError("")
	}

	filePath := path.Join(
		r.storagePath,
		strconv.FormatInt((int64)(user.GetID()), 10),
		req.ID,
	)
	return filesystem.RemoveFile(filePath)
}

func (r repository) Move(ctx context.Context, req MoveRequest) error {
	return nil
	// q := r.db.DB().Delete("filesystem",dbx.HashExp{"":req.})
}

func (r repository) GetCategories(ctx context.Context) ([]string, error) {
	loger := r.logger.WithContext(ctx)
	user := auth.CurrentUser(ctx)

	q := r.db.DB().Select("category_name").From("filesystem").Where(dbx.HashExp{"owner_id": user.GetID()}).Build()

	rows, err := q.Rows()
	var categories []string
	categories = append(categories, "default")
	if err != nil && err == sql.ErrNoRows {
		return categories, nil
	}
	if err != nil {
		loger.WithError(err).Error("GetCategories error")
		return nil, errors.InternalServerError("")
	}

	for rows.Next() {
		var cat string
		rows.Scan(&cat)
		if !slices.Contains(categories, cat) {
			categories = append(categories, cat)
		}
	}

	return categories, nil
}

func (r repository) GetCategoryContent(ctx context.Context, category_name string) ([]*entity.FileNode, error) {
	loger := r.logger.WithContext(ctx)
	user := auth.CurrentUser(ctx)

	q := r.db.DB().Select("*").From("filesystem").Where(dbx.HashExp{"owner_id": user.GetID(), "category_name": category_name}).Build()

	var nodes []*entity.FileNode

	rows, err := q.Rows()

	if err == sql.ErrNoRows {
		return nil, errors.NotFound("category not found")
	}

	if err != nil {
		loger.WithError(err).Error("Get category sql error")
		return nil, errors.InternalServerError("")

	}

	for rows.Next() {
		var nd entity.FileNode
		rows.ScanStruct(&nd)
		nodes = append(nodes, &nd)
	}
	return nodes, nil
}

func (r repository) CreateURL(ctx context.Context, id int) (string, error) {
	return "", nil
}

func (r repository) CheckUserStoragePath(ctx context.Context) {
	user := auth.CurrentUser(ctx)
	userDataPath := path.Join(
		r.storagePath,
		strconv.FormatInt((int64)(user.GetID()), 10),
	)
	if !filesystem.Exists(userDataPath) {
		filesystem.CreateDirectory(userDataPath)
	}
}
