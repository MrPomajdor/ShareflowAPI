package auth

import (
	"context"
	"fmt"
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
	Login(ctx context.Context, email, password string) (string, error)
	// Register registers a user using full name, email, password and AuthCode
	// An error is returned if the registration does not succeed.
	Register(ctx context.Context, fname, lname, email, password, authcode string) error
}

// Identity represents an authenticated user identity.
type Identity interface {
	// GetID returns the user ID.
	GetID() int
	// GetFirstName returns the user first name.
	GetFirstName() string
	// GetLastName returns the user last name.
	GetLastName() string
	// GetName returns the user email.
	GetEmail() string
}

type service struct {
	signingKey      string
	tokenExpiration int
	database        *dbcontext.DB
	logger          *logrus.Logger
}

// NewService creates a new authentication service.
func NewService(signingKey string, tokenExpiration int, db *dbcontext.DB, logger *logrus.Logger) Service {
	return service{signingKey, tokenExpiration, db, logger}
}

// Login authenticates a user and generates a JWT token if authentication succeeds.
// Otherwise, an error is returned.
func (s service) Login(ctx context.Context, username, password string) (string, error) {
	identity, err := s.authenticate(ctx, username, password)
	if identity != nil {
		return s.generateJWT(identity)
	}
	return "", errors.Unauthorized(err.Error())
}

// Register creates a user
func (s service) Register(ctx context.Context, fname, lname, email, password, authcode string) error {
	return s.register(ctx, fname, lname, email, password, authcode)
}

// authenticate authenticates a user using username and password.
// If username and password are correct, an identity is returned. Otherwise, nil is returned.
func (s service) authenticate(ctx context.Context, email, password string) (Identity, error) {
	logger := s.logger.WithContext(ctx).WithField("user", email)
	hashed, err := crypt.HashPassword(password)
	if err != nil {
		logger.Error("Failed to hash password")
		return nil, fmt.Errorf("failed to hash password")
	}
	q := s.database.With(ctx).NewQuery("SELECT * FROM users WHERE email={:email}")
	q.Bind(dbx.Params{
		"email":    email,
		"password": hashed,
	})
	User := entity.User{}
	rowErr := q.Row(&User.ID, &User.Email, &User.HashedPassword, &User.FirstName, &User.LastName, &User.AuthCode, &User.CreatedAt, &User.LastLogin, &User.LastLoginIP)
	if rowErr != nil {
		logger.WithError(rowErr).Trace("Row error")
		return nil, fmt.Errorf("invalid email or password")
	}
	if !crypt.CheckPasswordHash(password, User.HashedPassword) {
		return nil, fmt.Errorf("invalid email or password")
	}

	return User, nil

}

// Register registers a user using full name, email, password and AuthCode
// Returns an error if registration fails
func (s service) register(ctx context.Context, fname, lname, email, password, authcode string) error {
	logger := s.logger.WithContext(ctx).WithField("user", email)
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
		"authcode": authcode,
		"email":    email,
	})
	qErr := q.Row(&authcode_db.ID, &authcode_db.Authcode, &authcode_db.Email, &authcode_db.CreationDate, &authcode_db.ExpirationDate, &authcode_db.Used)
	if qErr != nil {
		logger.WithFields(logrus.Fields{"reason": "invalid authcode", "error": qErr}).Error("user creation failed")
		return fmt.Errorf("registration failed - invalid authcode")
	}

	q2 := s.database.With(ctx).NewQuery("SELECT COUNT(*) FROM users WHERE email={:email}")
	q2.Bind(dbx.Params{
		"email": email,
	})
	var count int
	q2.Row(&count)
	if count != 0 {
		logger.Error("Account already exists")
		return fmt.Errorf("account already exists")
	}

	hashed, hash_err := crypt.HashPassword(password)
	if hash_err != nil {
		logger.Error("Failed to hash password")
		return fmt.Errorf("failed to hash password")
	}

	q3 := s.database.With(ctx).NewQuery("INSERT INTO `users`(`email`, `password`, `first_name`, `last_name`, `auth_code`) VALUES ({:email},{:password},{:first_name},{:last_name},{:auth_code})")
	q3.Bind(dbx.Params{
		"email":      email,
		"password":   hashed,
		"first_name": fname,
		"last_name":  lname,
		"auth_code":  authcode,
	})
	_, err := q3.Execute()
	if err == nil {
		return nil
	}
	logger.WithField("reason", "not sure").Error("user creation failed")
	return fmt.Errorf("registration failed")
}

// generateJWT generates a JWT that encodes an identity.
func (s service) generateJWT(identity Identity) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       identity.GetID(),
		"firsname": identity.GetFirstName(),
		"lastname": identity.GetLastName(),
		"email":    identity.GetEmail(),
		"exp":      time.Now().Add(time.Duration(s.tokenExpiration) * time.Hour).Unix(),
	}).SignedString([]byte(s.signingKey))
}
