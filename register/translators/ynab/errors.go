package ynab

import "fmt"

func isExpected(err error) bool {
	exErr, ok := err.(expectedError)
	return ok && exErr.expected()
}

type expectedError interface {
	expected() bool
}

type notAccountError struct {
	w, g string
}

func newIsNotAccountError(want, got string) *notAccountError {
	return &notAccountError{want, got}
}

func (e notAccountError) Error() string {
	return fmt.Sprintf("want account %q, got account %q", e.w, e.g)
}

func (e notAccountError) expected() bool {
	return true
}
