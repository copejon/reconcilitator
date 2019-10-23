package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"main/datehash"
	"main/register"
	"main/register/usaa"
	"main/register/ynab"
	"os"
	"time"
)

func loadRegister(r register.Register, fp string) {
	f, err := os.OpenFile(fp, os.O_RDONLY, 0)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	err = r.Load(f)
	if err != nil {
		panic(err)
	}
}

func formatDate(t time.Time) string {
	return fmt.Sprintf("%d/%d/%d", t.Month(), t.Day(), t.Year())
}

func main() {

	args := os.Args[1:]

	ynabReg := ynab.NewYNABRegister()
	usaaReg := usaa.NewUSAARegister()

	loadRegister(ynabReg, args[0])
	loadRegister(usaaReg, args[1])

	ynabHash, err := datehash.NewDateHashMap(ynabReg)
	if err != nil {
		panic(err)
	}
	usaaHash, err := datehash.NewDateHashMap(usaaReg)
	if err != nil {
		panic(err)
	}

	cleared := ynabHash.ClearHashedEntries(usaaHash)
	tbl := table.NewWriter()
	tbl.SetAutoIndex(true)
	tbl.SetTitle("YNAB Entries vs USAA")
	tbl.AppendHeader(table.Row{"Date", "Payee", "Amount", "Cleared", "Sort"})
	tbl.AppendFooter(table.Row{"Cleared: ", cleared})

	tbl.SortBy([]table.SortBy{{
		Name: "Sort",
		Mode: table.DscNumeric,
	}})

	tbl.SetRowPainter(func(row table.Row) text.Colors {
		if cleared, ok := row[3].(bool); ok && ! cleared {
			return text.Colors{text.BgRed, text.Bold}
		}
		return nil
	})
	tbl.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:        "Amount",
			Transformer: text.NewNumberTransformer("%.2f"),
		},
	})

	for _, entries := range ynabHash {
		for _, e := range entries {
			tbl.AppendRow(table.Row{formatDate(e.Date()), e.Payee(), e.Amount(), e.IsCleared(), e.Date().Unix()})
		}
	}
	fmt.Println(tbl.Render())
}
