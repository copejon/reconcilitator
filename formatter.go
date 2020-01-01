package main

import (
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"main/register"
	"time"
)

const (
	date = iota + 1 // row count begins at 1
	_               // payee
	amount
	cleared
	sort
)

func columnHeaders() table.Row {
	return table.Row{"Date", "Payee", "Amount", "Cleared", "Sort"}
}

const layout = "01/02/2006"

const currencyFormat = "%.2f"

func columnConfig() []table.ColumnConfig {
	return []table.ColumnConfig{
		{
			Number:      date,
			Transformer: text.NewTimeTransformer(layout, nil),
		},
		{
			Number:      amount,
			Transformer: text.NewNumberTransformer(currencyFormat),
		},
	}
}

func sortConfig() []table.SortBy {
	return []table.SortBy{{
		Number: sort,
		Mode:   table.DscNumeric,
	}}
}

func paintUncleared(row table.Row) text.Colors {
	if cleared, ok := row[cleared-1].(bool); ok && !cleared {
		return text.Colors{text.BgRed, text.Bold}
	}
	return nil
}

type tableFormatter struct {
	w     table.Writer
	title string
	opts  *tableFormatterOpts
	reg   register.Register
}

type tableFormatterOpts struct {
	since     time.Time
	numDays   int
	autoIndex bool
}

func defaultOpts() *tableFormatterOpts {
	return &tableFormatterOpts{
		since:     time.Time{},
		numDays:   -1,
		autoIndex: true,
	}
}

func (t *tableFormatter) setWriter() (tw table.Writer) {
	t.w = table.NewWriter()
	t.w.SetTitle(t.title)
	t.w.SortBy(sortConfig())
	t.w.SetAutoIndex(t.opts.autoIndex)
	t.w.SetRowPainter(paintUncleared)
	t.w.AppendHeader(columnHeaders())
	t.w.SetColumnConfigs(columnConfig())
	return
}

func (t *tableFormatter) loadTable() {
	opts := t.opts
	for _, day := range t.reg.Entries() {
		for _, e := range day {
			// if the entry falls on or after opts.since, append it to the table
			if e.Date().After(opts.since) ||
				e.Date().Equal(opts.since) {
				t.w.AppendRow(table.Row{
					e.Date(),
					e.Payee(),
					e.Amount(),
					e.Cleared(),
					e.Date().Unix(),
				})
			}
		}
	}
}

func NewTableFormatter(r register.Register, title string, opts *tableFormatterOpts) *tableFormatter {
	tf := &tableFormatter{
		title: title,
		opts:  opts,
		reg:   r,
	}
	tf.w = tf.setWriter()
	tf.loadTable()
	return tf
}

func NewDefaultTableFormatter(r register.Register, title string) *tableFormatter {
	return NewTableFormatter(r, title, defaultOpts())
}
