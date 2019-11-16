package imgcontent

import (
	"path/filepath"
	"regexp"
	"strings"
)

type ContentPath string

var regxNewline = regexp.MustCompile(`\r\n|\r|\n`) //throw panic if fail

type Filename struct {
	Value string
	Valid bool
}

func NewFilename(raw string) (fname *Filename, err error) {

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

	return &Filename{Value: name, Valid: true}, err
}

type ErrIllegalFilename string

func (e ErrIllegalFilename) Error() string {
	return string(e)
}
