package imgcontent

import (
	"context"
	"io"
	"net/url"
)

// ContentPathBuilder define how to Build a content path.
type ContentPathBuilder interface {
	// Build builds content-path then return it.
	Build(ctx context.Context, filename Filename) (ContentPath, error)
}

// ContentWriter define how to Write image-content.
type ContentWriter interface {
	// Write writes out the content obtained via the reader to the path.
	Write(ctx context.Context, reader io.Reader, path ContentPath) (url.URL, error)
}
