package main

import (
	"os"
	"time"

	"github.com/urfave/cli/v2"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	_ "gocloud.dev/blob/fileblob"
	_ "gocloud.dev/blob/s3blob"
)

var version string

func App() *cli.App {
	return &cli.App{
		Name:           "teknum-conf",
		Version:        version,
		Description:    "CLI for working with Teknologi Umum Conference backend",
		DefaultCommand: "server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config-file-path",
				EnvVars: []string{"CONFIGURATION_FILE"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "server",
				Action: ServerHandlerAction,
			},
			{
				Name: "healthcheck",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "port",
						Value: "8080",
					},
					&cli.DurationFlag{
						Name:  "timeout",
						Value: time.Second * 15,
					},
				},
				Action: HealthcheckHandlerAction,
			},
			{
				Name: "blast-email",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "subject",
						Value:    "",
						Usage:    "Email subject",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "plaintext-body",
						Value:    "",
						Usage:    "Path to plaintext body file",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "html-body",
						Value:    "",
						Usage:    "Path to HTML body file",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "recipients",
						Value:    "",
						Usage:    "Path to CSV file containing list of emails",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "single-recipient",
						Value:    "",
						Required: false,
					},
				},
				Usage:     "blast-email [subject] [template-plaintext] [template-html-body] [csv-file list destination of emails]",
				ArgsUsage: "[subject] [template-plaintext] [template-html-body] [path-csv-file]",
				Action:    BlastMailHandlerAction,
			},
		},
		Copyright: `   Copyright 2023 Teknologi Umum

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.`,
	}
}

func main() {
	if err := App().Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run app")
	}
}
