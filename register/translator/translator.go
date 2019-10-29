package translator

import "main/register/entry"

type Translator interface {
	ToEntry(s []string) (*entry.Entry, error)
}
