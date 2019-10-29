package entry

import (
	"github.com/google/uuid"
	"time"
)

type Entry struct {
	date    time.Time
	payee   string
	amount  float64
	cleared bool
	id      uuid.UUID
}

func (e *Entry) Cleared() bool {
	return e.cleared
}

func (e *Entry) SetCleared() {
	if !e.cleared {
		e.cleared = true
	}
}

func (e *Entry) Id() uuid.UUID {
	return e.id
}

func (e *Entry) SetId(id uuid.UUID) {
	e.id = id
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
	return e.date
}

func (e *Entry) SetDate(date time.Time) {
	e.date = date
}

func NewEntry() *Entry {
	return &Entry{
		id: uuid.New(),
	}
}
