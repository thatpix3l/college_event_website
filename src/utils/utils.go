package utils

import (
	"github.com/a-h/templ"
)

// Given `subpath`, prefix with API path
func ApiPath(subpath string) string {
	return "/api/v1/" + subpath
}

// Same thing as `ApiPath`, but sanitized for URLs
func ApiPathSafe(subpath string) templ.SafeURL {
	return templ.SafeURL(ApiPath(subpath))
}
