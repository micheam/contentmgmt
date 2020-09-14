package imgcontent

import (
	"context"
	"io"
	"net/url"
)

type ContentPathBuilder interface {
	Build(ctx context.Context, filename Filename) (ContentPath, error)
}

type ContentWriter interface {
	Write(ctx context.Context, reader io.Reader, path ContentPath) (url.URL, error)
}

type UploadResultPresenter interface {
	Complete(ctx context.Context, data UploadOutput) error
}
