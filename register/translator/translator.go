package translator

import (
	"io"
	"main/register/entry"
)

type Translator interface {
	Translate(r io.Reader) ([]*entry.Entry, error)
}
