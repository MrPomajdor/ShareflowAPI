package storage

import (
	"strconv"

	"github.com/MrPomajdor/ShareFlowAPI/internal/errors"
	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/content"
)

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler) {
	r.Use(content.TypeNegotiator(content.JSON))
	r.Use(authHandler)
	r.Get("/storage/", ListCategories(service))
	r.Get("/storage/files", List(service))
	r.Put("/storage/files", Upload(service))
	r.Delete("/storage/files/<id>", Remove(service))
	r.Get("/storage/files/<id>", GetFile(service))
	r.Get("/storage/createurl", CreateURL(service))
}

func Upload(s Service) routing.Handler {
	return func(c *routing.Context) error {
		ctx := c.Request.Context()
		logger := s.GetLogger().WithContext(ctx)

		file, handler, err := c.Request.FormFile("file")
		if err != nil {
			logger.WithError(err).Trace("invalid form file")
			return errors.BadRequest("Invalid form file")
		}
		category := c.Request.FormValue("category")
		if category == "" {
			category = "default"
		}

		node, err := s.Upload(ctx, file, handler, category)
		if err != nil {
			return err
		}
		return c.Write(node)

	}
}

func Remove(s Service) routing.Handler {
	return func(c *routing.Context) error {
		ctx := c.Request.Context()

		param := c.Param("id")
		parsed, err_parse := strconv.ParseInt(param, 10, len(param))
		if err_parse != nil {
			return errors.BadRequest("Invalid parameter id")
		}

		return s.Remove(ctx, int(parsed))
	}
}

func List(s Service) routing.Handler {
	return func(c *routing.Context) error {
		ctx := c.Request.Context()
		loger := s.GetLogger().WithContext(ctx)

		var request ListRequest
		reader := routing.JSONDataReader{}
		if err := reader.Read(c.Request, &request); err != nil {
			loger.WithContext(c.Request.Context()).WithField("error", err.Error()).Error("invalid request")
			return errors.BadRequest("")
		}

		if err := request.Validate(); err != nil {
			loger.WithError(err).Error("asdasdasd")
			return errors.BadRequest("illegal request values")
		}

		if request.Category == "" {
			request.Category = "default"
		}
		nodes, err := s.GetCategoryContent(ctx, request.Category)

		if err != nil {
			loger.WithError(err).Error("List error")
			return errors.InternalServerError("")
		}

		return c.Write(nodes)
	}
}

func ListCategories(s Service) routing.Handler {
	return func(c *routing.Context) error {
		ctx := c.Request.Context()
		loger := s.GetLogger().WithContext(ctx)

		categories, err := s.GetCategories(ctx)

		if err != nil {
			loger.WithError(err).Error("List categories error")
			return err
		}

		return c.Write(categories)
	}
}

func CreateURL(s Service) routing.Handler {
	return func(c *routing.Context) error {
		return nil
	}
}

func GetFile(s Service) routing.Handler {
	return func(c *routing.Context) error {
		return nil
	}
}
