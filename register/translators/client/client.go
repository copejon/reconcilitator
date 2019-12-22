package client

import (
	"go.bmvs.io/ynab"
	"go.bmvs.io/ynab/api"
	"go.bmvs.io/ynab/api/transaction"
	"time"
)

type Config struct {
	budgetID,
	token string
	since *time.Time
}

type Reader struct {
	cfg *Config
	*transactionReader
}

func NewClientReader(cfg *Config) (*Reader, error) {
	c := ynab.NewClient(cfg.token)
	date, err := api.DateFromString(cfg.since.String())
	if err != nil {
		return nil, err
	}
	f := &transaction.Filter{
		Since: &date,
		Type:  nil,
	}
	trans, err := c.Transaction().GetTransactions(cfg.budgetID, f)
	rdr := newTransactionReader(trans)
	return &Reader{
		cfg:               cfg,
		transactionReader: rdr}, nil
}
