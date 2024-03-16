package utils

import (
	"errors"

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

// Chain one error with a new error containing custom message
func ErrInfo(err error, msg string) error {
	return errors.Join(errors.New(msg), err)
}

// Chain database error with no acquisition error
func ErrDb(err error) error {
	return ErrInfo(err, "unable to acquire database connection")
}
