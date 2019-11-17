package contentmgmt

import (
	"cloud.google.com/go/storage"
	"context"
	"io"
	"net/url"
)

var _ ContentWriter = (*GCPContentRepository)(nil)
var _ ContentWriter = GCPContentRepository{}

type GCPContentRepository struct {
	BucketName string
	Client     *storage.Client
}

func (g GCPContentRepository) Write(ctx context.Context, file io.Reader, path ContentPath) (url url.URL, err error) {

	bucket := g.Client.Bucket(g.BucketName)
	wc := bucket.Object(string(path)).NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		return
	}
	if err = wc.Close(); err != nil {
		return
	}

	url.Scheme = "https"
	url.Host = "storage.googleapis.com"
	url.Path = g.BucketName + "/" + string(path)

	return url, nil
}
