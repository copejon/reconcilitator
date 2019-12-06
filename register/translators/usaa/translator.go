package usaa

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
	_ = iota
	_
	date
	_
	payee
	_
	amount
)

var _ iface.Translator = &translator{}

type translator struct{}

func NewTranslator() *translator {
	return &translator{}
}

func (t *translator) Translate(r io.Reader) (entries []*entry.Entry, err error) {
	rdr := csv.NewReader(r)
	rdr.LazyQuotes = true
	entries, err = t.readRecordsToEntries(rdr)
	if err != nil {
		err = fmt.Errorf("failed translation")
	}
	return
}

// TODO this entire func is a dupe of ynab's.
//  should write a generic function that executes an anon func
func (t *translator) readRecordsToEntries(r *csv.Reader) ([]*entry.Entry, error) {
	entries := make([]*entry.Entry, 0, 4096)
	var e *entry.Entry
	var err error
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
	record, err := r.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read record: %v", err)
	}
	return t.parseRecord(record)
}

func (t *translator) parseRecord(rec []string) (*entry.Entry, error) {
	var err error
	errMsg := `couldn't generate entry: %v\n'`
	e := entry.NewEntry()

	e.SetPayee(rec[payee])
	d, err := register.ParseDate(rec[date])
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	e.SetDate(d)
	// USAA uses double negation to show positive values.  Makes perfect sense.
	a, err := register.ParseCurrency(strings.TrimPrefix(rec[amount], "--"))
	if err != nil {
		return nil, fmt.Errorf(errMsg, err)
	}
	e.SetAmount(a)
	return e, nil
}
