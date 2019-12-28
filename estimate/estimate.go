package estimate

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Result struct {
	Author string  `json:"author"`
	Hours  float64 `json:"hours"`
	Days   float64 `json:"days"`
}

type Estimate interface {
	Estimate(map[string][]time.Time) []Result
}

type Formatter interface {
	String(res []Result) string
}

// Format the result as a JSON object
type JSONFormatter struct{}

func (f JSONFormatter) String(results []Result) string {
	total := Result{Author: "all"}
	for _, res := range results {
		total.Hours += res.Hours
		total.Days += res.Days
	}

	res := struct {
		Developers []Result `json:"developers"`
		Overall    Result   `json:"overall"`
	}{
		results,
		total,
	}

	b, err := json.Marshal(res)
	if err != nil {
		fmt.Printf("error marshaling result: %s", err.Error())
		os.Exit(1)
	}
	return string(b)
}

type StringFormatter struct{}

func (f StringFormatter) String(results []Result) string {
	total := Result{Author: "all"}
	var builder strings.Builder
	for _, res := range results {
		builder.WriteString(fmt.Sprintf("commits by %s", res.Author))
		builder.WriteString(fmt.Sprintf("\n=== %.2f days (%.2f hours)", res.Days, res.Hours))
		builder.WriteString("\n\n")
		total.Hours += res.Hours
		total.Days += res.Days
	}
	builder.WriteString(fmt.Sprintf("overall %.2f days (%.2f hours)", total.Days, total.Hours))

	return builder.String()
}
