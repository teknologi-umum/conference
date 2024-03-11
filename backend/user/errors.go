package user

import (
	"errors"
	"strings"
)

type ValidationError struct {
	Errors []string
}

func (v ValidationError) Error() string {
	return strings.Join(v.Errors, ", ")
}

var ErrUserEmailNotFound = errors.New("user email not found")
