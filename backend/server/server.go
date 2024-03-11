package server

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"conf/administrator"
	"conf/features"
	"conf/mailer"
	"conf/ticketing"
	"conf/user"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

type ServerConfig struct {
	UserDomain          *user.UserDomain
	TicketDomain        *ticketing.TicketDomain
	AdministratorDomain *administrator.AdministratorDomain
	FeatureFlag         *features.FeatureFlag
	MailSender          *mailer.Mailer
	Environment         string
	ValidateTicketKey   string
	Hostname            string
	Port                string
}

type ServerDependency struct {
	userDomain          *user.UserDomain
	ticketDomain        *ticketing.TicketDomain
	administratorDomain *administrator.AdministratorDomain
	featureFlag         *features.FeatureFlag
	mailSender          *mailer.Mailer
	validateTicketKey   string
}

func NewServer(config *ServerConfig) (*http.Server, error) {
	if config.UserDomain == nil {
		return nil, fmt.Errorf("nil UserDomain")
	}

	if config.TicketDomain == nil {
		return nil, fmt.Errorf("nil TicketDomain")
	}

	if config.AdministratorDomain == nil {
		return nil, fmt.Errorf("nil AdministratorDomain")
	}

	if config.FeatureFlag == nil {
		return nil, fmt.Errorf("nil FeatureFlag")
	}

	if config.MailSender == nil {
		return nil, fmt.Errorf("nil MailSender")
	}

	if config.ValidateTicketKey == "" {
		return nil, fmt.Errorf("nil ValidateTicketKey")
	}

	dependencies := &ServerDependency{
		userDomain:          config.UserDomain,
		ticketDomain:        config.TicketDomain,
		administratorDomain: config.AdministratorDomain,
		featureFlag:         config.FeatureFlag,
		mailSender:          config.MailSender,
		validateTicketKey:   config.ValidateTicketKey,
	}

	r := chi.NewRouter()

	r.Use(sentryhttp.New(sentryhttp.Options{Repanic: false}).Handle)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	// NOTE: Only need to handle CORS, everything else is being handled by the API gateway
	corsAllowedOrigins := []string{"https://conference.teknologiumum.com", "https://conf.teknologiumum.com"}
	if config.Environment != "production" {
		corsAllowedOrigins = append(corsAllowedOrigins, "http://localhost:3000")
	}
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   corsAllowedOrigins,
		AllowedMethods:   []string{http.MethodPost},
		AllowedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           3600, // 1 day
	}).Handler)

	r.Use(middleware.Heartbeat("/api/public/ping"))

	r.Post("/api/public/register-user", dependencies.RegisterUser)
	r.Post("/api/public/upload-payment-proof", dependencies.UploadPaymentProof)
	r.Post("/api/public/scan-ticket", dependencies.DayTicketScan)

	r.Post("/api/administrator/login", dependencies.AdministratorLogin)

	return &http.Server{
		Addr:              net.JoinHostPort(config.Hostname, config.Port),
		Handler:           r,
		ReadTimeout:       time.Minute * 3,
		ReadHeaderTimeout: time.Minute,
		WriteTimeout:      time.Hour,
		IdleTimeout:       time.Hour,
	}, nil
}
