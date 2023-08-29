#!/bin/bash
go run . blast-email --smtp.hostname=127.0.0.1 --smtp.port=1025 --smtp.from=admin@localhost --smtp.password="" \
--subject="TeknumConf - Attendee Waitlist" --plaintext-body ../emails/attendee_waitlist.txt --html-body ../emails/attendee_waitlist.html --recipients ../emails/attendees-sample.csv
