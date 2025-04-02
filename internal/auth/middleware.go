package auth

import (
	"context"

	"github.com/MrPomajdor/ShareFlowAPI/internal/entity"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/auth"
	"github.com/golang-jwt/jwt"
)

// Handler returns a JWT-based authentication middleware.
func Handler(verificationKey string) routing.Handler {
	return auth.JWT(verificationKey, auth.JWTOptions{TokenHandler: handleToken})
}

// handleToken stores the user identity in the request context so that it can be accessed elsewhere.
func handleToken(c *routing.Context, token *jwt.Token) error {
	ctx := WithUser(
		c.Request.Context(),
		token.Claims.(jwt.MapClaims)["id"].(int),
		token.Claims.(jwt.MapClaims)["firstname"].(string),
		token.Claims.(jwt.MapClaims)["lastname"].(string),
		token.Claims.(jwt.MapClaims)["email"].(string),
	)
	c.Request = c.Request.WithContext(ctx)
	return nil
}

type contextKey int

const (
	userKey contextKey = iota
)

// WithUser returns a context that contains the user identity from the given JWT.
func WithUser(ctx context.Context, id int, firstnamename, lastname, email string) context.Context {
	return context.WithValue(ctx, userKey, entity.User{ID: id, FirstName: firstnamename, LastName: lastname, Email: email})
}

// CurrentUser returns the user identity from the given context.
// Nil is returned if no user identity is found in the context.
func CurrentUser(ctx context.Context) Identity {
	if user, ok := ctx.Value(userKey).(entity.User); ok {
		return user
	}
	return nil
}
