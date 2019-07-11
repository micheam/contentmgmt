package console

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/micheam/imgcontent/entities"
	"github.com/micheam/imgcontent/interfaces/cloudstorage"
	"github.com/micheam/imgcontent/usecases"

	"github.com/urfave/cli"
)

var uploadCmd = cli.Command{
	Name:      "upload",
	Usage:     "画像ファイルをアップロードします",
	ArgsUsage: "<filepath>",
	Action:    doUpload,
}

func doUpload(c *cli.Context) error {

	if c.NArg() > 1 {
		return fmt.Errorf("too many args")
	}

	ctx := context.Background()
	now := time.Now()

	contentPathBuilder := DefaultContentPathBuilder{now}
	contentWriter := cloudstorage.GCPContentRepository{}
	uploadPresenter := ConsoleUploadPresenter{}

	usecase := usecases.Upload{
		ContentPathBuilder: contentPathBuilder,
		ContentWriter:      contentWriter,
		UploadPresenter:    uploadPresenter,
	}

	filepath := c.Args().First()
	filename, err := entities.NewFilename(filepath)
	if err != nil {
		return err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	request := usecases.UploadInput{
		Filename: *filename,
		Reader:   file,
	}

	return usecase.Handle(ctx, request)
}

var _ usecases.UploadPresenter = (*ConsoleUploadPresenter)(nil)
var _ usecases.UploadPresenter = ConsoleUploadPresenter{}

type ConsoleUploadPresenter struct{}

func (c ConsoleUploadPresenter) Complete(ctx context.Context, data usecases.UploadOutput) error {
	log.Printf("data: %v", data)
	panic("not implemented yet")
}
