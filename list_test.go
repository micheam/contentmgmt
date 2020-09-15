package imgcontent

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockContentReader struct {
	mock.Mock
}

func (r *mockContentReader) List(ctx context.Context, prefix string) (<-chan ImageContent, error) {
	args := r.Called(ctx, prefix)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(chan ImageContent), args.Error(1)
}

func Test_list_Exec(t *testing.T) {

	// helper func: Create handler function
	var handler = func(t *testing.T, expected []ImageContent) ListResultHandler {
		return func(cctx context.Context, data ListOutput) error {
			_contents := make([]ImageContent, 0, len(expected))
			for c := range data.Contents {
				_contents = append(_contents, c)
			}
			assert.EqualValues(t, expected, _contents)
			return nil
		}
	}

	// helper func: Create channel, which streaming contents.
	var stream = func(contents []ImageContent) chan ImageContent {
		stream := make(chan ImageContent, len(contents))
		go func() {
			for _, c := range contents {
				stream <- c
			}
			close(stream)
		}()
		return stream
	}

	t.Run("Prefix may be empty", func(t *testing.T) {
		var (
			ctx    = context.TODO()
			reader = new(mockContentReader)
			input  = ListInput{Prefix: ""}
			cs     = []ImageContent{{Name: "foo.jpg"}, {Name: "bar.jpg"}}
			h      = handler(t, cs)
		)
		reader.Mock.On("List", ctx, input.Prefix).Return(stream(cs), nil)
		assert.NoError(t, NewList(reader).Exec(ctx, input, h))
	})
	t.Run("Never Error on NotFound", func(t *testing.T) {
		var (
			ctx    = context.TODO()
			reader = new(mockContentReader)
			input  = ListInput{Prefix: "foooooooooo"}
			cs     = []ImageContent{}
			h      = handler(t, cs)
		)
		reader.Mock.On("List", ctx, input.Prefix).Return(stream(cs), nil)
		assert.NoError(t, NewList(reader).Exec(ctx, input, h))
	})
	t.Run("must return error which happen inside", func(t *testing.T) {
		var (
			reader = new(mockContentReader)
			ctx    = context.TODO()
			input  = ListInput{}
			h      = func(cctx context.Context, data ListOutput) error {
				t.Fail()
				return nil
			}
		)
		aErr := errors.New("this is an Error")
		reader.Mock.On("List", ctx, "").Return(nil, aErr)
		gotErr := NewList(reader).Exec(ctx, input, h)
		assert.Error(t, gotErr)
		assert.ErrorIs(t, aErr, gotErr)
	})
}
