package main

import (
	"os"

	"github.com/konojunya/cli-oauth/menu"
)

func main() {
	app := menu.Getapp()
	app.Run(os.Args)
}
