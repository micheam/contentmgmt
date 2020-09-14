package gcs

import (
	"io"
	"net/url"

	"context"
	"errors"
	"os"

	"cloud.google.com/go/storage"
	"github.com/micheam/imgcontent"
	"google.golang.org/api/option"
)

// NewClient creates a new Google Cloud Storage client.
func NewClient(ctx context.Context) (*storage.Client, error) {
	cred := os.Getenv("IMGCONTENT_GCS_CREDENTIALS")
	if cred == "" {
		return nil, errors.New("env IMGCONTENT_GCS_CREDENTIALS is not set")
	}
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(cred))
	if err != nil {
		return nil, err
	}
	return client, nil
}

var _ imgcontent.ContentWriter = (*ContentRepository)(nil)

// ContentRepository ...
type ContentRepository struct {
	BucketName string
	Client     *storage.Client
}

// NewContentWriter init ContentWriter then return it.
func NewContentWriter(bucketName string, client *storage.Client) imgcontent.ContentWriter {
	return &ContentRepository{BucketName: bucketName, Client: client}
}

// Write ...
func (g *ContentRepository) Write(ctx context.Context, file io.Reader, path imgcontent.ContentPath) (url url.URL, err error) {

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
