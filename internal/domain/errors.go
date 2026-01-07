package domain

import "errors"

var (
	ErrAlreadyExists = errors.New("this alias is already taken")
)
