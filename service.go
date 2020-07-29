package main

import (
	"errors"
	"strings"
)

//StringService interface provides an  apbstraction of the business logic
type StringService interface {
	Uppercase(string) (string, error)
	Count(string) int
}

//StringService implementation
type stringService struct{}

//ErrEmpty is returned when the input string is empty
var ErrEmpty = errors.New("Empty string")

func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}
