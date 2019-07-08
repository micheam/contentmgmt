package interfaces

import (
	"context"
	en "github.com/micheam/imgcontent/entities"
	uc "github.com/micheam/imgcontent/usecases"
)

var _ uc.ContentPathBuilder = (*DefaultContentPathBuilder)(nil)

type DefaultContentPathBuilder struct{}

func (c DefaultContentPathBuilder) Build(
	ctx context.Context, filename en.Filename) (path en.ContentPath, err error) {

	return
}
