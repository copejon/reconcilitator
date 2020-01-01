package main

import (
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"main/register"
	"time"
)

type tableFormatter struct {
	w   table.Writer
	reg register.Register
}

type tableFormatterOpts struct {
	since     time.Time
	numDays   int
	autoIndex bool
}

const (
	date = iota + 1 // row count begins at 1
	_               // payee
	amount
	cleared
	sort
)

var columnHeaders = table.Row{"Date", "Payee", "Amount", "Cleared", "Sort"}

const layout = "01/02/2006"

var columnConfig = []table.ColumnConfig{
	{
		Number:      date,
		Transformer: text.NewTimeTransformer(layout, nil),
	},
	{
		Number:      amount,
		Transformer: text.NewNumberTransformer("%.2f"),
	},
}

func paintUncleared(row table.Row) text.Colors {
	if cleared, ok := row[cleared-1].(bool); ok && !cleared {
		return text.Colors{text.BgRed, text.Bold}
	}
	return nil
}

func NewTableFormatter(r register.Register, title string, opts *tableFormatterOpts) *tableFormatter {
	tw := table.NewWriter()
	tw.SetAutoIndex(opts.autoIndex)
	tw.SetTitle(title)
	tw.AppendHeader(columnHeaders)

	tw.SortBy([]table.SortBy{{
		Number: sort,
		Mode:   table.DscNumeric,
	}})
	tw.SetRowPainter(paintUncleared)
	tw.SetColumnConfigs(columnConfig)

	for _, day := range r.Entries() {
		for _, e := range day {
			// if the entry falls on or after opts.since, append it to the table
			if e.Date().After(opts.since) ||
				e.Date().Equal(opts.since) {
				tw.AppendRow(table.Row{
					e.Date(),
					e.Payee(),
					e.Amount(),
					e.Cleared(),
					e.Date().Unix(),
				})
			}
		}
	}
	return &tableFormatter{
		w:   tw,
		reg: r,
	}
}

func defaultOpts() *tableFormatterOpts {
	return &tableFormatterOpts{
		since:     time.Time{},
		numDays:   -1,
		autoIndex: true,
	}
}

func NewDefaultTableFormatter(r register.Register, title string) *tableFormatter {
	return NewTableFormatter(r, title, defaultOpts())
}
