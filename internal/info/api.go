package info

import (
	"github.com/MrPomajdor/ShareFlowAPI/internal/errors"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/sirupsen/logrus"
)

type resource struct {
	service Service
	logger  *logrus.Logger
}

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger *logrus.Logger) {
	res := resource{service, logger}
	r.Use(authHandler)
	r.Get("/me/info", res.Info(service, logger))
	r.Post("/me/update", res.Update(service, logger))
}

func (r resource) Info(s Service, logger *logrus.Logger) routing.Handler {
	return func(c *routing.Context) error {
		userData := r.service.Info(c.Request.Context())
		if userData == nil {
			return errors.InternalServerError("")
		}
		return c.Write(userData)
	}
}

func (r resource) Update(s Service, logger *logrus.Logger) routing.Handler {
	return func(c *routing.Context) error {
		var req struct {
			Field string `json:"field"`
			Value string `json:"value"`
		}

		if err := c.Read(&req); err != nil {
			logger.WithContext(c.Request.Context()).WithField("error", err.Error()).Error("invalid request")
			return errors.BadRequest("")
		}
		if err := s.Update(c.Request.Context(), req.Field, req.Value, logger); err != nil {
			return err
		}
		return nil
	}
}
