package cloudstorage

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/micheam/imgcontent/entities"
	"github.com/micheam/imgcontent/usecases"
	"io"
	"log"
	"net/url"
)

var _ usecases.ContentWriter = (*GCPContentRepository)(nil)
var _ usecases.ContentWriter = GCPContentRepository{}

type GCPContentRepository struct{}

func (g GCPContentRepository) Write(ctx context.Context, file io.Reader, path entities.ContentPath) (url url.URL, err error) {

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// TODO: specify `bucketName` from conf file or env var.
	bucketName := "micheam-image-content"

	bucket := client.Bucket(bucketName)
	wc := bucket.Object(string(path)).NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		return
	}
	if err = wc.Close(); err != nil {
		return
	}

    url.Scheme = "https"
    url.Host = "storage.googleapis.com"
    url.Path = bucketName + "/" + string(path)

	return url, nil
}
