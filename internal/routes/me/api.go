package info

import (
	"github.com/MrPomajdor/ShareFlowAPI/internal/errors"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler) {
	r.Use(authHandler)
	r.Get("/me/info", Info(service))
	r.Post("/me/update", Update(service))
}

func Info(s Service) routing.Handler {
	return func(c *routing.Context) error {
		userData := s.Info(c.Request.Context())
		if userData == nil {
			return errors.InternalServerError("")
		}
		return c.Write(userData)
	}
}

func Update(s Service) routing.Handler {
	return func(c *routing.Context) error {
		var req UpdateDataRequest

		if err := c.Read(&req); err != nil {
			s.GetLogger().WithContext(c.Request.Context()).WithField("error", err.Error()).Error("invalid request")
			return errors.BadRequest("")
		}
		if err := s.Update(c.Request.Context(), req, s.GetLogger()); err != nil {
			return err
		}
		return nil
	}
}
