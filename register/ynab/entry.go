package ynab

import (
	"fmt"
	"github.com/google/uuid"
	"main/register"
	iface "main/register/entry"
	"time"
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

func NewEntry(s []string) (*entry, error) {
	if s[acct] != "Primary Checking" {
		return nil, nil
	}

	var err error
	errMsg := `couldn't generate entry: %v\n`
	e := new(entry)

	e.payee = s[payee]
	e.date, err = register.ParseDate(s[date])
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	e.inflow, err = register.ParseCurrency(s[inflow])
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	of, err := register.ParseCurrency(s[outflow])
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	// YNAB stores outflows at positive numbers, so we negate them here
	if of > 0 {
		of = 0 - of
	}
	e.outflow = of

	e.uuid = uuid.New()

	return e, nil
}

var _ iface.Entry = &entry{}

type entry struct {
	date    time.Time
	inflow  float64
	outflow float64
	payee   string
	cleared bool
	uuid    uuid.UUID
}

func (e *entry) ID() uuid.UUID {
	return e.uuid
}

func (e *entry) IsCleared() bool {
	return e.cleared
}

func (e *entry) Clear() {
	if ! e.cleared {
		e.cleared = true
	}
}

func (e *entry) Date() time.Time {
	return e.date
}

func (e *entry) Payee() string {
	return e.payee
}

func (e *entry) Amount() float64 {
	if e.outflow < 0 {
		return e.outflow
	}
	return e.inflow
}

