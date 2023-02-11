package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/TheBoringDude/scuffed-go/requester"
	"github.com/leaanthony/spinner"
	"github.com/urfave/cli/v2"
)

var (
	optionsDriveApi    string
	optionsDriveApiKey string
)

type GetDriveFilesResponse struct {
	Error bool     `json:"error"`
	Data  []string `json:"data"`
	Code  int      `json:"code"`
}

var DriveCommand = &cli.Command{
	Name:        "drive",
	Description: "Manage the screenshot images saved on your drive.",
	Subcommands: []*cli.Command{DriveGetCommand},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "api",
			Value:       "",
			Usage:       "your deta space koyo instance app url (https://your-koyo-app.instance.deta.app)",
			Destination: &optionsDriveApi,
			DefaultText: "config.api || ",
		},
		&cli.StringFlag{
			Name:        "apiKey",
			Value:       "",
			Usage:       "your space app api key (generate from your dashboard)",
			Destination: &optionsDriveApiKey,
			DefaultText: "config.apiKey || ",
		},
	},
	Action: func(ctx *cli.Context) error {
		if optionsDriveApi == "" {
			altApi := myConfig.String("api", "")
			if altApi == "" {
				return errors.New("no KOYO API set, please set with --api option")
			}

			optionsDriveApi = altApi
		}

		if optionsDriveApiKey == "" {
			altApiKey := myConfig.String("apiKey", "")
			if altApiKey == "" {
				return errors.New("missing Space APP Api Key, please generate an `api_key` from your space app settings and set it with the --apiKey flag")
			}

			optionsDriveApiKey = altApiKey
		}

		url := ""
		if strings.HasSuffix(optionsDriveApi, "/") {
			url = optionsDriveApi + "api/drive/files"
		} else {
			url = optionsDriveApi + "/api/drive/files"
		}

		resp := GetDriveFilesResponse{}
		r := requester.NewRequester(&http.Client{})

		request, err := http.NewRequest("GET", url, nil)
		request.Header.Add("X-Space-App-Key", optionsApiKey)
		if err != nil {
			return err
		}

		if err = r.Request(request, &resp); err != nil {
			return err
		}

		fmt.Println("\n Website Screenshots stored in Drive")
		fmt.Println(" ------------------------------------------")
		for _, v := range resp.Data {
			fmt.Println("    - " + v)
		}

		return nil
	},
}

var (
	optionsDriveGetApi    string
	optionsDriveGetApiKey string
)

var DriveGetCommand = &cli.Command{
	Name:        "get",
	Description: "Download a screenshot from your drive.",
	UsageText:   "koyo drive get [filename1] [filename2] [filename3] ...",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "api",
			Value:       "",
			Usage:       "your deta space koyo instance app url (https://your-koyo-app.instance.deta.app)",
			Destination: &optionsDriveGetApi,
			DefaultText: "config.api || ",
		},
		&cli.StringFlag{
			Name:        "apiKey",
			Value:       "",
			Usage:       "your space app api key (generate from your dashboard)",
			Destination: &optionsDriveGetApiKey,
			DefaultText: "config.apiKey || ",
		},
	},
	Action: func(ctx *cli.Context) error {
		if ctx.NArg() == 0 {
			return errors.New("no files to download")
		}

		if optionsDriveGetApi == "" {
			altApi := myConfig.String("api", "")
			if altApi == "" {
				return errors.New("no KOYO API set, please set with --api option")
			}

			optionsDriveGetApi = altApi
		}

		if optionsDriveGetApiKey == "" {
			altApiKey := myConfig.String("apiKey", "")
			if altApiKey == "" {
				return errors.New("missing Space APP Api Key, please generate an `api_key` from your space app settings and set it with the --apiKey flag")
			}

			optionsDriveGetApiKey = altApiKey
		}

		url := ""
		if strings.HasSuffix(optionsDriveGetApi, "/") {
			url = optionsDriveGetApi + "api/drive/files/"
		} else {
			url = optionsDriveGetApi + "/api/drive/files/"
		}

		spinner := spinner.New()

		for _, v := range ctx.Args().Slice() {
			spinner.Start(fmt.Sprintf("Fetching %s", v))

			request, err := http.NewRequest("GET", url+v, nil)
			request.Header.Add("X-Space-App-Key", optionsApiKey)
			if err != nil {
				spinner.Errorf("Failed to initialize request to download %s", v)
				return err
			}

			client := &http.Client{}
			resp, err := client.Do(request)
			if err != nil {
				spinner.Errorf("Failed to download %s", v)
				return err
			}

			if resp.StatusCode != 200 {
				var err APIError
				json.NewDecoder(resp.Body).Decode(&err)

				spinner.Error(err.Message)
				return errors.New(err.Message)
			}

			defer resp.Body.Close()

			file, err := os.Create(v)
			if err != nil {
				spinner.Error(err.Error())
				return err
			}

			// write response to file
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				spinner.Error(err.Error())
				return err
			}

			spinner.Success(fmt.Sprintf("Saved screenshot to `%s`", v))

		}

		return nil
	},
}
