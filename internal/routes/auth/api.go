package auth

import (
	"net/http"

	"github.com/MrPomajdor/ShareFlowAPI/internal/errors"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

// RegisterHandlers registers handlers for different HTTP requests.
func RegisterHandlers(rg *routing.RouteGroup, service Service) {
	rg.Post("/login", login(service))
	rg.Post("/register", register(service))
}

// login returns a handler that handles user login request.
func login(s Service) routing.Handler {
	return func(c *routing.Context) error {
		var request LoginRequest
		if err := c.Read(&request); err != nil {
			s.GetLogger().WithContext(c.Request.Context()).WithField("error", err.Error()).Error("invalid request")
			return errors.BadRequest("")
		}

		if err := request.Validate(); err != nil {
			return errors.BadRequest("invalid request values")
		}

		token, err := s.Login(c.Request.Context(), request)
		if err != nil {
			return err
		}
		return c.Write(struct {
			Token string `json:"token"`
		}{token})
	}
}

// register returns a handler that handles user register request.
func register(s Service) routing.Handler {
	return func(c *routing.Context) error {
		var req RegisterRequest
		if err := c.Read(&req); err != nil {
			s.GetLogger().WithContext(c.Request.Context()).WithField("error", err.Error()).Error("invalid request")
			return errors.BadRequest("")
		}
		if err := req.Validate(); err != nil {
			return errors.BadRequest("invalid request values")
		}
		err := s.Register(c.Request.Context(), req)
		if err != nil {
			return err
		}
		return c.WriteWithStatus(struct {
			Message string `json:"message"`
		}{"Registration successfull"}, http.StatusOK)
	}
}
