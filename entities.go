package imgcontent

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

// Path is a path of image content.
type Path string

var regxNewline = regexp.MustCompile(`\r\n|\r|\n`) //throw panic if fail

// Name is a name of image content.
type Name struct {
	Value string
	Valid bool
}

// NewName create new content-name.
// this will return error if specified raw is invalid.
func NewName(raw string) (fname *Name, err error) {

	if raw == "" {
		return fname, ErrIllegalFilename("file name must not be empty")
	}

	rawName := filepath.Base(raw)

	if regxNewline.MatchString(rawName) {
		return fname, ErrIllegalFilename("file name must not contain new-line.")
	}

	name := strings.TrimSpace(rawName)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "#", "_")

	return &Name{Value: name, Valid: true}, err
}

// ImageContent ...
type ImageContent struct {
	Name Name
	Path Path
	url.URL
}
