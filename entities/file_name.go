package entities

import (
	"path/filepath"
	"regexp"
	"strings"
)

var regxNewline = regexp.MustCompile(`\r\n|\r|\n`) //throw panic if fail

type Filename struct {
	Value string
	Valid bool
}

func NewFilename(raw string) (fname *Filename, err error) {

	if raw == "" {
		return fname, IllegalFilename("File name must not be empty")
	}

    rawName := filepath.Base(raw)

	if regxNewline.MatchString(rawName) {
		return fname, IllegalFilename("File name must not contain new-line.")
	}

	name := strings.TrimSpace(rawName)
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "#", "_")

	return &Filename{Value: name, Valid: true}, err
}

type IllegalFilename string

func (e IllegalFilename) Error() string {
	return string(e)
}
