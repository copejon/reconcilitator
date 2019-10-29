package register

import (
	"encoding/csv"
	"fmt"
	"io"
	"main/register/entry"
	"main/register/translator"
	"strconv"
	"strings"
	"time"
)

type Register interface {
	Load(io.Reader) error
	Read() (*entry.Entry, error)
	ReadAll() ([]*entry.Entry, error)
}

type register struct {
	entries    []*entry.Entry
	xltr       translator.Translator
	readCursor int
}

var _ Register = &register{}

func (r *register) Load(rdr io.Reader) error {
	csvRdr := csv.NewReader(rdr)
	csvRdr.LazyQuotes = true

	for {
		line, err := csvRdr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("failed to load register: %v\n", err)
		}

		e, err := r.xltr.ToEntry(line)
		if err != nil {
			return fmt.Errorf("failed to get new entry: %v\n", err)
		}
		if e == nil {
			// translators are allowed to return nil entries as a result of filtering, etc.
			continue
		}
		r.entries = append(r.entries, e)
	}
	return nil
}

func (r *register) Read() (*entry.Entry, error) {
	if r.readCursor >= len(r.entries) {
		return nil, io.EOF
	}
	e := r.entries[r.readCursor]
	r.readCursor++
	return e, nil
}

func (r *register) ReadAll() ([]*entry.Entry, error) {
	if r.readCursor >= len(r.entries) {
		return nil, io.EOF
	}
	remainder := r.entries[r.readCursor:]
	r.readCursor = len(r.entries)
	return remainder, nil
}

func NewRegister(t translator.Translator) *register {
	return &register{
		entries:    make([]*entry.Entry, 0),
		readCursor: 0,
		xltr:       t,
	}
}

//const bomBytes = "\xEF\xBB\xBF"
//
//func stripUTF8BOM(s string) string {
//	b := bytes.TrimLeft([]byte(s), bomBytes)
//	return string(b)
//}

func ParseCurrency(c string) (float64, error) {
	s := strings.Trim(c, "$")
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		err = fmt.Errorf("error parsing currency: %v\n", err)
	}
	return f, err
}

func ParseDate(d string) (time.Time, error) {
	const format = `01/02/2006`
	t, err := time.Parse(format, d)
	if err != nil {
		err = fmt.Errorf("error parsing entry date: %v\n", err)
	}
	t = t.Truncate(24 * time.Hour)
	return t, err
}
