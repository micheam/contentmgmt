package console

import (
	"github.com/urfave/cli"
)

var Version string = "0.1.0"

func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "imgcontent"
	app.Usage = "image content management"
	app.Version = Version
	app.Author = "Michto Maeda"
	app.Email = "michito.maeda@gmail.com"
	app.Commands = Commands
	return app
}

var Commands = []cli.Command{
	uploadCmd,
}

