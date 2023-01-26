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
	optionsWidth     int
	optionsHeight    int
	optionsFullpage  bool
	optionsDriver    string
	optionsImageType string
	optionsWhiteHole string
	optionsApi       string
	optionsOutput    string
)

var ScreenshotFlags = []cli.Flag{
	&cli.IntFlag{
		Name:        "width",
		Value:       1280,
		Usage:       "set the viewport width",
		Destination: &optionsWidth,
		DefaultText: "config.width || 1280",
	},
	&cli.IntFlag{
		Name:        "height",
		Value:       800,
		Usage:       "set the viewport height",
		Destination: &optionsHeight,
		DefaultText: "config.height || 800",
	},
	&cli.BoolFlag{
		Name:        "fullPage",
		Value:       false,
		Usage:       "screenshot website as full page",
		Destination: &optionsFullpage,
		DefaultText: "config.fullPage || false",
	},
	&cli.StringFlag{
		Name:        "driver",
		Value:       "playwright",
		Usage:       "the driver for the screenshot api to use (playwright | puppeteer)",
		Destination: &optionsDriver,
		DefaultText: "config.driver || playwright",
	},
	&cli.StringFlag{
		Name:        "imageType",
		Value:       "png",
		Usage:       "screenshot output image type (png | jpeg)",
		Destination: &optionsImageType,
		DefaultText: "config.imageType || png",
	},
	&cli.StringFlag{
		Name:        "whiteHole",
		Value:       "",
		Usage:       "your WhiteHole integration url",
		Destination: &optionsWhiteHole,
		DefaultText: "config.whiteHole || ",
	},
	&cli.StringFlag{
		Name:        "api",
		Value:       "",
		Usage:       "your deta space koyo instance app url (https://your-koyo-app.instance.deta.app)",
		Destination: &optionsApi,
		DefaultText: "config.api || ",
	},
	&cli.StringFlag{
		Name:        "output",
		Value:       "",
		Usage:       "output filename",
		DefaultText: "<website>.<imageType>",
		Destination: &optionsOutput,
	},
}

type APIQuery struct {
	Width     int    `url:"width,omitempty"`
	Height    int    `url:"height,omitempty"`
	FullPage  bool   `url:"fullPage,omitempty"`
	Driver    string `url:"driver,omitempty"`
	ImageType string `url:"imageType,omitempty"`
	WhiteHole string `url:"whiteHole,omitempty"`
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
		Width:     myConfig.Int("width", optionsWidth),
		Height:    myConfig.Int("height", optionsHeight),
		FullPage:  myConfig.Bool("fullPage", optionsFullpage),
		Driver:    driver,
		ImageType: imageType,
		WhiteHole: optionsWhiteHole,
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
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		spinner.Error(err.Error())
		return errors.New("there was a problem while trying to fetch your screenshot")
	}

	if resp.StatusCode != 200 {
		var err APIError
		json.NewDecoder(resp.Body).Decode(&err)

		spinner.Error(err.Message)
		return errors.New(err.Message)
	}

	defer resp.Body.Close()

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
