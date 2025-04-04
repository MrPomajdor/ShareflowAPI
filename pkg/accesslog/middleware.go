// Package accesslog provides a middleware that records every RESTful API call in a log message.
package accesslog

import (
	"net/http"
	"time"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/go-ozzo/ozzo-routing/v2/access"
	"github.com/sirupsen/logrus"
)

// Handler returns a middleware that records an access log message for every HTTP request being processed.
func Handler(logger *logrus.Logger) routing.Handler {
	return func(c *routing.Context) error {
		start := time.Now()

		rw := &access.LogResponseWriter{ResponseWriter: c.Response, Status: http.StatusOK}
		c.Response = rw

		// associate request ID and session ID with the request context
		// so that they can be added to the log messages
		ctx := c.Request.Context()

		c.Request = c.Request.WithContext(ctx)

		err := c.Next()

		// generate an access log message
		logger.WithContext(ctx).WithFields(logrus.Fields{"duration": time.Since(start).Milliseconds(), "status": rw.Status}).
			Infof("%s %s %s %d %d", c.Request.Method, c.Request.URL.Path, c.Request.Proto, rw.Status, rw.BytesWritten)

		return err
	}
}
