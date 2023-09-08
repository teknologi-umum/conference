package main

import (
	"bytes"
	"net"
	"net/smtp"
	"strconv"
	"time"
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

func (m *Mailer) messageBuilder(mail *Mail) []byte {
	// Refer to https://www.rfc-editor.org/rfc/rfc5322
	var msg bytes.Buffer
	msg.WriteString("MIME-version: 1.0\n")
	msg.WriteString("Date: ")
	msg.WriteString(time.Now().Format(time.RFC1123Z))
	msg.WriteString("\n")
	msg.WriteString("From: \"Teknologi Umum Conference\" <conference@teknologiumum.com>\n")
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

	return msg.Bytes()
}

func (m *Mailer) Send(mail *Mail) error {
	if m.configuration.SmtpFrom == "" || m.configuration.SmtpPassword == "" {
		return smtp.SendMail(
			net.JoinHostPort(m.configuration.SmtpHostname, m.configuration.SmtpPort),
			nil, // no auth
			"conference@teknologiumum.com",
			[]string{mail.RecipientEmail},
			m.messageBuilder(mail),
		)
	}

	return smtp.SendMail(
		net.JoinHostPort(m.configuration.SmtpHostname, m.configuration.SmtpPort),
		smtp.PlainAuth("", m.configuration.SmtpFrom, m.configuration.SmtpPassword, m.configuration.SmtpHostname),
		"conference@teknologiumum.com",
		[]string{mail.RecipientEmail},
		m.messageBuilder(mail),
	)
}
