package core

import "errors"

var (
	ErrValidation = errors.New("validation")
)

type Errors []Error

type Error struct {
	Key string
	Err error
}
