package main

import (
	"fmt"
	"os"

	"conf/mailer"
	"conf/user"
	"github.com/flowchartsman/handlebars/v3"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func BlastMailHandlerAction(cCtx *cli.Context) error {
	config, err := GetConfig(cCtx.String("config-file-path"))
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn:              "",
		Debug:            config.Environment != "production",
		AttachStacktrace: true,
		SampleRate:       1.0,
		Release:          version,
		Environment:      config.Environment,
		DebugWriter:      log.Logger,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if config.Environment != "production" {
				log.Debug().Interface("exceptions", event.Exception).Msg(event.Message)
			}

			return event
		},
	})
	if err != nil {
		return fmt.Errorf("initializing Sentry: %w", err)
	}

	subject := cCtx.String("subject")
	plaintext := cCtx.String("plaintext-body")
	htmlBody := cCtx.String("html-body")
	mailCsv := cCtx.String("recipients")
	singleRecipient := cCtx.String("single-recipient")

	if subject == "" {
		log.Fatal().Msg("Subject is required")
	}
	if plaintext == "" {
		log.Fatal().Msg("Plaintext template is required")
	}
	if htmlBody == "" {
		log.Fatal().Msg("Html template is required")
	}
	if mailCsv == "" && singleRecipient == "" {
		log.Fatal().Msg("Recipient is required")
	}

	plaintextContent, err := os.ReadFile(plaintext)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read plaintext template")
	}

	plaintextTemplate, err := handlebars.Parse(string(plaintextContent))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse plaintext template")
	}

	htmlContent, err := os.ReadFile(htmlBody)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read html template")
	}

	htmlTemplate, err := handlebars.Parse(string(htmlContent))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse html template")
	}

	var userList []user.User

	if mailCsv != "" {
		emailList, err := os.ReadFile(mailCsv)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to read email list")
		}

		userList, err = csvReader(string(emailList), true)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to parse email list")
		}
	} else {
		userList = append(userList, user.User{
			Email: singleRecipient,
		})
	}

	mailSender := mailer.NewMailSender(&mailer.MailConfiguration{
		SmtpHostname: config.Mailer.Hostname,
		SmtpPort:     config.Mailer.Port,
		SmtpFrom:     config.Mailer.From,
		SmtpPassword: config.Mailer.Password,
	})

	for _, userItem := range userList {
		mail := &mailer.Mail{
			RecipientName:  userItem.Name,
			RecipientEmail: userItem.Email,
			Subject:        subject,
			PlainTextBody:  string(plaintextContent),
			HtmlBody:       string(htmlContent),
		}

		// Parse email template information
		emailTemplate := map[string]any{
			"ticketPrice":                         config.EmailTemplate.TicketPrice,
			"ticketStudentCollegePrice":           config.EmailTemplate.TicketStudentCollegePrice,
			"ticketStudentHighSchoolPrice":        config.EmailTemplate.TicketStudentHighSchoolPrice,
			"ticketStudentCollegeDiscount":        config.EmailTemplate.TicketStudentCollegeDiscount,
			"ticketStudentHighSchoolDiscount":     config.EmailTemplate.TicketStudentHighSchoolDiscount,
			"percentageStudentCollegeDiscount":    config.EmailTemplate.PercentageStudentCollegeDiscount,
			"percentageStudentHighSchoolDiscount": config.EmailTemplate.PercentageStudentHighSchoolDiscount,
			"conferenceEmail":                     config.EmailTemplate.ConferenceEmail,
			"bankAccounts":                        config.EmailTemplate.BankAccounts,
		}
		// Execute handlebars template only if userItem.Name is not empty
		if userItem.Name != "" {
			emailTemplate["name"] = userItem.Name
		}

		mail.PlainTextBody = plaintextTemplate.MustExec(emailTemplate)
		mail.HtmlBody = htmlTemplate.MustExec(emailTemplate)

		err := mailSender.Send(cCtx.Context, mail)
		if err != nil {
			log.Error().Err(err).Msgf("failed to send email to %s", userItem.Email)
			continue
		}

		log.Info().Msgf("Sending email to %s", userItem.Email)
	}
	log.Info().Msg("Blasting email done")
	return nil
}
