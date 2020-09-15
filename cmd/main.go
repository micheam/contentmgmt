package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/micheam/imgcontent"
	"github.com/micheam/imgcontent/console"
	"github.com/micheam/imgcontent/gcs"
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
		listCmd,
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
	Action: func(c *cli.Context) (err error) {
		if c.NArg() > 1 {
			return fmt.Errorf("too many args")
		}

		var f *os.File
		if f, err = os.Open(c.Args().First()); err != nil {
			return err
		}
		defer f.Close()

		var (
			ctx    = context.Background()
			client *storage.Client
		)
		if client, err = gcs.NewClient(ctx); err != nil {
			log.Fatalf("Failed to create new Google Cloud Storage client: %v", err)
		}

		var fname *imgcontent.Name
		if fname, err = imgcontent.NewName(f.Name()); err != nil {
			log.Fatal(err.Error())
		}

		return imgcontent.
			NewUpload(
				imgcontent.TimeBaseContentPathBuilder{},
				gcs.NewContentWriter(gcs.BucketName(ctx), client),
			).
			Exec(
				ctx,
				imgcontent.UploadInput{Name: *fname, Reader: f},
				console.UploadResultHandler(),
			)
	},
}

var listCmd = cli.Command{
	Name:      "list",
	Usage:     "list existing image content",
	ArgsUsage: "[prefix]",
	Action: func(c *cli.Context) (err error) {
		var (
			ctx    = context.Background()
			client *storage.Client
		)
		if client, err = gcs.NewClient(ctx); err != nil {
			log.Fatalf("Failed to create new Google Cloud Storage client: %v", err)
		}

		var (
			reader  = gcs.NewContentReader(gcs.BucketName(ctx), client)
			usecase = imgcontent.NewList(reader)
			input   = imgcontent.ListInput{
				Prefix: c.Args().First(),
			}
		)
		return usecase.Exec(ctx, input, console.ListResultHandler())
	},
}
