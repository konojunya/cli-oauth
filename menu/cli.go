package menu

import (
	"github.com/konojunya/cli-oauth/action"
	"github.com/urfave/cli"
)

func Getapp() *cli.App {
	app := cli.NewApp()
	config(app)
	app.Commands = getCommands()
	return app
}

func config(app *cli.App) {
	app.Name = "oauth-sample"
	app.Usage = "OAuthを試す"
	app.Version = "1.0.0"
}

func getCommands() []cli.Command {
	return []cli.Command{
		{
			Name:   "tweet",
			Usage:  "Tweet with cli",
			Action: action.Tweet,
		},
	}
}
