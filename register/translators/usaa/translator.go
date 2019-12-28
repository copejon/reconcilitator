package usaa

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
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

// Translate takes a reader who's output is expected to be CSV formatted
// returns a slice of entries or an error
func (t *translator) Translate(r io.Reader) (entries []*entry.Entry, err error) {
	nr, err := normalizeLineFeed(r)
	if err != nil {
		return nil, err
	}
	rdr := csv.NewReader(nr)
	rdr.LazyQuotes = true
	entries, err = t.readRecordsToEntries(rdr)
	if err != nil {
		return nil, fmt.Errorf("failed translation: %v", err)
	}
	return entries, nil
}

// TODO this entire func is a dupe of ynab's.
//  should write a generic function that executes an anon func
func (t *translator) readRecordsToEntries(r *csv.Reader) ([]*entry.Entry, error) {
	entries := make([]*entry.Entry, 0, 4096)
	var e *entry.Entry
	var err error
	for {
		e, err = t.readRecordToEntry(r)
		if err != nil || e == nil {
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
	fmt.Printf("rec: <%s> len: %d\n", record, len(record))
	if err != nil {
		return nil, err
	}
	if len(record) == 0 {
		return nil, nil
	}
	record[amount] = strings.TrimSpace(record[amount])
	return t.parseRecord(record)
}

func (t *translator) parseRecord(rec []string) (*entry.Entry, error) {

	if len(rec) == 0 {
		return nil, nil
	}
	var err error
	errMsg := `couldn't generate entry: %v\n'`
	e := entry.NewEntry()

	e.SetPayee(rec[payee])

	const pattern = "01/02/2006"

	d, err := register.ParseDate(rec[date], pattern)
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

const (
	crByte = 0x0d
	lfByte = 0x0a
)

func normalizeLineFeed(r io.Reader) (*bytes.Reader, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error normalizing linefeeds: %v", err)
	}
	buf = bytes.ReplaceAll(buf, []byte{crByte}, []byte{lfByte})
	return bytes.NewReader(buf), nil
}
