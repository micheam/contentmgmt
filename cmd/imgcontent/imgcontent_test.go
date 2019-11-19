package main

import (
	"context"
	"github.com/micheam/contentmgmt"
	"reflect"
	"testing"
	"time"
)

func TestDefaultContentPathBuilder_Build(t *testing.T) {

	basetime, _ := time.Parse("2006-01-02 03:04:05", "2014-10-27 11:12:13")

	type args struct {
		ctx      context.Context
		filename contentmgmt.Filename
	}
	tests := []struct {
		name     string
		c        ConsoleContentPathBuilder
		args     args
		wantPath contentmgmt.ContentPath
		wantErr  bool
	}{
		{
			name:     "return path with time prefix and filename",
			c:        ConsoleContentPathBuilder{BaseTime: &basetime},
			args:     args{ctx: context.TODO(), filename: contentmgmt.Filename{Value: "filename.png", Valid: true}},
			wantPath: contentmgmt.ContentPath("2014/10/27/111213.filename.png"),
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
