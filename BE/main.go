package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog/log"

	"conf/user"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// TODO: move this out from the main function
type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {

	config, err := GetConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get config")
	}

	app := &cli.App{
		Name:  "teknum-conf",
		Usage: "say a greeting",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Value:       config.DBName,
				Usage:       "db name",
				Destination: &config.DBName,
			},
			&cli.StringFlag{
				Name:        "db-user",
				Value:       config.DBUser,
				Usage:       "db user",
				Destination: &config.DBUser,
			},
			&cli.StringFlag{
				Name:        "db-password",
				Value:       config.DBPassword,
				Usage:       "db password",
				Destination: &config.DBPassword,
			},
			&cli.StringFlag{
				Name:        "db-host",
				Value:       config.DBHost,
				Usage:       "db host",
				Destination: &config.DBHost,
			},
			&cli.StringFlag{
				Name:        "port",
				Value:       config.Port,
				Usage:       "port",
				Destination: &config.Port,
			},
		},
		Commands: []*cli.Command{
			{
				Name: "server",
				Action: func(cCtx *cli.Context) error {
					conn, err := pgxpool.New(
						context.Background(),
						fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", config.DBUser, config.DBPassword, config.DBHost, config.Port, config.DBName),
					)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to connect to database")
					}
					defer conn.Close()
					e := echo.New()

					userDomain := user.New(conn)
					// TODO: move handler out from the main function
					e.POST("users", func(c echo.Context) error {
						payload := user.CreateParticipantRequest{}
						if err := c.Bind(&payload); err != nil {
							return err
						}

						err := userDomain.CreateParticipant(c.Request().Context(), payload)
						if err != nil {
							if errors.Is(err, user.ErrValidation) {
								return c.JSON(400, ErrorResponse{Error: err.Error()})
							}

							return c.JSON(500, ErrorResponse{Error: "Internal server error"})
						}

						return c.NoContent(201)
					})

					exitSig := make(chan os.Signal, 1)
					signal.Notify(exitSig, os.Interrupt, os.Kill)

					go func() {
						<-exitSig
						ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
						defer cancel()

						if err := e.Shutdown(ctx); err != nil {
							log.Error().Err(err).Msg("Failed to shutdown server")
						}
					}()

					if err := e.Start(net.JoinHostPort("", config.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
						log.Fatal().Err(err).Msg("Failed to start server")
					}
					return nil
				},
			},
			{
				Name: "migrate",
				Action: func(cCtx *cli.Context) error {
					conn, err := pgxpool.New(
						context.Background(),
						fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", config.DBUser, config.DBPassword, config.DBHost, config.Port, config.DBName),
					)
					if err != nil {
						log.Fatal().Err(err).Msg("Failed to connect to database")
					}
					defer conn.Close()

					migrate := MigrationNew(conn)
					err = migrate.Migrate(context.Background())
					if err != nil {
						conn.Close()
						log.Fatal().Err(err).Msg("Failed to migrate database")
						return err
					}
					log.Info().Msg("Migrating database")
					return nil
				},
			},
			{
				Name: "dump-attendees",
				Action: func(cCtx *cli.Context) error {
					//TODO: get attendees from postgres and dump to csv in stdout
					log.Info().Msg("Dumping attendees")
					return nil
				},
			},
			{
				Name:      "blast-email",
				Usage:     "blast-email [template] [file list destination of emails]",
				ArgsUsage: "[template] [file list destination of emails]",
				Action: func(cCtx *cli.Context) error {
					templateArg := cCtx.Args().Get(0)
					emailList := cCtx.Args().Tail()

					if templateArg == "" {
						log.Fatal().Msg("Template is required")
					}

					if len(emailList) == 0 {
						log.Fatal().Msg("Email list is required for blasting email minimum 1 email")
					}

					for _, email := range emailList {
						log.Info().Msgf("Sending email to %s", email)
					}

					//TODO: send email to all attendees. Implement email your logic here
					log.Info().Msg("Blasting email")
					return nil
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("Failed to run app")
	}
}
