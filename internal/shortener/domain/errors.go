package domain

import "errors"

var (
	ErrShortNotFound = errors.New("shorten url not found")
	ErrShortCodeAlreadyExists = errors.New("shorten code already exists")
)
