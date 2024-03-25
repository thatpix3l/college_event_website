package api

import (
	"net/http"

	app "github.com/thatpix3l/collge_event_website/src/gen_templ"
)

// Authentication middleware
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(StdHttpFunc("/", "*", func(hs HandlerState) error {

		next.ServeHTTP(hs.Local.ResponseWriter, hs.Local.Request)
		return nil

		// allow if accessing any resources that don't require authentication or authorization
		for _, path := range noAuthPaths {
			if hs.Local.Request.URL.Path == path {
				next.ServeHTTP(hs.Local.ResponseWriter, hs.Local.Request)
				return nil
			}
		}

		// Should only be here if requested resource requires authentication and authorization

		// Exit early if no authentication token provided
		givenToken, err := hs.Local.Request.Cookie("authentication_token")
		if err == http.ErrNoCookie {

			// alert user that they are not authenticated yet
			hs.Local.ResponseWriter.WriteHeader(http.StatusUnauthorized)
			if _, err := hs.Local.ResponseWriter.Write([]byte("401 - Any access besides homepage requires authentication")); err != nil {
				return err
			}
			return err

		}

		// Exit Early if given token is invalid
		if !validToken(givenToken.Value) {
			if err := hs.Local.RespondHtml(app.StatusMessage("danger", "403 - Not authorized to access resource"), http.StatusForbidden); err != nil {
				return err
			}
		}

		// Token exists and is valid, continue
		next.ServeHTTP(hs.Local.ResponseWriter, hs.Local.Request)

		return nil

	}))
}
