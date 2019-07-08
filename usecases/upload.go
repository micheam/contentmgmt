package usecases

import (
	"context"
	"github.com/micheam/imgcontent/entities"
	"github.com/pkg/errors"
	"io"
	"net/url"
)

type Upload struct { // {{{2
	ContentPathBuilder
	ContentWriter
	UploadPresenter
}

type UploadInput struct {
	entities.Filename
	Reader io.Reader
}

type UploadOutput struct {
	entities.Filename
	url.URL
}

func (u Upload) Handle(ctx context.Context, input UploadInput) error { // {{{2

	path, err := u.ContentPathBuilder.Build(ctx, BuildContentPathInput{
		Filename: input.Filename,
	})
	if err != nil {
		return errors.Cause(err)
	}

	writeResult, err := u.ContentWriter.Write(ctx, ContentWriteInput{
		ContentPath: path,
		Reader:      input.Reader,
	})

	result := UploadOutput{
		Filename: input.Filename,
		URL:      writeResult.URL,
	}
	return u.UploadPresenter.Complete(ctx, result)
}

type ContentPathBuilder interface { // {{{2
	Build(ctx context.Context, input BuildContentPathInput) (entities.ContentPath, error)
}

type BuildContentPathInput struct {
	entities.Filename
}

type ContentWriter interface { // {{{2
	Write(ctx context.Context, input ContentWriteInput) (ContentWriteOutput, error)
}

type ContentWriteInput struct {
	Reader io.Reader
	entities.ContentPath
}

type ContentWriteOutput struct {
	url.URL
}

type UploadPresenter interface { // {{{2
	Complete(ctx context.Context, data UploadOutput) error
}
