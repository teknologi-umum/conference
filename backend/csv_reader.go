package main

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"

	"github.com/rs/zerolog/log"
)

func csvReader(file string) (err error, users []User) {
	var userList []User
	r := csv.NewReader(strings.NewReader(file))
	header, err := r.Read()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read csv header")
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
		userList = append(userList, User{
			Name:  m["name"],
			Email: m["email"],
		})
	}
	return nil, userList
}
