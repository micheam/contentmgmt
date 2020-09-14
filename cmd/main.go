package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"cloud.google.com/go/storage"
	"github.com/micheam/imgcontent"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	// Version ...
	Version string = "0.1.0"
)

func main() {
	app := cli.NewApp()
	app.Name = "imagecontent"
	app.Usage = "manage img content"
	app.Version = Version
	app.Author = "Michto Maeda"
	app.Email = "michito.maeda@gmail.com"
	app.Commands = []cli.Command{
		uploadCmd,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var uploadCmd = cli.Command{
	Name:      "upload",
	Usage:     "画像ファイルをアップロードします",
	ArgsUsage: "<filepath>",
	Action: func(c *cli.Context) error {

		if c.NArg() > 1 {
			return fmt.Errorf("too many args")
		}

		filepath := c.Args().First()
		file, err := os.Open(filepath)
		if err != nil {
			return err
		}
		defer file.Close()

		ctx := context.Background()

		// init contentPathBuilder
		contentPathBuilder := ConsoleContentPathBuilder{}

		// init contentWriter
		// TODO(michema): spec conf from file
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", os.Getenv("IMGCONTENT_GCS_CREDENTIALS"))
		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}

		bucketName := os.Getenv("IMGCONTENT_GCS_BUCKET")
		contentWriter := imgcontent.GCPContentRepository{
			BucketName: bucketName,
			Client:     client,
		}

		// init presenter
		// TODO(michema): spec template from file
		resultTemplate := `![{{.Alt}}]({{.Scheme}}://{{.Host}}/{{.Path}})`
		uploadPresenter := ConsoleUploadResultPresenter{Tmpl: &resultTemplate}

		// create InputData
		filename, err := imgcontent.NewFilename(file.Name())
		if err != nil {
			log.Fatal(err.Error())
		}

		request := imgcontent.UploadInput{
			Filename: *filename,
			Reader:   file,
		}

		return imgcontent.
			NewUpload(contentPathBuilder, contentWriter, uploadPresenter).
			Exec(ctx, request)
	},
}

// ConsoleContentPathBuilder ...
type ConsoleContentPathBuilder struct {
	BaseTime *time.Time
	Tmpl     *string
}

// DefaultContentPathTemplate ...
const DefaultContentPathTemplate = `{{.BaseTime.Format "2006/01/02/030405"}}.{{.FileName}}`

// Build ...
func (c ConsoleContentPathBuilder) Build(
	ctx context.Context, fname imgcontent.Filename) (path imgcontent.ContentPath, err error) {

	basetime := time.Now()
	if c.BaseTime != nil {
		basetime = *c.BaseTime
	}

	model := struct {
		BaseTime time.Time
		FileName string
	}{
		BaseTime: basetime,
		FileName: fname.Value,
	}

	tmpl := DefaultContentPathTemplate
	if c.Tmpl != nil {
		tmpl = *c.Tmpl
	}

	log.Printf("content path: template %q", tmpl)

	t := template.Must(template.New("content-path-template").Parse(tmpl))
	var buf bytes.Buffer
	if err = t.Execute(&buf, model); err != nil {
		return path, errors.Wrap(err, "failed to build content path")
	}

	path = imgcontent.ContentPath(buf.String())
	log.Printf("content path: computed %q", path)

	return
}

// ConsoleUploadResultPresenter ...
type ConsoleUploadResultPresenter struct {
	Tmpl *string
}

// DefaultResultTemplate ...
const DefaultResultTemplate = `{{.Scheme}}://{{.Host}}/{{.Path}}`

// Complete ...
func (c ConsoleUploadResultPresenter) Complete(ctx context.Context, data imgcontent.UploadOutput) error {

	model := struct {
		Scheme, Host, Path, Alt string
	}{
		Scheme: data.URL.Scheme,
		Host:   data.URL.Host,
		Path:   data.URL.EscapedPath(),
		Alt:    data.Filename.Value,
	}

	tmpl := DefaultResultTemplate
	if c.Tmpl != nil {
		tmpl = *c.Tmpl
	}

	log.Printf("result: template %q", tmpl)

	t := template.Must(template.New("result-template").Parse(tmpl))
	t.Execute(os.Stdout, model)

	return nil
}
