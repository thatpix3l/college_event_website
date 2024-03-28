package utils

import (
	"errors"
	"strings"

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

// Shorthand to prepend given error with multiple string errors.
func ErrPrep(err error, msgs ...string) error {
	for i := len(msgs) - 1; i >= 0; i-- {
		err = errors.Join(errors.New(msgs[i]), err)
	}
	return err
}

// Make the first character uppercase.
func ToUpperFirst(s string) string {
	sBytes := []byte(s)
	sBytes[0] = []byte(strings.ToUpper(s))[0]
	return string(sBytes)
}
