package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"net"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
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
	Attachments    []Attachment
}

type ContentDisposition uint8

const (
	ContentDispositionUnspecified ContentDisposition = iota
	ContentDispositionInline
	ContentDispositionAttachment
)

type Attachment struct {
	Name               string
	Description        string
	ContentType        string
	ContentDisposition ContentDisposition
	ContentId          string
	SHA256Checksum     []byte
	// Payload must be a type of raw binary. Do not encode this to hex or base64.
	Payload []byte
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

func (m *Mailer) messageBuilder(ctx context.Context, mail *Mail) []byte {
	span := sentry.StartSpan(ctx, "mailer.message_builder")
	defer span.Finish()

	// Refer to https://www.rfc-editor.org/rfc/rfc5322
	mixedBoundary := strings.ReplaceAll(uuid.NewString(), "-", "")
	relatedBoundary := strings.ReplaceAll(uuid.NewString(), "-", "")
	alternateBoundary := strings.ReplaceAll(uuid.NewString(), "-", "")

	var msg bytes.Buffer
	msg.WriteString("MIME-version: 1.0\n")
	msg.WriteString("Date: ")
	msg.WriteString(time.Now().Format(time.RFC1123Z))
	msg.WriteString("\n")
	msg.WriteString("From: \"Teknologi Umum Conference\" <conference@teknologiumum.com>\n")
	msg.WriteString("To: " + strconv.Quote(mail.RecipientName) + " <" + mail.RecipientEmail + ">\n")
	msg.WriteString("Subject: " + mail.Subject + "\n")
	msg.WriteString("Content-Type: multipart/mixed; boundary= " + strconv.Quote(mixedBoundary) + "\n\n")
	msg.WriteString("--" + mixedBoundary + "\n")
	msg.WriteString("Content-Type: multipart/related; boundary=" + strconv.Quote(relatedBoundary) + "\n\n")
	msg.WriteString("--" + relatedBoundary + "\n")
	msg.WriteString("Content-Type: multipart/alternative; boundary=" + strconv.Quote(alternateBoundary) + "\n\n")
	msg.WriteString("--" + alternateBoundary + "\n")
	msg.WriteString("Content-Type: text/plain; charset=\"us-ascii\"\n")
	msg.WriteString("Content-Transfer-Encoding: 8bit\n")
	msg.WriteString("\n")
	msg.WriteString(mail.PlainTextBody)
	msg.WriteString("\n")
	msg.WriteString("--" + alternateBoundary + "\n")
	msg.WriteString("Content-Type: text/html; charset=\"utf-8\"\n")
	msg.WriteString("Content-Transfer-Encoding: 8bit\n")
	msg.WriteString("\n")
	msg.WriteString(mail.HtmlBody)
	msg.WriteString("\n")
	msg.WriteString("--" + alternateBoundary + "--\n")
	for _, attachment := range mail.Attachments {
		if attachment.ContentDisposition == ContentDispositionInline {
			msg.WriteString("\n--" + relatedBoundary + "\n")
			msg.WriteString("Content-Type: " + attachment.ContentType + "\n")
			msg.WriteString("Content-Disposition: inline; filename=" + strconv.Quote(attachment.Name) + ";\n")
			msg.WriteString("Content-Description: " + attachment.Description + "\n")
			msg.WriteString("Content-ID: <" + attachment.ContentId + ">\n")
			msg.WriteString("Content-Transfer-Encoding: base64\n")
			msg.WriteString("\n")
			msg.WriteString(base64.StdEncoding.EncodeToString(attachment.Payload))
			msg.WriteString("\n")
			msg.WriteString("\n--" + relatedBoundary + "--\n")
		}

		if attachment.ContentDisposition == ContentDispositionAttachment {
			msg.WriteString("\n--" + mixedBoundary + "\n")
			msg.WriteString("Content-Type: " + attachment.ContentType + "\n")
			msg.WriteString("Content-Disposition: attachment; filename=" + strconv.Quote(attachment.Name) + ";\n")
			msg.WriteString("Content-Description: " + attachment.Description + "\n")
			msg.WriteString("Content-Transfer-Encoding: base64\n")
			msg.WriteString("\n")
			msg.WriteString(base64.StdEncoding.EncodeToString(attachment.Payload))
			msg.WriteString("\n")
			msg.WriteString("\n--" + mixedBoundary + "--\n")
		}
	}
	return msg.Bytes()
}

func (m *Mailer) Send(ctx context.Context, mail *Mail) error {
	span := sentry.StartSpan(ctx, "mailer.send")
	defer span.Finish()

	if m.configuration.SmtpFrom == "" || m.configuration.SmtpPassword == "" {
		return smtp.SendMail(
			net.JoinHostPort(m.configuration.SmtpHostname, m.configuration.SmtpPort),
			nil, // no auth
			"conference@teknologiumum.com",
			[]string{mail.RecipientEmail},
			m.messageBuilder(ctx, mail),
		)
	}

	return smtp.SendMail(
		net.JoinHostPort(m.configuration.SmtpHostname, m.configuration.SmtpPort),
		smtp.PlainAuth("", m.configuration.SmtpFrom, m.configuration.SmtpPassword, m.configuration.SmtpHostname),
		"conference@teknologiumum.com",
		[]string{mail.RecipientEmail},
		m.messageBuilder(ctx, mail),
	)
}
