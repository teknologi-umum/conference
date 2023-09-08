package main

import (
	"errors"
	"net/http"

	sentryecho "github.com/getsentry/sentry-go/echo"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerConfig struct {
	userDomain *UserDomain
}

func NewServer(config *ServerConfig) *echo.Echo {
	e := echo.New()

	sentryMiddleware := sentryecho.New(sentryecho.Options{Repanic: false})
	e.Use(sentryMiddleware)

	// NOTE: Only need to handle CORS, everything else is being handled by the API gateway
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://conf.teknologiumum.com"},
		AllowMethods:     []string{http.MethodPost},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: false,
		MaxAge:           3600, // 1 day
	}))

	e.Use(middleware.RequestID())

	e.GET("/ping", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	e.POST("/users", func(c echo.Context) error {
		requestId := c.Response().Header().Get(echo.HeaderXRequestID)
		sentryHub := sentryecho.GetHubFromContext(c)
		sentryHub.Scope().SetTag("request-id", requestId)

		type payload struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}
		p := payload{}
		if err := c.Bind(&p); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message":    "Invalid request body",
				"errors":     err.Error(),
				"request_id": requestId,
			})
		}

		err := config.userDomain.CreateParticipant(
			c.Request().Context(),
			CreateParticipantRequest{
				Name:  p.Name,
				Email: p.Email,
			},
		)
		if err != nil {
			var validationError *ValidationError
			if errors.As(err, &validationError) {
				return c.JSON(http.StatusBadRequest, echo.Map{
					"message":    "Validation error",
					"errors":     validationError.Errors,
					"request_id": requestId,
				})
			}

			sentryHub.CaptureException(err)
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message":    "Internal server error",
				"errors":     "Internal server error",
				"request_id": requestId,
			})
		}

		return c.NoContent(http.StatusCreated)
	})

	return e
}
