package ynab

import (
	"fmt"
	"main/register"
	"main/register/entry"
	"main/register/translator"
	"strings"
)

const (
	acct = iota
	flag
	date
	payee
	groupCategory
	group
	category
	memo
	outflow
	inflow
	cleared
)

type ynabTranslator struct {
	account string
}

var _ translator.Translator = &ynabTranslator{}

func NewTranslator(account string) *ynabTranslator {
	return &ynabTranslator{
		account: account,
	}
}

const header string = "Account"

func (t *ynabTranslator) ToEntry(s []string) (*entry.Entry, error) {
	// Ynab CSVs contain all tracked accounts as one, ignore all but the specified one
	if s[acct] != t.account || strings.Contains(s[acct], header) {
		return nil, nil
	}

	errMsg := `couldn't generate entry: %v\n`
	e := entry.NewEntry()

	e.SetPayee(s[payee])
	d, err := register.ParseDate(s[date])
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	e.SetDate(d)

	amt, err := register.ParseCurrency(s[inflow])
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	if amt == 0 {
		amt, err = register.ParseCurrency(s[outflow])
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
		// YNAB stores outflows at positive numbers, so we negate them here
		if amt > 0 {
			amt = 0 - amt
		}
	}
	e.SetAmount(amt)
	return e, nil

}
