package main

import (
	"fmt"
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/spf13/pflag"
	"main/register"
	"main/register/dateMapper"
	uxltr "main/register/translators/usaa"
	yxltr "main/register/translators/ynab"
	"os"
	"time"
)

var (
	ynabToken,
	ynabBudget,
	ynabAccount,
	ynabCSV,
	bankCSV string
	afterDate time.Duration
)

const (
	tokenFlag   = "token"
	accountFlag = "account"
	budgetFlag  = "budget"
	ynabCSVFlag = "ynab"
	usaaCSVFlag = "bank"
)

func loadRegisterFromFile(r register.Register, fp string) {
	f, err := os.OpenFile(fp, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = r.Load(f)
	if err != nil {
		panic(err)
	}
}

func formatDate(t time.Time) string {
	return fmt.Sprintf("%d/%d/%d", t.Month(), t.Day(), t.Year())
}

func init() {
	pflag.DurationVar(&afterDate, "since", 0, "The date after which to compare entries")
	pflag.StringVarP(&ynabToken, tokenFlag, "t", "", "Your YNAB token.")
	pflag.StringVarP(&ynabAccount, accountFlag, "a", "", "Case Sensitive. Because YNAB exports all accounts as one .csv, we need this to target the relative entries")
	pflag.StringVarP(&ynabBudget, budgetFlag, "b", "", "Your YNAB Budget")
	pflag.StringVarP(&ynabCSV, ynabCSVFlag, "y", "", "Path to YNAB CSV file")
	pflag.StringVarP(&bankCSV, usaaCSVFlag, "f", "", "Path to bank CSV file")
	pflag.Parse()
}

func main() {

	ynabReg := register.NewRegister(yxltr.NewTranslator(ynabAccount))
	loadRegisterFromFile(ynabReg, ynabCSV)

	usaaReg := register.NewRegister(uxltr.NewTranslator())
	loadRegisterFromFile(usaaReg, bankCSV)

	cleared := ynabReg.Clear(usaaReg)
	tbl := table.NewWriter()
	tbl.SetAutoIndex(true)
	tbl.SetTitle("YNAB Entries Test")
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

	startDate := dateMapper.MostRecentStartTime(ynabReg.Entries(), usaaReg.Entries())

	fmt.Printf("back date: %v\n", startDate)
	for _, day := range ynabReg.Entries() {
		for _, e := range day {
			if e.Date().Before(startDate) {
				continue
			}
			tbl.AppendRow(table.Row{formatDate(e.Date()), e.Payee(), e.Amount(), e.Cleared(), e.Date().Unix()})
		}
	}
	fmt.Println(tbl.Render())
}
