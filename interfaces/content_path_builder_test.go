package interfaces

import (
	"context"
	"reflect"
	"testing"

	e "github.com/micheam/imgcontent/entities"
	uc "github.com/micheam/imgcontent/usecases"
)

func TestDefaultContentPathBuilder_Build(t *testing.T) {
	type args struct {
		ctx   context.Context
		input uc.BuildContentPathInput
	}
	tests := []struct {
		name     string
		c        DefaultContentPathBuilder
		args     args
		wantPath e.ContentPath
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := DefaultContentPathBuilder{}
			gotPath, err := c.Build(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultContentPathBuilder.Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotPath, tt.wantPath) {
				t.Errorf("DefaultContentPathBuilder.Build() = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}
