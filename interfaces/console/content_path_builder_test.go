package console

import (
	"context"
	"reflect"
	"testing"
	"time"

	en "github.com/micheam/imgcontent/entities"
)

func TestDefaultContentPathBuilder_Build(t *testing.T) {

	basedate, _ := time.Parse("2006-01-02", "2014-10-27")

	type args struct {
		ctx      context.Context
		filename en.Filename
	}
	tests := []struct {
		name     string
		c        DefaultContentPathBuilder
		args     args
		wantPath en.ContentPath
		wantErr  bool
	}{
		{
			name:     "return path with time prefix and filename",
			c:        DefaultContentPathBuilder{BaseDate: basedate},
			args:     args{ctx: context.TODO(), filename: en.Filename{Value: "filename.png", Valid: true}},
			wantPath: en.ContentPath("2014/10/27/filename.png"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, err := tt.c.Build(tt.args.ctx, tt.args.filename)
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
