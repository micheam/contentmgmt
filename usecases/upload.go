package usecases

import (
	"context"
	"github.com/micheam/imgcontent/entities"
	"github.com/pkg/errors"
	"io"
	"net/url"
)

type Upload struct {
	ContentPathBuilder
	ContentWriter
	UploadPresenter
}

type UploadInput struct {
	entities.Filename
	Reader io.Reader
}

type UploadOutput struct {
	entities.Filename
	url.URL
}

func (u Upload) Handle(ctx context.Context, input UploadInput) error {

	path, err := u.ContentPathBuilder.Build(ctx, input.Filename)
	if err != nil {
		return errors.Cause(err)
	}

	url, err := u.ContentWriter.Write(ctx, input.Reader, path)
    if err != nil {
        return errors.Wrap(err, "Failed to Write content")
    }

	return u.UploadPresenter.Complete(ctx, UploadOutput{
		Filename: input.Filename,
		URL:      url,
	})
}

type ContentPathBuilder interface {
	Build(ctx context.Context, filename entities.Filename) (entities.ContentPath, error)
}

type ContentWriter interface {
	Write(ctx context.Context, reader io.Reader, path entities.ContentPath) (url.URL, error)
}

type UploadPresenter interface {
	Complete(ctx context.Context, data UploadOutput) error
}
