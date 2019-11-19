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
	"github.com/micheam/contentmgmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	Version string = "0.1.0"
)

func main() {
	app := cli.NewApp()
	app.Name = "imagecontent"
	app.Usage = "manage img content"
	app.Version = Version
	app.Author = "Michto Maeda"
	app.Email = "https://github.com/micheam"
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
	Usage:     "upload file as a web content",
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
		contentWriter := contentmgmt.GCPContentRepository{
			BucketName: bucketName,
			Client:     client,
		}

		// init presenter
		// TODO(michema): spec template from file
		resultTemplate := `![{{.Alt}}]({{.Scheme}}://{{.Host}}/{{.Path}})`
		uploadPresenter := ConsoleUploadResultPresenter{Tmpl: &resultTemplate}

		// init interacter
		usecase := contentmgmt.UploadUsecase{
			PathBuilder: contentPathBuilder,
			Writer:      contentWriter,
			Presenter:   uploadPresenter,
		}

		// create InputData
		filename, err := contentmgmt.NewFilename(file.Name())
		if err != nil {
			log.Fatal(err.Error())
		}

		request := contentmgmt.UploadInput{
			Filename: *filename,
			Reader:   file,
		}

		return usecase.Handle(ctx, request)
	},
}

type ConsoleContentPathBuilder struct {
	BaseTime *time.Time
	Tmpl     *string
}

const DefaultContentPathTemplate = `{{.BaseTime.Format "2006/01/02/030405"}}.{{.FileName}}`

func (c ConsoleContentPathBuilder) Build(
	ctx context.Context, fname contentmgmt.Filename) (path contentmgmt.ContentPath, err error) {

	basetime := time.Now()
	if c.BaseTime != nil {
		basetime = time.Now()
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

	path = contentmgmt.ContentPath(buf.String())
	log.Printf("content path: computed %q", path)

	return
}

type ConsoleUploadResultPresenter struct {
	Tmpl *string
}

const DefaultResultTemplate = `{{.Scheme}}://{{.Host}}/{{.Path}}`

func (c ConsoleUploadResultPresenter) Complete(ctx context.Context, data contentmgmt.UploadOutput) error {

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
