package dateMapper

import (
	"main/register/entry"
	"time"
)

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

func (d DateMapper) Push(entry *entry.Entry) {
	d[entry.Date()] = append(d[entry.Date()], entry)
}

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

func (d DateMapper) ClearEntry(againstEntry *entry.Entry) bool {
	daysEntries := d[againstEntry.Date()]
	var cleared bool
	for _, e := range daysEntries {
		if e.ImportID() == againstEntry.ImportID() {
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
