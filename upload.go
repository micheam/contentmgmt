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
		Exec(ctx context.Context, input UploadInput, cb UploadResultHandler) error
	}
	// UploadResultHandler define how to handle upload result data
	UploadResultHandler func(ctx context.Context, data UploadOutput) error
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
}

// NewUpload return initialized Upload usecase.
func NewUpload(PathBuilder ContentPathBuilder, Writer ContentWriter) Upload {
	return &upload{PathBuilder: PathBuilder, Writer: Writer}
}

func (u upload) Exec(ctx context.Context, input UploadInput, cb UploadResultHandler) error {
	var (
		err  error
		path ContentPath
		url  url.URL
	)
	if path, err = u.PathBuilder.Build(ctx, input.Filename); err != nil {
		return errors.Cause(err)
	}

	if url, err = u.Writer.Write(ctx, input.Reader, path); err != nil {
		return errors.Wrap(err, "failed to write content")
	}

	return cb(ctx, UploadOutput{
		Filename: input.Filename,
		URL:      url,
	})
}
