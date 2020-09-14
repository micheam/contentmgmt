package main

import (
	"context"
	"testing"
	"time"

	"github.com/micheam/imgcontent"
	"github.com/stretchr/testify/assert"
)

func TestDefaultContentPathBuilder_Build(t *testing.T) {
	basetime, _ := time.Parse("2006-01-02 03:04:05", "2014-10-27 11:12:13")
	sut := ConsoleContentPathBuilder{BaseTime: &basetime}
	ctx := context.TODO()
	filename := imgcontent.Filename{Value: "filename.png", Valid: true}
	got, err := sut.Build(ctx, filename)
	assert.NoError(t, err)
	assert.EqualValues(t, imgcontent.ContentPath("2014/10/27/111213.filename.png"), got)
}
