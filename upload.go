package imgcontent

import (
	"context"
	"io"
	"net/url"

	"github.com/pkg/errors"
)

type (
	// Upload ...
	Upload interface {
		Exec(ctx context.Context, input UploadInput) error
	}
	// UploadInput ...
	UploadInput struct {
		Filename
		Reader io.Reader
	}
	// UploadOutput ...
	UploadOutput struct {
		Filename
		url.URL
	}
)

type upload struct {
	PathBuilder ContentPathBuilder
	Writer      ContentWriter
	Presenter   UploadResultPresenter
}

// NewUpload return initialized Upload usecase.
func NewUpload(
	PathBuilder ContentPathBuilder,
	Writer ContentWriter,
	Presenter UploadResultPresenter, // TODO(micheam): CallBack であることを明確にする
) Upload {
	return &upload{
		PathBuilder: PathBuilder,
		Writer:      Writer,
		Presenter:   Presenter,
	}
}

func (u upload) Exec(ctx context.Context, input UploadInput) error {

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
