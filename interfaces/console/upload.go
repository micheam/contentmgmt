package console

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var uploadCmd = cli.Command{
	Name:      "upload",
	Usage:     "画像ファイルをアップロードします",
	ArgsUsage: "<filepath>",
	Action:    doUpload,
}

func doUpload(c *cli.Context) error {

	if c.NArg() > 1 {
		return fmt.Errorf("too many args")
	}

	filepath := c.Args().First()
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	ctx := context.Background()
	return handleUpload(ctx, file)
}
