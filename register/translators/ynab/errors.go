package ynab

import (
	"fmt"
	tranerr "main/register/translators/translator_errors"
)

var _ tranerr.ExpectedError = &notAccountError{}

type notAccountError struct {
	w, g string
}

func (e *notAccountError) IsExpected(err error) (expected bool) {
	switch err.(type) {
	case tranerr.ExpectedError:
		return true
	}
	return
}

func newIsNotAccountError(want, got string) *notAccountError {
	return &notAccountError{want, got}
}

func (e notAccountError) Error() string {
	return fmt.Sprintf("want account %q, got account %q", e.w, e.g)
}

func isExpected(err error) bool {
	switch err.(type) {
	case tranerr.ExpectedError:
		return true
	}
	return false
}
