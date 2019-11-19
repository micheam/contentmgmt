package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"

	"cloud.google.com/go/storage"
	"github.com/atotto/clipboard"
	"github.com/micheam/contentmgmt"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	Version string = "0.2.0"
)

func main() {
	app := cli.NewApp()
	app.Name = "imagecontent"
	app.Usage = "manage img content"
	app.Version = Version
	app.Author = "Michto Maeda"
	app.Email = "https://github.com/micheam"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "verbose",
			Usage:  "Print detail log",
			Hidden: true,
		},
	}
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
	Aliases:   []string{"u"},
	Usage:     "upload file as a web content",
	ArgsUsage: "<filepath>",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format",
			Usage: "Display result with specified format. [mkd,html,adoc]",
		},
		cli.BoolFlag{
			Name:  "clipboard,c",
			Usage: "Write result to clipboard",
		},
	},
	Action: func(c *cli.Context) error {

		if c.GlobalBool("verbose") {
			log.SetOutput(os.Stderr)
		} else {
			log.SetOutput(ioutil.Discard)
		}

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
		var resultTemplate string
		switch c.String("format") {
		case "mkd":
			resultTemplate = MarkdownResultTemplate
		case "html":
			resultTemplate = HTMLResultTemplate
		case "asd":
			resultTemplate = AsciidocResultTemplate
		case "":
			resultTemplate = DefaultResultTemplate
		default:
			return errors.Errorf("Unknown format %q. You must specify one of 'mkd','html' or 'asd'.\n"+
				"See: 'imgcontent upload help' for detail", c.String("format"))
		}
		uploadPresenter := ConsoleUploadResultPresenter{
			Tmpl:             resultTemplate,
			WriteToClipboard: c.Bool("clipboard"),
		}

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
	Tmpl             string
	WriteToClipboard bool
}

const (
	DefaultResultTemplate  = `{{.Scheme}}://{{.Host}}/{{.Path}}`
	MarkdownResultTemplate = `![{{.Filename}}]({{.Scheme}}://{{.Host}}/{{.Path}})`
	HTMLResultTemplate     = `<img alt="{{.Filename}}" src="{{.Scheme}}://{{.Host}}/{{.Path}}" />`
	AsciidocResultTemplate = `image::{{.Scheme}}://{{.Host}}/{{.Path}}[{{.Filename}}]`
)

func (c ConsoleUploadResultPresenter) Complete(ctx context.Context, data contentmgmt.UploadOutput) error {

	var err error

	model := struct {
		Scheme, Host, Path, Filename string
	}{
		Scheme:   data.URL.Scheme,
		Host:     data.URL.Host,
		Path:     data.URL.EscapedPath(),
		Filename: data.Filename.Value,
	}

	tmpl := DefaultResultTemplate
	if c.Tmpl != "" {
		tmpl = c.Tmpl
	}

	log.Printf("result: template %q", tmpl)

	t := template.Must(template.New("result-template").Parse(tmpl))

	var buf bytes.Buffer
	if err = t.Execute(&buf, model); err != nil {
		return errors.Wrap(err, "fialed to execute result tempalte")
	}

	result := buf.String()
	if c.WriteToClipboard {
		log.Printf("write result to clipboard")
		if err = clipboard.WriteAll(result); err != nil {
			return errors.Wrap(err, "failed to write result to clipboard")
		}
	}

	fmt.Fprintf(os.Stdout, result)
	return nil
}
