package api

import (
	"net/http"
	"strings"

	app "github.com/thatpix3l/collge_event_website/src/gen_templ"
)

// Authentication middleware
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(StdHttpFunc("/", "*", func(hs HandlerState) error {

		// If accessing resource that doesn't require authentication, allow
		for path, methods := range noAuth {
			for _, method := range methods {

				// Request path allowed
				allowedPath := hs.Local.Request.URL.Path == path

				// Request method allowed
				allowedMethod := hs.Local.Request.Method == strings.ToUpper(method)

				if allowedPath && allowedMethod {
					next.ServeHTTP(hs.Local.ResponseWriter, hs.Local.Request)
					return nil
				}
			}

		}

		// If not authenticated for all other resources, deny
		if err := hs.Authenticated(); err != nil {
			hs.Local.RespondHtml(app.StatusMessage("danger", "invalid authentication token"), http.StatusBadRequest)
			return err
		}

		// Should only be here if authenticated
		next.ServeHTTP(hs.Local.ResponseWriter, hs.Local.Request)

		return nil

	}))
}
