package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

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
		panic(err)
	}

	conn, err := pgxpool.New(
		context.Background(),
		fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", config.DBUser, config.DBPassword, config.DBHost, config.Port),
	)
	if err != nil {
		panic(err)
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
			log.Println(err.Error())
		}
	}()

	if err := e.Start(net.JoinHostPort("", config.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln(err)
	}
}
