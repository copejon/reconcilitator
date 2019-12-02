package client

import (
	"fmt"
	"go.bmvs.io/ynab"
	"go.bmvs.io/ynab/api"
	"go.bmvs.io/ynab/api/transaction"
	"io"
	"main/register/entry"
	"time"
)

type translator struct {
}

type TranslatorConfig struct {
	Token,
	Budget,
	Account string
	Since time.Time
}

func NewTranslator(t *TranslatorConfig) (*translator, error) {
	c := ynab.NewClient(t.Token)
	if c == nil {
		return nil, fmt.Errorf("error getting ynab client")
	}
	f := &transaction.Filter{Since: &api.Date{
		Time: t.Since,
	}}
	trans, err := c.Transaction().GetTransactionsByAccount(t.Budget, t.Account, f)
	if err != nil {
		return nil, err
	}
	return &translator{c: c, transactions: trans}, nil
}

func (t *translator) Translate(r io.Reader) ([]*entry.Entry, error) {

}

func (t *translator) ToEntry(s []string) (*entry.Entry, error) {

}
