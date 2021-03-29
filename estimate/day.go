package estimate

import (
	"fmt"
	"strings"
	"time"
)

type DayEstimate struct{}

func (d DayEstimate) Estimate(byAuthors map[string][]time.Time) []Result {
	results := make([]Result, len(byAuthors))
	c := 0
	for k := range byAuthors {
		r := &results[c]
		r.Author = k
		if strings.Count(k, "@") > 1 {
			p := strings.Split(k, "@")
			r.Author = fmt.Sprintf("%s@%s", p[0], p[1])
			r.Issue = p[2]
		}
		prev := time.Time{}
		v := byAuthors[k]
		for _, t := range v {
			if prev.IsZero() {
				prev = t
				continue
			}
			if prev.YearDay() != t.YearDay() {
				r.Days += 1.0
			}
			prev = t
		}
		r.Days += 1.0
		r.Hours = r.Days * 8.0
		c++
	}
	return results
}
