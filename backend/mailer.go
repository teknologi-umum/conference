package main

import (
	"bytes"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type MailConfiguration struct {
	SmtpHostname string
	SmtpPort     string
	SmtpFrom     string
	SmtpPassword string
}

type Mail struct {
	RecipientName  string
	RecipientEmail string
	Subject        string
	PlainTextBody  string
	HtmlBody       string
}

type Mailer struct {
	configuration *MailConfiguration
}

func NewMailSender(configuration *MailConfiguration) *Mailer {
	mailer := &Mailer{
		configuration: configuration,
	}

	return mailer
}

func (m *Mailer) messageBuilder(mail *Mail) string {
	// Refer to https://www.rfc-editor.org/rfc/rfc5322
	var msg bytes.Buffer
	msg.WriteString("MIME-version: 1.0\n")
	msg.WriteString("Date: ")
	msg.WriteString(time.Now().Format(time.RFC1123Z))
	msg.WriteString("\n")
	msg.WriteString("From: \"Teknologi Umum Conference\" <" + m.configuration.SmtpFrom + ">\n")
	msg.WriteString("To: " + strconv.Quote(mail.RecipientName) + " <" + mail.RecipientEmail + ">\n")
	msg.WriteString("Subject: " + mail.Subject + "\n")
	msg.WriteString("Content-Type: multipart/mixed; boundary=\"mixed_boundary\"\n\n")
	msg.WriteString("--mixed_boundary\n")
	msg.WriteString("Content-Type: multipart/related; boundary=\"related_boundary\"\n\n")
	msg.WriteString("--related_boundary\n")
	msg.WriteString("Content-Type: multipart/alternative; boundary=\"alternative_boundary\"\n\n")
	msg.WriteString("--alternative_boundary\n")
	msg.WriteString("Content-Type: text/plain; charset=\"us-ascii\"\n")
	msg.WriteString("Content-Transfer-Encoding: 8bit\n")
	msg.WriteString("\n")
	msg.WriteString(mail.PlainTextBody)
	msg.WriteString("\n")
	msg.WriteString("--alternative_boundary\n")
	msg.WriteString("Content-Type: text/html; charset=\"utf-8\"\n")
	msg.WriteString("Content-Transfer-Encoding: 8bit\n")
	msg.WriteString("\n")
	msg.WriteString(mail.HtmlBody)
	msg.WriteString("\n")
	msg.WriteString("--alternative_boundary--\n")

	return msg.String()
}

func (m *Mailer) Send(mail *Mail) error {
	client, err := smtp.Dial(net.JoinHostPort(m.configuration.SmtpHostname, m.configuration.SmtpPort))
	if err != nil {
		return fmt.Errorf("dialing: %w", err)
	}
	defer func(client *smtp.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg("Closing smtp client connection")
		}
	}(client)

	err = client.Mail(m.configuration.SmtpFrom)
	if err != nil {
		return fmt.Errorf("sending mail from: %w", err)
	}

	// If password is not empty, we'll try and use authentication
	if m.configuration.SmtpPassword != "" {
		err := client.Auth(smtp.CRAMMD5Auth(m.configuration.SmtpFrom, m.configuration.SmtpPassword))
		if err != nil {
			// If CRAM-MD5 fails, we'll try plain auth
			e := client.Auth(smtp.PlainAuth("", m.configuration.SmtpFrom, m.configuration.SmtpPassword, m.configuration.SmtpHostname))
			if e != nil {
				return fmt.Errorf("authenticating: %w", e)
			}
		}
	}

	err = client.Rcpt(mail.RecipientEmail)
	if err != nil {
		return fmt.Errorf("sending mail to: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("sending data: %w", err)
	}

	_, err = fmt.Fprint(wc, m.messageBuilder(mail))
	if err != nil {
		return fmt.Errorf("sending data: %w", err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("closing connection: %w", err)
	}

	return nil
}
