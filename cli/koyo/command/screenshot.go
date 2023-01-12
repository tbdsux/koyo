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
	},
	&cli.IntFlag{
		Name:        "height",
		Value:       800,
		Usage:       "set the viewport height",
		Destination: &optionsHeight,
	},
	&cli.BoolFlag{
		Name:        "fullPage",
		Value:       false,
		Usage:       "screenshot website as full page",
		Destination: &optionsFullpage,
	},
	&cli.StringFlag{
		Name:        "driver",
		Value:       "playwright",
		Usage:       "the driver for the screenshot api to use (playwright | puppeteer)",
		Destination: &optionsDriver,
	},
	&cli.StringFlag{
		Name:        "imageType",
		Value:       "png",
		Usage:       "screenshot output image type (png | jpeg)",
		Destination: &optionsImageType,
	},
	&cli.StringFlag{
		Name:        "whiteHole",
		Value:       "",
		Usage:       "your WhiteHole integration url",
		Destination: &optionsWhiteHole,
	},
	&cli.StringFlag{
		Name:        "api",
		Value:       "",
		Usage:       "your deta space koyo instance app url (https://your-koyo-app.instance.deta.app)",
		Destination: &optionsApi,
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
		return errors.New("no KOYO API set, please set with --api option")
	}

	spinner := spinner.New()
	website := c.Args().Get(0)

	imageType := optionsImageType
	if optionsImageType != "png" && optionsImageType != "jpeg" {
		imageType = "png"
	}

	queryOptions := APIQuery{
		Width:     optionsWidth,
		Height:    optionsHeight,
		FullPage:  optionsFullpage,
		Driver:    optionsDriver,
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
