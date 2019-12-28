package estimate

import (
	"time"
)

type DayEstimate struct {}

func (d DayEstimate) Estimate(byAuthors map[string][]time.Time) []Result {
	results := make([]Result, len(byAuthors))
	c := 0
	for k, _ := range byAuthors {
		r := &results[c]
		r.Author = k
		prev := time.Time{}
		v := byAuthors[k]
		for _, t := range v {
			if prev.IsZero() {
				prev = t
				r.Days += 1.0
				continue
			}
			if prev.YearDay() != t.YearDay() {
				r.Days += 1.0
			}
			prev = t
		}
		r.Hours = r.Days * 8.0
		c++
	}
	return results
}
