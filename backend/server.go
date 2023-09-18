package main

import (
	"errors"
	"mime"
	"net/http"
	"path"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerConfig struct {
	UserDomain   *UserDomain
	TicketDomain *TicketDomain
}

type ServerDependency struct {
	userDomain   *UserDomain
	ticketDomain *TicketDomain
}

func NewServer(config *ServerConfig) *echo.Echo {
	if config.UserDomain == nil || config.TicketDomain == nil {
		// For production backend application, please don't do what I just did.
		// Do a proper nil check and validation for each of your config and dependencies.
		// NEVER call panic(), just return error.
		// I'm in a hackathon (basically in a rush), so I'm doing this.
		// Let me remind you again: don't do what I just did.
		panic("one of the domain dependency is nil")
	}

	dependencies := &ServerDependency{
		userDomain:   config.UserDomain,
		ticketDomain: config.TicketDomain,
	}

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

	e.POST("/users", dependencies.RegisterUser)
	e.POST("/bukti-transfer", dependencies.UploadBuktiTransfer)
	return e
}

type RegisterUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *ServerDependency) RegisterUser(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	sentryHub := sentryecho.GetHubFromContext(c)
	sentryHub.Scope().SetTag("request-id", requestId)

	p := RegisterUserRequest{}
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message":    "Invalid request body",
			"errors":     err.Error(),
			"request_id": requestId,
		})
	}

	err := s.userDomain.CreateParticipant(
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
}

func (s *ServerDependency) UploadBuktiTransfer(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	sentryHub := sentryecho.GetHubFromContext(c)
	sentryHub.Scope().SetTag("request-id", requestId)

	email := c.FormValue("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message":    "Validation error",
			"errors":     "Email field is required",
			"request_id": requestId,
		})
	}
	photoFormFile, err := c.FormFile("photo")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message":    "Validation error",
				"errors":     "Photo field is required",
				"request_id": requestId,
			})
		}

		sentryHub.CaptureException(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message":    "Reading form file",
			"errors":     err.Error(),
			"request_id": requestId,
		})
	}

	photoExtension := path.Ext(photoFormFile.Filename)
	// Guard the content type, the only content type allowed is images.
	var extensionFound = false
	for _, ext := range []string{".gif", ".jpeg", ".jpg", ".png", ".webp"} {
		if ext == photoExtension {
			extensionFound = true
			break
		}
	}
	if !extensionFound {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message":    "Unknown photo file type",
			"errors":     "Unknown photo file type",
			"request_id": requestId,
		})
	}

	photoContentType := mime.TypeByExtension(photoExtension)

	photoFile, err := photoFormFile.Open()
	if err != nil {
		sentryHub.CaptureException(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message":    "Opening photo file",
			"errors":     err.Error(),
			"request_id": requestId,
		})
	}
	defer func() {
		err := photoFile.Close()
		if err != nil {
			log.Error().Err(err).Str("request_id", requestId).Msg("Closing photo file")
		}
	}()

	err = s.ticketDomain.StorePaymentReceipt(c.Request().Context(), email, photoFile, photoContentType)
	if err != nil {
		var validationError *ValidationError
		if errors.As(err, &validationError) {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message":    "Validation error",
				"errors":     validationError.Error(),
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
}
