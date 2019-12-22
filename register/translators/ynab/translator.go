package ynab

import (
	"encoding/csv"
	"fmt"
	"io"
	"main/register"
	"main/register/entry"
	iface "main/register/translator"
	"strings"
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
		if e == nil && err != nil && !isExpected(err){
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
	// YNAB CSVs contain all tracked accounts, ignore entries for other accounts
	if t.isHeader(record) {
		return nil, nil
	}
	if record[acct] != t.account {
		return nil, newIsNotAccountError(t.account, record[acct])
	}
	errMsg := `error parsing record: %v`
	e := entry.NewEntry()

	e.SetPayee(record[payee])

	const dateFormat = "01/02/2006"

	d, err := register.ParseDate(record[date], dateFormat)
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

func (t *translator) isHeader(record []string) bool {
	return strings.Contains(record[0], "Account")
}
