package main

import (
	"log"
	"os"

	"github.com/tbdsux/koyo/cli/koyo/command"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:   "koyo",
		Usage:  "Screenshot any website with ease.",
		Flags:  command.ScreenshotFlags,
		Action: command.Screenshot,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
