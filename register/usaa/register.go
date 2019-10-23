package usaa

import (
	"encoding/csv"
	"fmt"
	"io"
	"main/register"
	iface "main/register/entry"
)

var _ register.Register = &usaaRegister{}

func NewUSAARegister() *usaaRegister {
	return &usaaRegister{
		entries: make([]iface.Entry, 0, 64),
	}
}

type usaaRegister struct {
	entries []iface.Entry
	readCursor int
}

func (u *usaaRegister) Load(r io.Reader) error {
	rdr := csv.NewReader(r)
	rdr.LazyQuotes = true

	for {
		raw, err := rdr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("failed to load register: %v\n", err)
		}
		ent, err := NewEntry(raw)
		if err != nil {
			return fmt.Errorf("failed to get new entry: %v\n", err)
		}
		u.entries = append(u.entries, ent)
	}
	return nil
}

func (u *usaaRegister) Read() (iface.Entry, error) {
	if u.readCursor >= len(u.entries){
		return nil, io.EOF
	}
	e := u.entries[u.readCursor]
	u.readCursor++
	return e, nil
}

func (u *usaaRegister) ReadAll() ([]iface.Entry, error) {
	if u.readCursor >= len(u.entries){
		return nil, io.EOF
	}
	remainder := u.entries[u.readCursor:]
	u.readCursor = len(u.entries)
	return remainder, nil
}

