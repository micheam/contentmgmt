package console

import (
	"cloud.google.com/go/storage"
	"context"
	"log"
	"os"
	"time"

	"github.com/micheam/imgcontent/entities"
	"github.com/micheam/imgcontent/interfaces/cloudstorage"
	"github.com/micheam/imgcontent/usecases"
)

func handleUpload(ctx context.Context, file *os.File) error {

	// init contentPathBuilder
	now := time.Now()
	contentPathBuilder := DefaultContentPathBuilder{BaseDate: now}

	// init contentWriter
    os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", os.Getenv("IMGCONTENT_GCS_CREDENTIALS"))
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	bucketName := os.Getenv("IMGCONTENT_GCS_BUCKET")
	contentWriter := cloudstorage.GCPContentRepository{
		BucketName: bucketName,
		Client:     client,
	}

	// init presenter
	uploadPresenter := ConsoleUploadPresenter{}

	// init interacter
	usecase := usecases.Upload{
		ContentPathBuilder: contentPathBuilder,
		ContentWriter:      contentWriter,
		UploadPresenter:    uploadPresenter,
	}

	// create InputData
	filename, err := entities.NewFilename(file.Name())
	if err != nil {
		log.Fatal(err.Error())
	}

	request := usecases.UploadInput{
		Filename: *filename,
		Reader:   file,
	}

	return usecase.Handle(ctx, request)
}
