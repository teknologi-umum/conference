package server

import (
	"net"
	"net/http"
	"time"

	"conf/ticketing"
	"conf/user"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

type ServerConfig struct {
	UserDomain                *user.UserDomain
	TicketDomain              *ticketing.TicketDomain
	Environment               string
	FeatureRegistrationClosed bool
	ValidateTicketKey         string
	Hostname                  string
	Port                      string
}

type ServerDependency struct {
	userDomain         *user.UserDomain
	ticketDomain       *ticketing.TicketDomain
	registrationClosed bool
	validateTicketKey  string
}

func NewServer(config *ServerConfig) *http.Server {
	if config.UserDomain == nil || config.TicketDomain == nil {
		// For production backend application, please don't do what I just did.
		// Do a proper nil check and validation for each of your config and dependencies.
		// NEVER call panic(), just return error.
		// I'm in a hackathon (basically in a rush), so I'm doing this.
		// Let me remind you again: don't do what I just did.
		panic("one of the domain dependency is nil")
	}

	dependencies := &ServerDependency{
		userDomain:         config.UserDomain,
		ticketDomain:       config.TicketDomain,
		registrationClosed: config.FeatureRegistrationClosed,
		validateTicketKey:  config.ValidateTicketKey,
	}

	r := chi.NewRouter()

	r.Use(sentryhttp.New(sentryhttp.Options{Repanic: false}).Handle)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	// NOTE: Only need to handle CORS, everything else is being handled by the API gateway
	corsAllowedOrigins := []string{"https://conference.teknologiumum.com"}
	if config.Environment != "production" {
		corsAllowedOrigins = append(corsAllowedOrigins, "http://localhost:3000")
	}
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   corsAllowedOrigins,
		AllowedMethods:   []string{http.MethodPost},
		AllowCredentials: false,
		MaxAge:           3600, // 1 day
	}).Handler)

	r.Use(middleware.Heartbeat("/api/ping"))

	r.Post("/users", dependencies.RegisterUser)
	r.Post("/bukti-transfer", dependencies.UploadPaymentProof)
	r.Post("/scan-tiket", dependencies.DayTicketScan)

	return &http.Server{
		Addr:              net.JoinHostPort(config.Hostname, config.Port),
		Handler:           r,
		ReadTimeout:       time.Minute * 3,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Hour,
		IdleTimeout:       time.Hour,
	}
}
