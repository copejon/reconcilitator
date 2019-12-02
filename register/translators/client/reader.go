package client

import (
	"go.bmvs.io/ynab/api/transaction"
	"io"
	"strconv"
	"strings"
)

type transactionReader struct {
	t      []*transaction.Transaction
	r, eof int
}

func newTransactionReader(t []*transaction.Transaction) *transactionReader {
	return &transactionReader{
		t:   t,
		r:   0,
		eof: len(t),
	}
}

const Delimiter = ";"

func (c *transactionReader) Read(p []byte) (n int, err error) {
	if c.r > c.eof {
		return 0, io.EOF
	}
	out, n := c.read()
	p = append(p, []byte(out)...)
	return len(out), nil
}

func (c *transactionReader) read() (string, int) {
	cur := c.t[c.r]
	c.r++
	date := cur.Date.String()
	amt := strconv.FormatInt(cur.Amount, 10)
	payee := *cur.PayeeName

	out := strings.Join([]string{date, amt, payee}, Delimiter)
	return out, len(out)
}
