package imgcontent

import (
	"context"
)

type (
	// List define how to list existing contents.
	List interface {
		Exec(ctx context.Context, input ListInput, cb ListResultHandler) error
	}
	// ListResultHandler define how to handle result data.
	ListResultHandler func(ctx context.Context, data ListOutput) error
	// ListInput holds input data of usecase.
	ListInput struct {
		Prefix string
	}
	// ListOutput holds output data of usecase.
	ListOutput struct {
		Contents <-chan ImageContent
	}
)

type list struct {
	Reader ContentReader
}

// NewList return initialized List usecase.
func NewList(rdr ContentReader) List {
	return &list{Reader: rdr}
}

func (u list) Exec(ctx context.Context, input ListInput, cb ListResultHandler) error {
	stream, err := u.Reader.List(ctx, input.Prefix)
	if err != nil {
		return err
	}
	return cb(ctx, ListOutput{Contents: stream})
}
