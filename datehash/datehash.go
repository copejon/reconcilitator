package datehash

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"main/register"
	"main/register/entry"
	"time"
)

func NewDateHashMap(r register.Register) (DateHash, error) {
	var m = make(DateHash)
	for {
		e, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("unable to hash entry: %v\n", err)
		}
		m[e.Date()] = append(m[e.Date()], e)
	}
	return m, nil
}

func clearEntries(lhs, rhs []entry.Entry) (cleared int) {
	for _, lhsEnt := range lhs {
		if lhsEnt.IsCleared() { // Safety catch, ignore cleared LHS entries
			break
		}
		for _, rhsEnt := range rhs {
			if ! rhsEnt.IsCleared() && lhsEnt.Amount() == rhsEnt.Amount() {
				rhsEnt.Clear()
				lhsEnt.Clear()
				cleared++
				break
			}
		}
	}
	return
}

func oldestEntry(d DateHash) time.Time {
	var oldest time.Time
	for t, _ := range d {
		if t.Before(oldest) {
			oldest = t
		}
	}
	return oldest
}

func MostRecentStartTime(d1, d2 DateHash) (t time.Time) {
	t = oldestEntry(d1)
	if altT := oldestEntry(d2); altT.After(t) {
		t = altT
	}
	return
}

type DateHash map[time.Time][]entry.Entry

func (d DateHash) ClearHashedEntries(rhs DateHash) (cleared int) {
	lhs := d

	startDate := MostRecentStartTime(lhs, rhs)

	for day, lhsDay := range lhs {
		// Only compare entries from dates within the same span, ignore dates prior to this subset
		if day.Before(startDate) {
			continue
		}
		if rhsEntries, ok := rhs[day]; ok {
			cleared += clearEntries(lhsDay, rhsEntries)
		}
	}
	return
}

func (d DateHash) DayHasEntry(t time.Time, amount float64) (bool, int) {
	day := d[t]
	var instances int
	for _, e := range day {
		if e.Amount() == amount {
			instances++
		}
	}
	return instances != 0, instances
}

func (d DateHash) ClearEntry(t time.Time, uuid uuid.UUID) bool {
	day := d[t]
	var cleared bool
	for _, e := range day {
		if e.ID() == uuid {
			e.Clear()
			cleared = true
			break
		}
	}
	return cleared
}
