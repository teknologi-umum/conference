package main

import (
	"conf/core"
	"conf/user"
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ErrorResponse struct {
	Key   string `json:"key"`
	Error string `json:"error"`
}

type ServerConfig struct {
	userDomain *user.UserDomain
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
		type payload struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		p := payload{}
		if err := c.Bind(&p); err != nil {
			return err
		}

		errs := config.userDomain.CreateParticipant(c.Request().Context(), user.CreateParticipant{
			Name:  p.Name,
			Email: p.Email,
		})

		errorValidations := []ErrorResponse{}
		if errs != nil {
			for _, err := range errs {
				if errors.Is(err.Err, core.ErrValidation) {
					errorValidations = append(errorValidations, ErrorResponse{
						Key:   err.Key,
						Error: err.Err.Error(),
					})
				}
			}

			if len(errorValidations) > 0 {
				return c.JSON(400, errorValidations)
			}

			return c.JSON(500, ErrorResponse{Error: "Internal server error"})
		}

		return c.NoContent(201)
	})

	return e
}
