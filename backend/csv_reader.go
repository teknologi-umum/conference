package main

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"

	"github.com/rs/zerolog/log"
)

func csvReader(file string) (users []User, err error) {
	r := csv.NewReader(strings.NewReader(file))
	header, err := r.Read()
	if err != nil {
		err = errors.New("failed to read csv header")
		return
	}
	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		m := make(map[string]string)
		for i, h := range header {
			m[h] = record[i]
		}
		if m["name"] == "" {
			log.Fatal().Msg("Username is required")
		}
		if m["email"] == "" {
			log.Fatal().Msg("Email is required")
		}
		users = append(users, User{
			Name:  m["name"],
			Email: m["email"],
		})
	}
	return
}
