package register

import (
	"bytes"
	"fmt"
	"io"
	"main/register/entry"
	"strconv"
	"strings"
	"time"
)

type Register interface {
	Load(io.Reader) error
	Read() (entry.Entry, error)
	ReadAll() ([]entry.Entry, error)
}

type sourceId int

const (
	srcUSAA sourceId = iota
	srcYNAB
)

func GuessSource(s string) (sourceId, error) {
	const (
		ynab = "\"Account\""
		usaa = "posted"
	)

	s = stripUTF8BOM(s)
	switch {
	case s == usaa:
		return srcUSAA, nil
	case s == ynab:
		return srcYNAB, nil
	default:
		return -1, fmt.Errorf("unrecognized column header: %v\n", s)
	}
}

const bomBytes = "\xEF\xBB\xBF"

func stripUTF8BOM(s string) string {
	b := bytes.TrimLeft([]byte(s), bomBytes)
	return string(b)
}

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
	t = t.Truncate(24*time.Hour)
	return t, err
}
