package usaa

import (
	"fmt"
	"main/register"
	"main/register/entry"
	"main/register/translator"
	"strings"
)

const (
	posted = iota
	_
	date
	_
	payee
	category
	amount
)

type usaaTranslator struct{}

var _ translator.Translator = &usaaTranslator{}

func NewTranslator() *usaaTranslator {
	return &usaaTranslator{}
}

func (u *usaaTranslator) ToEntry(s []string) (*entry.Entry, error) {
	var err error
	errMsg := `couldn't generate entry: %v\n'`
	e := entry.NewEntry()

	e.SetPayee(s[payee])
	d, err := register.ParseDate(s[date])
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	e.SetDate(d)
	// USAA uses double negation to show positive values.  Makes perfect sense.
	a, err := register.ParseCurrency(strings.TrimPrefix(s[amount], "--"))
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	e.SetAmount(a)
	return e, nil
}
