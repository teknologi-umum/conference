package main_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"gocloud.dev/blob"
	_ "gocloud.dev/blob/fileblob"

	main "conf"
)

var database *pgxpool.Pool
var bucket *blob.Bucket
var mailSender *main.Mailer

func TestMain(m *testing.M) {
	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		databaseUrl = "postgres://postgres:password@localhost:5432/conf?sslmode=disable"
	}

	tempDir, err := os.MkdirTemp(os.TempDir(), "teknologi-umum-conference")
	if err != nil {
		log.Fatalf("creating temporary directory: %s", err.Error())
		return
	}

	blobUrl, ok := os.LookupEnv("BLOB_URL")
	if !ok {
		blobUrl = "file://" + tempDir
	}

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

	database, err = pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		log.Fatalf("creating pgx pool instance: %s", err.Error())
		return
	}

	bucket, err = blob.OpenBucket(context.Background(), blobUrl)
	if err != nil {
		log.Fatalf("creating bucket instance: %s", err.Error())
	}

	mailSender = main.NewMailSender(&main.MailConfiguration{
		SmtpHostname: smtpHostname,
		SmtpPort:     smtpPort,
		SmtpFrom:     smtpFrom,
		SmtpPassword: smtpPassword,
	})

	// Migrate database
	if err := main.App().Run([]string{"teknum-conf", "migrate", "up"}); err != nil {
		log.Fatalf("migrating database: %s", err.Error())
		return
	}

	exitCode := m.Run()

	// Migrate database
	if err := main.App().Run([]string{"teknum-conf", "migrate", "down"}); err != nil {
		log.Fatalf("migrating database: %s", err.Error())
		return
	}

	os.RemoveAll(tempDir)
	bucket.Close()
	database.Close()

	os.Exit(exitCode)
}
