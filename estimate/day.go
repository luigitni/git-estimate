package estimate

import (
	"time"
)

const LogFmt = "2 Jan 2006 15:04:05"

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
				r.Hours += 1.0
				continue
			}
			if prev.YearDay() != t.YearDay() {
				r.Hours += 1.0
			}
			prev = t
		}
		r.Days = r.Hours / 8.0
		c++
	}
	return results
}
