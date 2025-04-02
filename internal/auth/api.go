package auth

import (
	"net/http"

	"github.com/MrPomajdor/ShareFlowAPI/internal/errors"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/sirupsen/logrus"
)

// RegisterHandlers registers handlers for different HTTP requests.
func RegisterHandlers(rg *routing.RouteGroup, service Service, logger *logrus.Logger) {
	rg.Post("/login", login(service, logger))
	rg.Post("/register", register(service, logger))
}

// login returns a handler that handles user login request.
func login(service Service, logger *logrus.Logger) routing.Handler {
	return func(c *routing.Context) error {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.Read(&req); err != nil {
			logger.WithContext(c.Request.Context()).WithField("error", err.Error()).Error("invalid request")
			return errors.BadRequest("")
		}

		token, err := service.Login(c.Request.Context(), req.Email, req.Password)
		if err != nil {
			return err
		}
		return c.Write(struct {
			Token string `json:"token"`
		}{token})
	}
}

// register returns a handler that handles user register request.
func register(service Service, logger *logrus.Logger) routing.Handler {
	return func(c *routing.Context) error {
		var req struct {
			Fname    string `json:"first_name"`
			Lname    string `json:"last_name"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Authcode string `json:"authcode"`
		}

		if err := c.Read(&req); err != nil {
			logger.WithContext(c.Request.Context()).WithField("error", err.Error()).Error("invalid request")
			return errors.BadRequest("")
		}

		err := service.Register(c.Request.Context(), req.Fname, req.Lname, req.Email, req.Password, req.Authcode)
		if err != nil {
			return errors.BadRequest(err.Error())
		}
		return c.Write(struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
		}{http.StatusOK, "Registration successfull"})
	}
}
