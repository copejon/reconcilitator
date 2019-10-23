package ynab

import (
	"encoding/csv"
	"fmt"
	"io"
	"main/register"
	iface "main/register/entry"
)

var _ register.Register = &ynabRegister{}

func NewYNABRegister() *ynabRegister {
	return &ynabRegister{
		entries: make([]iface.Entry, 0, 64),
	}
}

type ynabRegister struct {
	entries []iface.Entry
	readCursor int
}

func (y *ynabRegister) Load(r io.Reader) error {
	rdr := csv.NewReader(r)
	rdr.LazyQuotes = true

	//strip header row
	_, err := rdr.Read()
	if err != nil {
		return err
	}

	for {
		raw, err := rdr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("failed to load register: %v\n", err)
		}
		ent, err := NewEntry(raw)
		if err != nil {
			return fmt.Errorf("failed to load new entry: %v\n", ent)
		}
		if ent == nil {
			continue // Not the YNAB account we are examining
		}
		y.entries = append(y.entries, ent)
	}
	return nil
}

func (y *ynabRegister) Read() (iface.Entry, error) {
	if y.readCursor >= len(y.entries){
		return nil, io.EOF
	}
	e := y.entries[y.readCursor]
	y.readCursor++
	return e, nil
}

func (y *ynabRegister) ReadAll() ([]iface.Entry, error) {
	if y.readCursor >= len(y.entries){
		return nil, io.EOF
	}
	remainder := y.entries[y.readCursor:]
	y.readCursor = len(y.entries)
	return remainder, nil
}
