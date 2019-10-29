package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/pflag"
	"main/dateMapper"
	"main/register"
	uxltr "main/register/translators/usaa"
	yxltr "main/register/translators/ynab"
	"os"
	"time"
)

var (
	ynabAccount,
	ynabCSV,
	bankCSV string
)

func init() {
	pflag.StringVarP(&ynabAccount, "account", "a", "", "Case Sensitive. Because YNAB exports all accounts as one .csv, we need this to target the relative entries")
	pflag.StringVarP(&ynabCSV, "ynab", "y", "", "Path to YNAB CSV file")
	pflag.StringVarP(&bankCSV, "bank", "b", "", "Path to bank CSV file")
	pflag.Parse()
}

func main() {

	if ynabCSV == "" {
		panic(fmt.Errorf("missing required option: -ynab|-y"))
	}
	if ynabAccount == "" {
		panic(fmt.Errorf("missing required option: -account|-a"))
	}
	if bankCSV == "" {
		panic(fmt.Errorf("missing required option: -bank|-b"))
	}

	ynabReg := register.NewRegister(yxltr.NewTranslator(ynabAccount))
	loadRegister(ynabReg, ynabCSV)

	usaaReg := register.NewRegister(uxltr.NewTranslator())
	loadRegister(usaaReg, bankCSV)

	ynab, err := dateMapper.NewDateMapper(ynabReg)
	if err != nil {
		panic(err)
	}
	usaa, err := dateMapper.NewDateMapper(usaaReg)
	if err != nil {
		panic(err)
	}

	cleared := ynab.ClearEntries(usaa)
	tbl := table.NewWriter()
	tbl.SetAutoIndex(true)
	tbl.SetTitle("YNAB Entries vs USAA")
	tbl.AppendHeader(table.Row{"Date", "Payee", "Amount", "Cleared", "Sort"})
	tbl.AppendFooter(table.Row{fmt.Sprintf("Cleared: %d", cleared)})

	tbl.SortBy([]table.SortBy{{
		Name: "Sort",
		Mode: table.DscNumeric,
	}})

	tbl.SetRowPainter(func(row table.Row) text.Colors {
		if cleared, ok := row[3].(bool); ok && !cleared {
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

	backDate := dateMapper.MostRecentStartTime(ynab, usaa)

	fmt.Printf("back date: %v\n", backDate)
	for _, entries := range ynab {
		for _, e := range entries {
			if e.Date().Before(backDate) {
				break
			}
			tbl.AppendRow(table.Row{formatDate(e.Date()), e.Payee(), e.Amount(), e.Cleared(), e.Date().Unix()})
		}
	}
	fmt.Println(tbl.Render())
}

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
