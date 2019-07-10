package console

import (
	"context"
	en "github.com/micheam/imgcontent/entities"
	uc "github.com/micheam/imgcontent/usecases"
	"strings"
	"time"
)

var _ uc.ContentPathBuilder = (*DefaultContentPathBuilder)(nil)

type DefaultContentPathBuilder struct {
	BaseDate time.Time
}

func (c DefaultContentPathBuilder) Build(
	ctx context.Context, filename en.Filename) (path en.ContentPath, err error) {

	prefix := c.BaseDate.Format("2006/01/02")
	path = en.ContentPath(strings.Join([]string{prefix, filename.Value}, "/"))
	return
}
