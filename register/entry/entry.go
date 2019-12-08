package entry

import (
	"time"
)

type Entry struct {
	time      time.Time
	payee     string
	amount    float64
	cleared   bool
	importID  string
	occurance int
}

func (e *Entry) Occurance() int {
	return e.occurance
}

func (e *Entry) SetOccurance(occurance int) {
	e.occurance = occurance
}

func (e *Entry) ImportID() string {
	return e.importID
}

func (e *Entry) SetImportID(importID string) {
	e.importID = importID
}

func (e *Entry) Cleared() bool {
	return e.cleared
}

func (e *Entry) SetCleared() {
	if !e.cleared {
		e.cleared = true
	}
}

func (e *Entry) Amount() float64 {
	return e.amount
}

func (e *Entry) SetAmount(amount float64) {
	e.amount = amount
}

func (e *Entry) Payee() string {
	return e.payee
}

func (e *Entry) SetPayee(payee string) {
	e.payee = payee
}

func (e *Entry) Date() time.Time {
	return e.time
}

func (e *Entry) SetDate(date time.Time) {
	e.time = date
}

func NewEntry() *Entry {
	return &Entry{}
}
