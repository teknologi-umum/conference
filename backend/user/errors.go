package user

import "strings"

type ValidationError struct {
	Errors []string
}

func (v ValidationError) Error() string {
	return strings.Join(v.Errors, ", ")
}
