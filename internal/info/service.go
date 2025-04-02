package info

import (
	"context"
	"strings"

	"github.com/MrPomajdor/ShareFlowAPI/internal/auth"
	"github.com/MrPomajdor/ShareFlowAPI/internal/entity"
	"github.com/MrPomajdor/ShareFlowAPI/internal/errors"
	"github.com/MrPomajdor/ShareFlowAPI/pkg/dbcontext"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/sirupsen/logrus"
)

type Service interface {
	//Returns authorized users personal information including profile img url
	Info(ctx context.Context) interface{}
	//Update updates specified field in user information to specified value
	Update(ctx context.Context, field, value string, logger *logrus.Logger) error
}

type service struct {
	db     *dbcontext.DB
	logger *logrus.Logger
}

func NewService(logger *logrus.Logger, db *dbcontext.DB) Service {
	return service{db, logger}
}

func (s service) Info(ctx context.Context) interface{} {
	logging := logrus.WithContext(ctx)
	user := auth.CurrentUser(ctx)
	q := s.db.DB().NewQuery("SELECT * FROM users WHERE id={:id}")
	q.Bind(dbx.Params{"id": user.GetID()})
	dbUserData := entity.User{}
	if err := q.One(&dbUserData); err != nil {
		logging.WithError(err).Error("Error querying user's info from db")
		return nil
	}
	UserData := struct {
		FirstName  string
		LastName   string
		Email      string
		ProfileIMG string
	}{dbUserData.FirstName, dbUserData.LastName, dbUserData.Email, dbUserData.ProfileIMG}
	return UserData
}

func (s service) Update(ctx context.Context, field, value string, logger *logrus.Logger) error {
	logging := logger.WithContext(ctx)
	user := auth.CurrentUser(ctx)
	field = strings.ToLower(field)
	//check if the user is modifing the correct field
	var q *dbx.Query
	switch field {
	case "phone":
		q = s.db.DB().NewQuery("UPDATE users SET phone = {:value} WHERE id = {:id}")
	case "first_name":
		q = s.db.DB().NewQuery("UPDATE users SET first_name = {:value} WHERE id = {:id}")
	case "last_name":
		q = s.db.DB().NewQuery("UPDATE users SET last_name = {:value} WHERE id = {:id}")
	case "profile_img":
		q = s.db.DB().NewQuery("UPDATE users SET profile_img = {:value} WHERE id = {:id}")
	default:
		return errors.BadRequest("Illegal field")
	}

	q.Bind(dbx.Params{
		"value": value,
		"id":    user.GetID(),
	})
	if _, err := q.Execute(); err != nil {
		logging.WithError(err).Error("User data modification query error")
		return errors.InternalServerError("")
	}
	return nil

}
