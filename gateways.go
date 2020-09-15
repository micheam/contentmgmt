package imgcontent

import (
	"context"
	"io"
	"net/url"
)

// ContentPathBuilder define how to Build a content path.
type ContentPathBuilder interface {
	// Build builds content-path then return it.
	Build(ctx context.Context, filename Name) (Path, error)
}

// ContentWriter define how to Write image-content.
type ContentWriter interface {
	// Write writes out the content obtained via the reader to the path.
	// TODO(micheam): change return value to *imgcontent.ImageContent
	Write(ctx context.Context, reader io.Reader, path Path) (url.URL, error)
}

// ContentReader define how to Read image-content.
type ContentReader interface {
	List(ctx context.Context, prefix string) (<-chan ImageContent, error)
}
