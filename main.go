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
	bankCSVFlag = "bank"
)

func loadRegisterFromFile(r register.Register, fp string) error {
	f, err := os.OpenFile(fp, os.O_RDONLY, 0)
	if err != nil {
		panic(fmt.Errorf("could not open file: %v", err))
	}
	defer f.Close()
	err = r.Load(f)
	if err != nil {
		return fmt.Errorf("error loading file %s: %v", fp, err)
	}
	return nil
}

func formatDate(t time.Time) string {
	return fmt.Sprintf("%d/%d/%d", t.Month(), t.Day(), t.Year())
}

func init() {
	pflag.DurationVar(&afterDate, "since", 0, "The date after which to compare entries")
	pflag.StringVarP(&ynabToken, tokenFlag, "t", "", "Your YNAB token.")
	pflag.StringVarP(&ynabAccount, accountFlag, "a", "", "Case Sensitive. Because YNAB exports all accounts as one .csv_utils, we need this to target the relative entries")
	pflag.StringVarP(&ynabBudget, budgetFlag, "b", "", "Your YNAB Budget")
	pflag.StringVarP(&ynabCSV, ynabCSVFlag, "y", "", "Path to YNAB CSV file")
	pflag.StringVarP(&bankCSV, bankCSVFlag, "c", "", "Path to bank CSV file")
	pflag.Parse()
}

func main() {

	ynabReg := register.NewRegister(yxltr.NewTranslator(ynabAccount))
	err := loadRegisterFromFile(ynabReg, ynabCSV)
	if err != nil {
		panic(err)
	}

	usaaReg := register.NewRegister(uxltr.NewTranslator())
	err = loadRegisterFromFile(usaaReg, bankCSV)
	if err != nil {
		panic(err)
	}

	err = ynabReg.Clear(usaaReg)
	if err != nil {
		panic(err)
	}
	tbl := table.NewWriter()
	tbl.SetAutoIndex(true)
	tbl.SetTitle("YNAB Entries Test")
	tbl.AppendHeader(table.Row{"Date", "Payee", "Amount", "Cleared", "Sort"})

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
