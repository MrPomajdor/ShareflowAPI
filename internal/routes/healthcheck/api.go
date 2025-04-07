package healthcheck

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
)

// RegisterHandlers registers the handlers that perform healthchecks.
func RegisterHandlers(r *routing.Router, version string) {
	r.To("GET,HEAD", "/healthcheck", healthcheck(version))
	r.To("GET,HEAD", "/mentalcheck", mentalcheck())
}

// healthcheck responds to a healthcheck request.
func healthcheck(version string) routing.Handler {
	return func(c *routing.Context) error {
		return c.Write("OK " + version)
	}
}

func mentalcheck() routing.Handler {
	return func(c *routing.Context) error {
		resp := struct {
			Status      string  `json:"status"`
			Comments    string  `json:"comments"`
			Rating      float64 `json:"rating"`
			RatingScale string  `json:"rating_scale"`
		}{"on_the_edge", "Almost insane but still not quite", 3.5, "0,10"}
		return c.WriteWithStatus(resp, http.StatusTeapot)
	}
}
