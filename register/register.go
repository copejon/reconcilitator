package register

import (
	"fmt"
	"io"
	"main/register/dateMapper"
	"main/register/translator"
	"strconv"
	"strings"
	"time"
)

type Register interface {
	Load(io.Reader) error
	Clear(Register) error
	Entries() dateMapper.DateMapper
}

var _ Register = &register{}

// register implements the Register AND Translator interfaces
type register struct {
	translator.Translator
	entries    dateMapper.DateMapper
	readCursor int
}

func (r *register) Entries() dateMapper.DateMapper {
	return r.entries
}

func NewRegister(translator translator.Translator) *register {
	return &register{
		entries:    make(dateMapper.DateMapper),
		readCursor: 0,
		Translator: translator,
	}
}

func (r *register) Load(rdr io.Reader) error {
	entries, err := r.Translate(rdr)
	if err != nil {
		return fmt.Errorf("error loading from file: %v", err)
	}
	for _, e := range entries {
		if e == nil {
			// translators are allowed to return nil entries
			continue
		}
		r.entries.Push(e)
	}
	return nil
}

func (r *register) Clear(reg Register) error {
	r.entries.ClearEntries(reg.Entries())
	return nil
}

func ParseCurrency(c string) (float64, error) {
	s := strings.Trim(c, "$")
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		err = fmt.Errorf("error parsing currency: %v\n", err)
	}
	return f, err
}

func ParseDate(d, pattern string) (time.Time, error) {
	t, err := time.Parse(pattern, d)
	if err != nil {
		err = fmt.Errorf("error parsing entry date: %v\n", err)
	}
	t = t.Truncate(24 * time.Hour)
	return t, err
}
