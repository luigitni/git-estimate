package estimate

import (
	"fmt"
	"strings"
	"time"
)

type WorkingSession struct {
	Baseline float64
}

// sums up the hours in a given day, assuming the beginning of the session 2 hours earlier than the first commit
// and the end of the job in the last commit of the day
func (ws WorkingSession) Estimate(byAuthors map[string][]time.Time) []Result {
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
		next := time.Time{}
		v := byAuthors[k]
		for _, t := range v {
			if next.IsZero() {
				next = t
				continue
			}

			diff := next.Sub(t).Hours()
			if diff < 8 {
				r.Hours += diff
			} else {
				r.Hours += (time.Duration(ws.Baseline) * time.Hour).Hours()
			}
			next = t
		}

		// add the last/first padding to the commit
		r.Hours += (time.Duration(ws.Baseline) * time.Hour).Hours()

		r.Days = r.Hours / 8.0
		c++
	}
	return results
}
