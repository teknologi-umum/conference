package main

import (
	"errors"

	"conf/user"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type ServerConfig struct {
	userDomain *user.User
}

func NewServer(config *ServerConfig) *echo.Echo {
	e := echo.New()

	// NOTE: Only need to handle CORS, everything else is being handled by the API gateway
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://conf.teknologiumum.com"},
		AllowMethods:     []string{"POST"},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: false,
		MaxAge:           3600, // 1 day
	}))

	e.POST("users", func(c echo.Context) error {
		payload := user.CreateParticipantRequest{}
		if err := c.Bind(&payload); err != nil {
			return err
		}

		err := config.userDomain.CreateParticipant(c.Request().Context(), payload)
		if err != nil {
			if errors.Is(err, user.ErrValidation) {
				return c.JSON(400, ErrorResponse{Error: err.Error()})
			}

			return c.JSON(500, ErrorResponse{Error: "Internal server error"})
		}

		return c.NoContent(201)
	})

	return e
}
