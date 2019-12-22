package translator_errors

type ExpectedError interface {
	IsExpected(e error) bool
}
