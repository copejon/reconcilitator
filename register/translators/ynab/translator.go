package ynab

import (
	"encoding/csv"
	"fmt"
	"io"
	"main/register"
	"main/register/entry"
	"main/register/translator"
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
	return &ynabTranslator{account: account}
}

func (t *ynabTranslator) Translate(r io.Reader) ([]*entry.Entry, error) {
	rdr := csv.NewReader(r)
	rdr.LazyQuotes = true
	return t.readRecordsToEntries(rdr)
}

func (t *ynabTranslator) readRecordsToEntries(r *csv.Reader) (entries []*entry.Entry, err error) {
	entries = make([]*entry.Entry, 0, 4096)
	var e *entry.Entry
	for stop := false; !stop; {
		e, err = t.readRecordToEntry(r)
		if err == nil {
			entries = append(entries, e)
		}
		stop = stopReads(err)
	}
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("could not read records to entries: %v", err)
	}
	return entries, nil
}

func stopReads(err error) bool {
	return err != nil && !isExpected(err)
}

func (t *ynabTranslator) readRecordToEntry(r *csv.Reader) (*entry.Entry, error) {
	rec, err := r.Read()
	if err != nil {
		return nil, err
	}
	return t.parseRecord(rec)
}

func (t *ynabTranslator) parseRecord(record []string) (*entry.Entry, error) {

	// Ynab CSVs contain all tracked accounts as one, ignore all but the specified one

	if record[acct] != t.account {
		return nil, newIsNotAccountError(t.account, record[acct])
	}
	errMsg := `couldn't generate entry: %v\n`
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
		// YNAB stores outflows at positive numbers, so we negate them here
		if amt > 0 {
			amt = 0 - amt
		}
	}
	e.SetAmount(amt)
	return e, nil

}
