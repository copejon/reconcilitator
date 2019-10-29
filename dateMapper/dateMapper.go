package dateMapper

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"main/register"
	"main/register/entry"
	"time"
)

func NewDateMapper(r register.Register) (DateMapper, error) {
	var m = make(DateMapper)
	for {
		e, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("unable to map entry: %v\n", err)
		}
		m[e.Date()] = append(m[e.Date()], e)
	}
	return m, nil
}

func clearEntries(lhs, rhs []*entry.Entry) (cleared int) {
	for _, lhsEnt := range lhs {
		if lhsEnt.Cleared() { // Safety catch, ignore cleared LHS entries
			break
		}
		for _, rhsEnt := range rhs {
			if !rhsEnt.Cleared() && lhsEnt.Amount() == rhsEnt.Amount() {
				rhsEnt.SetCleared()
				lhsEnt.SetCleared()
				cleared++
				break
			}
		}
	}
	return
}

func oldestEntry(d DateMapper) time.Time {
	oldest := time.Now()
	for t := range d {
		if t.Before(oldest) {
			oldest = t
		}
	}
	fmt.Printf("oldest date: %v\n", oldest)
	return oldest
}

func MostRecentStartTime(d1, d2 DateMapper) (t time.Time) {
	t = oldestEntry(d1)
	if altT := oldestEntry(d2); altT.After(t) {
		t = altT
	}
	return
}

type DateMapper map[time.Time][]*entry.Entry

func (d DateMapper) ClearEntries(rhs DateMapper) (cleared int) {
	lhs := d

	startDate := MostRecentStartTime(lhs, rhs)

	for date, lhsDaysEntries := range lhs {
		// Only compare entries from dates within the same span, ignore dates prior to this subset
		if date.Before(startDate) {
			continue
		}
		if rhsDaysEntries, ok := rhs[date]; ok {
			cleared += clearEntries(lhsDaysEntries, rhsDaysEntries)
		}
	}
	return
}

func (d DateMapper) DayHasEntry(t time.Time, amount float64) (bool, int) {
	day := d[t]
	var instances int
	for _, e := range day {
		if e.Amount() == amount {
			instances++
		}
	}
	return instances != 0, instances
}

func (d DateMapper) ClearEntry(t time.Time, uuid uuid.UUID) bool {
	day := d[t]
	var cleared bool
	for _, e := range day {
		if e.Id() == uuid {
			e.SetCleared()
			cleared = true
			break
		}
	}
	return cleared
}

func (d DateMapper) Entries() (i int) {
	for _, e := range d {
		i += len(e)
	}
	return
}
