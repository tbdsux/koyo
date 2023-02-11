package main

import (
	"log"
	"os"

	"github.com/tbdsux/koyo/cli/koyo/command"
	"github.com/urfave/cli/v2"
)

func init() {
	command.InitConfig()
}

func main() {

	app := &cli.App{
		Name:      "koyo",
		Usage:     "Screenshot any website with ease.",
		UsageText: "koyo [...options] [website]",
		Flags:     command.ScreenshotFlags,
		Action:    command.Screenshot,
		Commands:  []*cli.Command{command.SetConfigCommand, command.DriveCommand},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
