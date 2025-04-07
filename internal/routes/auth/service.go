package auth

import (
	"context"
	"time"

	"github.com/MrPomajdor/ShareFlowAPI/internal/crypt"
	"github.com/MrPomajdor/ShareFlowAPI/internal/entity"
	"github.com/MrPomajdor/ShareFlowAPI/internal/errors"
	"github.com/MrPomajdor/ShareFlowAPI/pkg/dbcontext"
	"github.com/dgrijalva/jwt-go"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/sirupsen/logrus"
)

// Service encapsulates the authentication logic.
type Service interface {
	// authenticate authenticates a user using username and password.
	// It returns a JWT token if authentication succeeds. Otherwise, an error is returned.
	Login(ctx context.Context, req LoginRequest) (string, error)
	// Register registers a user using full name, email, password and AuthCode
	// An error is returned if the registration does not succeed.
	Register(ctx context.Context, req RegisterRequest) error
	// Returns service logger
	GetLogger() *logrus.Logger
}

type service struct {
	signingKey      string
	tokenExpiration int
	database        *dbcontext.DB
	logger          *logrus.Logger
}

func (s service) GetLogger() *logrus.Logger {
	return s.logger
}

// NewService creates a new authentication service.
func NewService(signingKey string, tokenExpiration int, db *dbcontext.DB, logger *logrus.Logger) Service {
	return service{signingKey, tokenExpiration, db, logger}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, req LoginRequest) (string, error) {
	identity, err := s.authenticate(ctx, req)
	if identity != nil {
		return s.generateJWT(identity)
	}
	return "", errors.Unauthorized(err.Error())
}

// Register creates a user
func (s service) Register(ctx context.Context, req RegisterRequest) error {
	return s.register(ctx, req)
}

// authenticate authenticates a user using username and password.
// If username and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, req LoginRequest) (entity.Identity, error) {
	if err := req.Validate(); err != nil {
		return nil, errors.BadRequest("")
	}
	logger := s.logger.WithContext(ctx).WithField("user", req.Email)
	hashed, err := crypt.HashPassword(req.Password)
	if err != nil {
		logger.Error("Failed to hash password")
		return nil, errors.BadRequest("failed to hash password")
	}
	q := s.database.With(ctx).NewQuery("SELECT * FROM users WHERE email={:email}")
	q.Bind(dbx.Params{
		"email":    req.Email,
		"password": hashed,
	})
	User := entity.User{}
	rowErr := q.One(&User) //q.Row(&User.ID, &User.Email, &User.HashedPassword, &User.FirstName, &User.LastName, &User.ProfileIMG, &User.AuthCode, &User.CreatedAt, &User.LastLogin, &User.LastLoginIP)
	if rowErr != nil {
		logger.WithError(rowErr).Trace("Row error")
		return nil, errors.Unauthorized("invalid email or password")
	}

	if !crypt.CheckPasswordHash(req.Password, User.HashedPassword) {
		logger.Trace("Invaid password for valid email")
		return nil, errors.Unauthorized("invalid email or password")
	}

	return User, nil

}

// Register registers a user using full name, email, password and AuthCode
// Returns an error if registration fails
func (s service) register(ctx context.Context, req RegisterRequest) error {
	if err := req.Validate(); err != nil {
		return errors.BadRequest("")
	}
	logger := s.logger.WithContext(ctx).WithField("user", req.Email)
	var authcode_db struct {
		ID             int
		Authcode       string
		Email          string
		CreationDate   string
		ExpirationDate string
		Used           bool
	}
	q := s.database.With(ctx).NewQuery("SELECT * FROM authcodes WHERE authcode={:authcode} AND email={:email} AND used=0")
	q.Bind(dbx.Params{
		"authcode": req.Authcode,
		"email":    req.Email,
	})
	qErr := q.Row(&authcode_db.ID, &authcode_db.Authcode, &authcode_db.Email, &authcode_db.CreationDate, &authcode_db.ExpirationDate, &authcode_db.Used)
	if qErr != nil {
		logger.WithFields(logrus.Fields{"reason": "invalid authcode", "error": qErr}).Error("user creation failed")
		return errors.BadRequest("registration failed - invalid authcode")
	}

	q2 := s.database.With(ctx).NewQuery("SELECT COUNT(*) FROM users WHERE email={:email}")
	q2.Bind(dbx.Params{
		"email": req.Email,
	})
	var count int
	q2.Row(&count)
	if count != 0 {
		logger.Error("Account already exists")
		return errors.BadRequest("account already exists")
	}

	hashed, hash_err := crypt.HashPassword(req.Password)
	if hash_err != nil {
		logger.Error("Failed to hash password")
		return errors.InternalServerError("failed to hash password")
	}

	q3 := s.database.With(ctx).NewQuery("INSERT INTO `users`(`email`, `password`, `first_name`, `last_name`, `auth_code`) VALUES ({:email},{:password},{:first_name},{:last_name},{:auth_code})")
	q3.Bind(dbx.Params{
		"email":      req.Email,
		"password":   hashed,
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"auth_code":  req.Authcode,
	})
	_, err := q3.Execute()

	q4 := s.database.With(ctx).NewQuery("UPDATE authcodes SET used = 1 WHERE authcode = {:authcode}")
	q4.Bind(dbx.Params{"authcode": req.Authcode})
	_, _ = q4.Execute() // There is no chance this won't succed if previous queries succeded.

	if err == nil {
		return nil
	}

	logger.WithField("reason", "not sure").Error("user creation failed")
	return errors.InternalServerError("registration failed")
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(identity entity.Identity) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        identity.GetID(),
		"firstname": identity.GetFirstName(),
		"lastname":  identity.GetLastName(),
		"email":     identity.GetEmail(),
		"exp":       time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(s.signingKey))
}
