package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/micheam/imgcontent"
	"github.com/micheam/imgcontent/console"
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
		f, err := os.Open(filepath)
		if err != nil {
			return err
		}
		defer f.Close()

		ctx := context.Background()

		// init contentPathBuilder
		contentPathBuilder := imgcontent.TimeBaseContentPathBuilder{}

		// init contentWriter
		// TODO(michema): spec conf from file
		os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", os.Getenv("IMGCONTENT_GCS_CREDENTIALS"))
		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}

		contentWriter := imgcontent.GCPContentRepository{
			BucketName: os.Getenv("IMGCONTENT_GCS_BUCKET"),
			Client:     client,
		}

		// create InputData
		fname, err := imgcontent.NewFilename(f.Name())
		if err != nil {
			log.Fatal(err.Error())
		}

		return imgcontent.
			NewUpload(contentPathBuilder, contentWriter).
			Exec(ctx, imgcontent.UploadInput{Filename: *fname, Reader: f}, console.UploadResultHandler())
	},
}
