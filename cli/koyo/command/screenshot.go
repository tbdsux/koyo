package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/leaanthony/spinner"
	"github.com/urfave/cli/v2"
)

var (
	optionsWidth        int
	optionsHeight       int
	optionsFullpage     bool
	optionsDriver       string
	optionsImageType    string
	optionsWhiteHole    string
	optionsApi          string
	optionsApiKey       string
	optionsOutput       string
	optionsSaveToDrive  bool
	optionsSaveNoOutput bool
)

var ScreenshotFlags = []cli.Flag{
	&cli.IntFlag{
		Name:        "width",
		Value:       1280,
		Usage:       "set the viewport width",
		Destination: &optionsWidth,
		DefaultText: "config.width || 1280",
		Category:    "Parameters",
	},
	&cli.IntFlag{
		Name:        "height",
		Value:       800,
		Usage:       "set the viewport height",
		Destination: &optionsHeight,
		DefaultText: "config.height || 800",
		Category:    "Parameters",
	},
	&cli.BoolFlag{
		Name:        "fullPage",
		Value:       false,
		Usage:       "screenshot website as full page",
		Destination: &optionsFullpage,
		DefaultText: "config.fullPage || false",
		Category:    "Parameters",
	},
	&cli.StringFlag{
		Name:        "driver",
		Value:       "playwright",
		Usage:       "the driver for the screenshot api to use (playwright | puppeteer)",
		Destination: &optionsDriver,
		DefaultText: "config.driver || playwright",
		Category:    "Parameters",
	},
	&cli.StringFlag{
		Name:        "imageType",
		Value:       "png",
		Usage:       "screenshot output image type (png | jpeg)",
		Destination: &optionsImageType,
		DefaultText: "config.imageType || png",
		Category:    "Parameters",
	},
	&cli.StringFlag{
		Name:        "whiteHole",
		Value:       "",
		Usage:       "your WhiteHole integration url",
		Destination: &optionsWhiteHole,
		DefaultText: "config.whiteHole || ",
		Category:    "Parameters",
	},
	&cli.StringFlag{
		Name:        "api",
		Value:       "",
		Usage:       "your deta space koyo instance app url (https://your-koyo-app.instance.deta.app)",
		Destination: &optionsApi,
		DefaultText: "config.api || ",
		Category:    "CLI Configurations",
	},
	&cli.StringFlag{
		Name:        "apiKey",
		Value:       "",
		Usage:       "your space app api key (generate from your dashboard)",
		Destination: &optionsApiKey,
		DefaultText: "config.apiKey || ",
		Category:    "CLI Configurations",
	},
	&cli.StringFlag{
		Name:        "output",
		Value:       "",
		Usage:       "output filename",
		DefaultText: "<website>.<imageType>",
		Destination: &optionsOutput,
		Category:    "CLI Configurations",
	},
	&cli.BoolFlag{
		Name:        "save-to-drive",
		Value:       false,
		Usage:       "save the screenshot to drive",
		DefaultText: "false",
		Destination: &optionsSaveToDrive,
		Category:    "Parameters",
	},
	&cli.BoolFlag{
		Name:        "no-output",
		Value:       false,
		Usage:       "do not return the screenshot image output (can be only used with --save-to-drive flag)",
		DefaultText: "false",
		Destination: &optionsSaveNoOutput,
		Category:    "Parameters",
	},
}

type APIQuery struct {
	Width        int    `url:"width,omitempty"`
	Height       int    `url:"height,omitempty"`
	FullPage     bool   `url:"fullPage,omitempty"`
	Driver       string `url:"driver,omitempty"`
	ImageType    string `url:"imageType,omitempty"`
	WhiteHole    string `url:"whiteHole,omitempty"`
	SaveToDrive  bool   `url:"saveToDrive,omitempty"`
	SaveNoOutput bool   `url:"saveNoOutput,omitempty"`
}

type APIBody struct {
	Website string `json:"website"`
}

type APIError struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var Screenshot = func(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("<website> argument required")
	}

	if optionsApi == "" {
		altApi := myConfig.String("api", "")
		if altApi == "" {
			return errors.New("no KOYO API set, please set with --api option")
		}

		optionsApi = altApi
	}

	if optionsApiKey == "" {
		altApiKey := myConfig.String("apiKey", "")
		if altApiKey == "" {
			return errors.New("missing Space APP Api Key, please generate an `api_key` from your space app settings and set it with the --apiKey flag")
		}

		optionsApiKey = altApiKey
	}

	if optionsWhiteHole == "" {
		altWhitehole := myConfig.String("whiteHole", "")
		if altWhitehole != "" {
			optionsWhiteHole = altWhitehole
		}
	}

	spinner := spinner.New()
	website := c.Args().Get(0)

	imageType := myConfig.String("imageType", optionsImageType)
	if imageType != "png" && imageType != "jpeg" {
		imageType = "png"
	}

	driver := myConfig.String("driver", optionsDriver)
	if driver != "playwright" && driver != "puppeteer" {
		driver = "playwright"
	}

	queryOptions := APIQuery{
		Width:        myConfig.Int("width", optionsWidth),
		Height:       myConfig.Int("height", optionsHeight),
		FullPage:     myConfig.Bool("fullPage", optionsFullpage),
		Driver:       driver,
		ImageType:    imageType,
		WhiteHole:    optionsWhiteHole,
		SaveToDrive:  optionsSaveToDrive,
		SaveNoOutput: optionsSaveNoOutput,
	}
	q, _ := query.Values(queryOptions)

	body := APIBody{
		Website: website,
	}
	jsonBody, _ := json.Marshal(body)

	url := ""
	if strings.HasSuffix(optionsApi, "/") {
		url = optionsApi + "api/screenshot?" + q.Encode()
	} else {
		url = optionsApi + "/api/screenshot?" + q.Encode()
	}

	spinner.Start("Fetching screenshot...")

	// send request to screenshot api
	r, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		spinner.Error(err.Error())
		return errors.New("there was a problem while parsing the request body")
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("X-Space-App-Key", optionsApiKey)

	client := &http.Client{}
	resp, err := client.Do(r)

	if err != nil {
		spinner.Error(err.Error())
		return errors.New("there was a problem while trying to fetch your screenshot")
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var err APIError
		json.NewDecoder(resp.Body).Decode(&err)

		fmt.Println(err)

		spinner.Error(err.Message)
		return errors.New(err.Message)
	}

	if optionsSaveToDrive && optionsSaveNoOutput {
		// if no output is passed, api will return a json ok response
		spinner.Success("Successfully saved screenshot to Drive, you can check your screenshots from your drive with the `drive` command.")
		return nil
	}

	// generate output filename
	filename := fmt.Sprintf("%s.%s", strings.ReplaceAll(strings.ReplaceAll(website, "https://", ""), "http://", ""), optionsImageType)
	if optionsOutput != "" {
		filename = fmt.Sprintf("%s.%s", optionsOutput, optionsImageType)
	}

	file, err := os.Create(filename)
	if err != nil {
		spinner.Error(err.Error())
		return err
	}

	defer file.Close()

	// write response to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		spinner.Error(err.Error())
		return err
	}

	spinner.Success(fmt.Sprintf("Saved screenshot to `%s`", filename))

	return nil
}
