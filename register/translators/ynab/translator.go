package ynab

import (
	"encoding/csv"
	"fmt"
	"io"
	"main/register"
	"main/register/entry"
	iface "main/register/translator"
)

const (
	acct = iota
	_
	date
	payee
	_
	_
	_
	_
	outflow
	inflow
	_
)

var _ iface.Translator = &translator{}

type translator struct {
	account string
}

func NewTranslator(account string) *translator {
	return &translator{account: account}
}

// Translate takes a reader who's output is expected to be CSV formated
// returns a slice of entries or an error
func (t *translator) Translate(r io.Reader) ([]*entry.Entry, error) {
	rdr := csv.NewReader(r)
	rdr.LazyQuotes = true
	return t.readRecordsToEntries(rdr)
}

func (t *translator) readRecordsToEntries(r *csv.Reader) (entries []*entry.Entry, err error) {
	entries = make([]*entry.Entry, 0, 4096)
	var e *entry.Entry
	for {
		e, err = t.readRecordToEntry(r)
		if err != nil {
			break
		}
		entries = append(entries, e)
	}
	if err != io.EOF {
		return nil, fmt.Errorf("could not read records to entries: %v", err)
	}
	return entries, nil
}

func (t *translator) readRecordToEntry(r *csv.Reader) (*entry.Entry, error) {
	rec, err := r.Read()
	if err != nil {
		return nil, err
	}
	return t.parseRecord(rec)
}

func (t *translator) parseRecord(record []string) (*entry.Entry, error) {
	// Ynab CSVs contain all tracked accounts as one, ignore all but the specified one
	if record[acct] != t.account {
		return nil, newIsNotAccountError(t.account, record[acct])
	}
	errMsg := `error parsing record`
	e := entry.NewEntry()

	e.SetPayee(record[payee])
	d, err := register.ParseDate(record[date])
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	e.SetDate(d)

	amt, err := register.ParseCurrency(record[inflow])
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	if amt == 0 {
		amt, err = register.ParseCurrency(record[outflow])
		if err != nil {
			return nil, fmt.Errorf(errMsg, err)
		}
		// YNAB stores outflows at positive numbers, manually negate them here
		amt = -amt
	}
	e.SetAmount(amt)
	return e, nil
}
