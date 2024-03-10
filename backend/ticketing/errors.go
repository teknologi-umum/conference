package ticketing

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

var ErrInvalidTicket = errors.New("invalid ticket")
var ErrUserEmailNotFound = errors.New("user email not found")
