package api

import (
	"net/http"
	"strings"

	app "github.com/thatpix3l/collge_event_website/src/gen_templ"
)

// Internal helper func to build middlware using my custom handler signatures
func middlewareBuilder(path string, method string, middleware HandlerFuncMiddleware) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(StdHttpFunc(path, method, func(hs HandlerState) error {
			return middleware(hs, next)
		}))

	}
}

// Post-request cleanup middleware
var Cleanup = middlewareBuilder("/", "*", func(hs HandlerState, next http.Handler) error {
	next.ServeHTTP(hs.Local.ResponseWriter, hs.Local.Request)

	if hs.Local.Conn != nil {
		hs.Local.Conn.Release()
		hs.Local.Conn = nil
	}

	return nil
})

// Authentication middleware.
var Authentication = middlewareBuilder("/", "*", func(hs HandlerState, next http.Handler) error {

	// If accessing resource that doesn't require authentication, allow.

	// For each path...
	for path, methods := range noAuth {
		// For each method in path...
		for _, method := range methods {

			// If request's path is allowed...
			allowedPath := hs.Local.Request.URL.Path == path

			// If request's method is allowed...
			allowedMethod := hs.Local.Request.Method == strings.ToUpper(method)

			// If both are allowed...
			if allowedPath && allowedMethod {
				// Allow
				next.ServeHTTP(hs.Local.ResponseWriter, hs.Local.Request)
				return nil
			}
		}

	}

	// If not authenticated for all other resources, deny.
	if err := hs.Authenticated(); err != nil {
		hs.Local.RespondHtml(app.StatusMessage("danger", "invalid authentication token"), http.StatusBadRequest)
		return err
	}

	return nil

})
