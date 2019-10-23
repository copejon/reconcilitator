package entry

import (
	"github.com/google/uuid"
	"time"
)

type Entry interface {
	Date() time.Time
	Payee() string
	Amount() float64
	Clear()
	IsCleared() bool
	ID() uuid.UUID
}
