package interfaces

import (
	"context"
	e "github.com/micheam/imgcontent/entities"
	uc "github.com/micheam/imgcontent/usecases"
)

var _ uc.ContentPathBuilder = (*DefaultContentPathBuilder)(nil)

type DefaultContentPathBuilder struct{}

func (c DefaultContentPathBuilder) Build(
	ctx context.Context, input uc.BuildContentPathInput) (path e.ContentPath, err error) {

	return
}
