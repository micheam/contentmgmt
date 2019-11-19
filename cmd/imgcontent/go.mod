module github.com/micheam/contentmgmt/cmd/imgcontent

go 1.12

require (
	cloud.google.com/go/storage v1.3.0
	github.com/atotto/clipboard v0.1.2
	github.com/micheam/contentmgmt v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.8.1
	github.com/urfave/cli v1.22.1
)

replace github.com/micheam/contentmgmt => ../../
