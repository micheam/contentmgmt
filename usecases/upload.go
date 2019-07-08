package usecases

import (
	"context"
	"github.com/micheam/imgcontent/entities"
	"github.com/pkg/errors"
)

type Upload struct {
	ContentPathBuilder
	ContentWriter
	UploadPresenter
}

type UploadInput struct {
	Filename string
	Content  []byte
}

type UploadOutput struct {
	Filename string
	Url      string
}

func (u Upload) Handle(ctx context.Context, input UploadInput) error {

	path := u.ContentPathBuilder.Build(input.Filename)
	img, err := entities.NewImageContent(path, input.Content)
	if err != nil {
		return errors.Cause(err)
	}

	writeResult, err := u.ContentWriter.Write(ctx, ContentWriteInput{
		ImageContent: img,
	})

	result := UploadOutput{
		Filename: input.Filename,
		Url:      writeResult.Url,
	}
	return u.UploadPresenter.Complete(ctx, result)
}

type ContentPathBuilder interface {
	Build(fname string) string
}

type ContentWriter interface {
	Write(ctx context.Context, input ContentWriteInput) (ContentWriteOutput, error)
}

type ContentWriteInput struct {
	*entities.ImageContent
}

type ContentWriteOutput struct {
	Url string
}

type UploadPresenter interface {
	Complete(ctx context.Context, data UploadOutput) error
}
