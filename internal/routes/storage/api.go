package storage

import (
	"github.com/MrPomajdor/ShareFlowAPI/internal/errors"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler) {
	r.Use(authHandler)
	r.Put("/storage/upload", Upload(service))
	r.Post("/storage/remove", Remove(service))
	r.Post("/storage/move", Move(service))
	r.Get("/storage/list", List(service))
	r.Get("/storage/createurl", CreateURL(service))
}

func Upload(s Service) routing.Handler {
	return func(c *routing.Context) error {
		ctx := c.Request.Context()
		logger := s.GetLogger().WithContext(ctx)
		var req UploadRequest
		if err := c.Read(&req); err != nil {
			return errors.InternalServerError()
		}
		file, handler, err := c.Request.FormFile("file")
		if err != nil {
			logger.WithError(err).Trace("Upload error")
			return errors.InternalServerError("")
		}
		defer file.Close()
		//filePath = s.UploadFilePath +

	}
}

func Remove(s Service) routing.Handler {
	return func(c *routing.Context) error {
		ctx := c.Request.Context()
		logger := s.GetLogger().WithContext(ctx)
		var req RemoveRequest

		if err := c.Read(&req); err != nil {
			logger.WithError(err).Trace("Bad Remove request")
			return errors.BadRequest("")
		}

		if err := req.Validate(); err != nil {
			logger.WithError(err).Trace("Remove request did not validate")
			return errors.BadRequest("invalid request values")
		}

		return s.Remove(ctx, req)
	}
}

func Move(s Service) routing.Handler {
	return func(c *routing.Context) error {
		return nil
	}
}

func List(s Service) routing.Handler {
	return func(c *routing.Context) error {
		ctx := c.Request.Context()
		loger := s.GetLogger().WithContext(ctx)
		rootFileNode, err := s.GetRoot(ctx)
		if err != nil {
			loger.WithError(err).Error("List error")
			return errors.InternalServerError("")
		}
		if _, err := rootFileNode.ToJSON(); err == nil {
			return c.Write(rootFileNode)
		} else {
			loger.WithError(err).Error("List error")
			return errors.InternalServerError("")
		}
	}
}

func CreateURL(s Service) routing.Handler {
	return func(c *routing.Context) error {
		return nil
	}
}
