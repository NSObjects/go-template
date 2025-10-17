package user

import (
	"errors"
	"strings"
)

var (
	ErrInvalidUsername = errors.New("username cannot be empty")
	ErrInvalidEmail    = errors.New("email cannot be empty")
	ErrInvalidAge      = errors.New("age must be non-negative")
)

type ID int64

func (id ID) Int64() int64 {
	return int64(id)
}

type Username struct {
	value string
}

func NewUsername(value string) (Username, error) {
	if strings.TrimSpace(value) == "" {
		return Username{}, ErrInvalidUsername
	}
	return Username{value: value}, nil
}

func (u Username) String() string {
	return u.value
}

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	if strings.TrimSpace(value) == "" {
		return Email{}, ErrInvalidEmail
	}
	return Email{value: value}, nil
}

func (e Email) String() string {
	return e.value
}

type Age struct {
	value int
}

func NewAge(value int) (Age, error) {
	if value < 0 {
		return Age{}, ErrInvalidAge
	}
	return Age{value: value}, nil
}

func (a Age) Int() int {
	return a.value
}
