package main_test

import (
	"os"
	"testing"

	main "conf"
)

func TestMailSender(t *testing.T) {
	smtpHostname, ok := os.LookupEnv("SMTP_HOSTNAME")
	if !ok {
		smtpHostname = "localhost"
	}
	smtpPort, ok := os.LookupEnv("SMTP_PORT")
	if !ok {
		smtpPort = "1025"
	}
	smtpFrom, ok := os.LookupEnv("SMTP_FROM")
	if !ok {
		smtpFrom = ""
	}
	smtpPassword, ok := os.LookupEnv("SMTP_PASSWORD")
	if !ok {
		smtpPassword = ""
	}

	mailSender := main.NewMailSender(&main.MailConfiguration{
		SmtpHostname: smtpHostname,
		SmtpPort:     smtpPort,
		SmtpFrom:     smtpFrom,
		SmtpPassword: smtpPassword,
	})

	t.Run("Happy Scenario", func(t *testing.T) {
		mail := &main.Mail{
			RecipientName:  "John Doe",
			RecipientEmail: "johndoe@example.com",
			Subject:        "Welcome to TeknumConf, you are on the waiting list",
			PlainTextBody: `Hello, {{ name }}!
Thank you for participating on TeknumConf, due to the limited seating quota, you are on a waitlist.
Not to worry, you will receive an email from us regarding the seating in about 7 days.
Please do contact us if you didn't receive any email by then.`,
			HtmlBody: `<!DOCTYPE html>
<html lang="en" xmlns="http://www.w3.org/1999/xhtml"
>
<head>
    <meta content="IE=edge" http-equiv="X-UA-Compatible">
    <meta content="width=device-width,initial-scale=1 user-scalable=yes" name="viewport">
    <meta content="telephone=no, date=no, address=no, email=no, url=no" name="format-detection">
    <meta name="x-apple-disable-message-reformatting">
    <meta content="light dark" name="color-scheme">
    <meta content="light dark" name="supported-color-schemes">
    <meta charset="UTF-8">
    <!--[if mso]>
    <noscript>
        <xml>
            <o:OfficeDocumentSettings>
                <o:PixelsPerInch>96</o:PixelsPerInch>
            </o:OfficeDocumentSettings>
        </xml>
    </noscript> <![endif]-->

    <style>
        :root {
            color-scheme: light dark;
            supported-color-schemes: light dark;
        }
    </style>
    <title>TeknumConf - Attendee Waitlist</title>
</head>
<body>
<h1>Hello, {{ name }}!</h1>
<p>Thank you for participating on TeknumConf, due to the limited seating quota, you are on a waitlist.</p>
<p>Not to worry, you will receive an email from us regarding the seating in about 7 days.
Please do contact us if you didn't receive any email by then.</p>
</body>
</html>`,
		}

		err := mailSender.Send(mail)
		if err != nil {
			t.Errorf("unexpected error: %s", err.Error())
		}
	})
}
