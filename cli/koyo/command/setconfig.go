package command

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	simplefiletest "github.com/TheBoringDude/simple-filetest"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yamlv3"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli/v2"
)

func HomeDir() string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalln(err)
	}

	return home
}

var (
	home     string = HomeDir()
	cfgFile  string = path.Join(home, ".koyo.yaml")
	myConfig        = config.NewWith("koyo", func(c *config.Config) {
		c.AddDriver(yamlv3.Driver)
	})
)

// Initialize the default config at $HOME/.koyo.yaml
func InitConfig() {
	if !simplefiletest.FileExists(cfgFile) {
		// create initial first config
		myConfig.Set("width", 1280)
		myConfig.Set("height", 800)
		myConfig.Set("fullPage", false)
		myConfig.Set("driver", "playwright")
		myConfig.Set("imageType", "png")
		myConfig.Set("whiteHole", "")
		myConfig.Set("api", "")
		myConfig.Set("apiKey", "")

		writeConfig()
	}

	// load config
	myConfig.LoadFiles(cfgFile)
}

func writeConfig() {
	buf := new(bytes.Buffer)
	myConfig.DumpTo(buf, config.Yaml)
	os.WriteFile(cfgFile, buf.Bytes(), 0755)
}

var (
	optShow bool
)

var SetConfigCommand = &cli.Command{
	Name:        "set",
	Description: "Set api defaults to config.",
	UsageText:   "koyo set config.[key] [value]",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "show",
			Usage:       "Show the available config keys.",
			Destination: &optShow,
			Value:       false,
		},
	},
	Action: func(ctx *cli.Context) error {
		// show the config if --show is true
		if optShow {
			fmt.Println("--- Current default config")
			fmt.Printf("width: %s\n", myConfig.String("width"))
			fmt.Printf("height: %s\n", myConfig.String("height"))
			fmt.Printf("fullPage: %s\n", myConfig.String("fullPage"))
			fmt.Printf("driver: %s\n", myConfig.String("driver"))
			fmt.Printf("imageType: %s\n", myConfig.String("imageType"))
			fmt.Printf("whiteHole: %s\n", myConfig.String("whiteHole"))
			fmt.Printf("api: %s\n", myConfig.String("api"))
			fmt.Printf("apiKey: %s\n", myConfig.String("apiKey"))

			return nil
		}

		// check if args length is correct
		if ctx.Args().Len() != 2 {
			return errors.New("invalid command argument. Use --help for guide")
		}

		key := ctx.Args().Get(0)
		value := ctx.Args().Get(1)

		if key == "config.width" {
			cValue, err := strconv.Atoi(value)
			if err != nil {
				log.Fatalf("Invalid parsing `width` value: %v \n", err)
			}

			myConfig.Set("width", cValue)
		}
		if key == "config.height" {
			cValue, err := strconv.Atoi(value)
			if err != nil {
				log.Fatalf("Invalid parsing `height` value: %v \n", err)
			}

			myConfig.Set("height", cValue)
		}
		if key == "config.fullPage" {
			cValue, err := strconv.ParseBool(value)
			if err != nil {
				log.Fatalf("Invalid parsing `fullPage` value: %v \n", err)
			}

			myConfig.Set("fullPage", cValue)
		}
		if key == "config.driver" {
			if value != "playwright" && value != "puppeteer" {
				log.Fatalln("Invalid `driver` value. Can only be (playwright | puppeteer)")
			}

			myConfig.Set("driver", value)
		}
		if key == "config.imageType" {
			if value != "png" && value != "jpeg" {
				log.Fatalln("Invalid `imageType` value. Can only be (png | jpeg)")
			}

			myConfig.Set("imageType", value)
		}
		if key == "config.whiteHole" {
			myConfig.Set("whiteHole", value)
		}
		if key == "config.api" {
			myConfig.Set("api", value)
		}
		if key == "config.apiKey" {
			myConfig.Set("apiKey", value)
		}

		// update new config change
		writeConfig()

		return nil
	},
}
