package gcs

import (
	"io"
	"log"
	"net/url"
	"strings"

	"context"
	"errors"
	"os"

	"cloud.google.com/go/storage"
	"github.com/micheam/imgcontent"
	"google.golang.org/api/iterator"
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

// BucketName return target bucket name
// TODO(micheam): package conf ?
func BucketName(ctx context.Context) string {
	return os.Getenv("IMGCONTENT_GCS_BUCKET")
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

// NewContentReader init ContentReader then return it.
func NewContentReader(bucketName string, client *storage.Client) imgcontent.ContentReader {
	return &ContentRepository{BucketName: bucketName, Client: client}
}

// Write ...
func (r *ContentRepository) Write(ctx context.Context, file io.Reader, path imgcontent.Path) (url url.URL, err error) {
	bucket := r.Client.Bucket(r.BucketName)
	wc := bucket.Object(string(path)).NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		return
	}
	if err = wc.Close(); err != nil {
		return
	}

	// TODO(micheam): duplicated
	url.Scheme = "https"
	url.Host = "storage.googleapis.com"
	url.Path = r.BucketName + "/" + string(path)

	return url, nil
}

var _ imgcontent.ContentReader = (*ContentRepository)(nil)

// List ...
func (r *ContentRepository) List(ctx context.Context, Prefix string) (<-chan imgcontent.ImageContent, error) {
	var (
		bucket = r.Client.Bucket(r.BucketName)
		itr    = bucket.Objects(ctx, &storage.Query{Prefix: Prefix})
		stream = make(chan imgcontent.ImageContent, 100)
	)
	go func() {
		defer close(stream)
		for {
			var (
				err error
				o   *storage.ObjectAttrs
			)
			if o, err = itr.Next(); err != nil {
				if err != iterator.Done {
					// TODO(micheam): エラーの時、どうしよう？
					log.Printf("[ERROR] failed to iterate: %v", err.Error())
				}
				break
			}

			// TODO(micheam): duplicated
			url := new(url.URL)
			url.Scheme = "https"
			url.Host = "storage.googleapis.com"
			url.Path = r.BucketName + "/" + string(o.Name)

			chunks := strings.Split(o.Name, "/")
			if chunks[len(chunks)-1] == "" {
				// MEMO(micheam): may be `directory`
				continue
			}
			name, _ := imgcontent.NewName(chunks[len(chunks)-1]) // MEMO(micheam): ignore err!
			stream <- imgcontent.ImageContent{
				Name: *name,
				Path: imgcontent.Path(o.Name),
				URL:  *url}
		}
	}()
	return stream, nil
}
