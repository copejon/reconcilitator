package client

import (
	"fmt"
	"go.bmvs.io/ynab"
	"go.bmvs.io/ynab/api"
	"go.bmvs.io/ynab/api/transaction"
	"main/register/entry"
	iface "main/register/translator"
	"time"
)

type translator struct {
	c            ynab.ClientServicer
	transactions []*transaction.Transaction
}

var _ iface.Translator = &translator{}

func NewTranslator(token, budget, account string, after time.Time) (*translator, error) {
	c := ynab.NewClient(token)
	if c == nil {
		return nil, fmt.Errorf("error getting ynab client")
	}
	f := &transaction.Filter{Since: &api.Date{
		Time: after,
	}}
	t, err := c.Transaction().GetTransactionsByAccount(budget, account, f)
	if err != nil {
		return nil, err
	}
	return &translator{c: c, transactions: t}, nil
}

func (t *translator) ToEntry(s []string) (*entry.Entry, error) {

}
