package contentmgmt

import (
	"context"
	"github.com/pkg/errors"
	"io"
	"net/url"
)

type UploadUsecase struct {
	PathBuilder ContentPathBuilder
	Writer      ContentWriter
	Presenter   UploadResultPresenter
}

type UploadInput struct {
	Filename
	Reader io.Reader
}

type UploadOutput struct {
	Filename
	url.URL
}

func (u UploadUsecase) Handle(ctx context.Context, input UploadInput) error {

	path, err := u.PathBuilder.Build(ctx, input.Filename)
	if err != nil {
		return errors.Cause(err)
	}

	url, err := u.Writer.Write(ctx, input.Reader, path)
	if err != nil {
		return errors.Wrap(err, "failed to write content")
	}

	return u.Presenter.Complete(ctx, UploadOutput{
		Filename: input.Filename,
		URL:      url,
	})
}

type ContentPathBuilder interface {
	Build(ctx context.Context, filename Filename) (ContentPath, error)
}

type ContentWriter interface {
	Write(ctx context.Context, reader io.Reader, path ContentPath) (url.URL, error)
}

type UploadResultPresenter interface {
	Complete(ctx context.Context, data UploadOutput) error
}
