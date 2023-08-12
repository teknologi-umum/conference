package main

import (
	"conf/config"
	"conf/user"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	config, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:     config.DBHost,
		Port:     config.DBPort,
		User:     config.DBUser,
		Password: config.DBPassword,
		Database: config.DBName,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		panic(err)
	}
	defer conn.Close()

	e := echo.New()

	userDomain := user.New(conn)
	e.POST("users", func(c echo.Context) error {
		payload := user.CreateParticipantRequest{}
		if err := c.Bind(&payload); err != nil {
			return err
		}

		err := userDomain.CreateParticipant(payload)
		if err != nil {
			if errors.Is(err, user.ErrInvalidUserName) || errors.Is(err, user.ErrInvalidEmail) {
				c.JSON(400, ErrorResponse{Error: err.Error()})
			}

			return c.JSON(500, ErrorResponse{Error: "Internal server error"})
		}

		return c.NoContent(201)
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Port)))
}
