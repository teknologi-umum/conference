package main

import (
	"encoding/csv"
	"errors"
	"io"
	"strings"

	"conf/user"
)

func csvReader(file string, mandatoryNameField bool) (users []user.User, err error) {
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
		var name, email string
		if mandatoryNameField {
			if m["name"] == "" {
				err = errors.New("name is required")
				return nil, err
			}
			name = m["name"]
		}

		if m["email"] == "" {
			err = errors.New("email is required")
			return nil, err
		}
		email = m["email"]
		users = append(users, user.User{
			Name:  name,
			Email: email,
		})
	}
	return
}
