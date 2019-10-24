package usaa

import (
	"fmt"
	"github.com/google/uuid"
	"main/register"
	iface "main/register/entry"
	"strings"
	"time"
)

var _ iface.Entry = &entry{}

type entry struct {
	date    time.Time
	payee   string
	amount  float64
	cleared bool
	uuid    uuid.UUID
}

const (
	posted = iota
	_
	date
	_
	payee
	category
	amount
)

func newEntry(s []string) (*entry, error) {
	var err error
	errMsg := `couldn't generate entry: %v\n'`
	e := new(entry)

	e.payee = s[payee]
	e.date, err = register.ParseDate(s[date])
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	// USAA uses double negation to show positive values.  Makes perfect sense.
	amt := strings.TrimPrefix(s[amount], "--")
	e.amount, err = register.ParseCurrency(amt)
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	e.uuid = uuid.New()
	return e, nil
}

func (e *entry) Clear() {
	if ! e.cleared {
		e.cleared = true
	}
}

func (e *entry) IsCleared() bool {
	return e.cleared
}

func (e *entry) Date() time.Time {
	return e.date
}

func (e *entry) Payee() string {
	return e.payee
}

func (e *entry) Amount() float64 {
	return e.amount
}

func (e *entry) ID() uuid.UUID {
	return e.uuid
}