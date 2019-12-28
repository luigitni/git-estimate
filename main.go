package main

import (
	"flag"
	"fmt"
	"git-estimate/estimate"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"sort"
	"time"
)

const (
	estSession = "session"
	estDay     = "day"
)

const timeArgsLayout = "2006-01-02T15-04"

func main() {
	path := flag.String("repo", ".", "git repository path. If no flag is specified the current folder is assumed")

	estMethod := flag.String("estimate", estSession,
		fmt.Sprintf("estimation method. Accepted values are %q and %q.", estSession, estDay))

	json := flag.Bool("json", false, "if true will output estimates in JSON format")

	baseline := flag.Float64("baseline", 2.0, "baseline value for session estimate")

	from := flag.String("from", "", "if provided computation starts from the given date and time. Format is yyyy-mm-ddThh-ii")

	flag.Parse()

	repo, err := git.PlainOpen(*path)

	var est estimate.Estimate
	switch *estMethod {
	case estSession:
		est = estimate.WorkingSession{Baseline: *baseline}
	case estDay:
		est = estimate.DayEstimate{}
	default:
		fmt.Printf("invalid estimation method. Accepted values are %q and %q", estSession, estDay)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("error opening repository at %s: %s", *path, err.Error())
		os.Exit(1)
	}

	// get the commit history
	var start time.Time
	if *from != "" {
		// parse the string
		start, err = time.Parse(timeArgsLayout, *from)
		if err != nil {
			fmt.Printf("unable to parse 'from' %s given. %s", *from, err.Error())
			os.Exit(1)
		}
	}

	iter, err := repo.Log(&git.LogOptions{All: true, Order: git.LogOrderCommitterTime})
	if err != nil {
		fmt.Printf("error reading log of repository: %s", err.Error())
		os.Exit(1)
	}

	defer iter.Close()

	byAuthors := make(map[string][]time.Time)

	// first group commits by authors, then for each author count the working days
	if err := iter.ForEach(func(commit *object.Commit) error {

		when := commit.Author.When
		if !start.IsZero() && when.Before(start) {
			return nil
		}

		sl, ok := byAuthors[commit.Author.Email]
		if !ok {
			sl = make([]time.Time, 0)
		}

		sl = append(sl, when)
		byAuthors[commit.Author.Email] = sl

		return nil
	}); err != nil {
		fmt.Printf("error reading commits: %s", err.Error())
		os.Exit(1)
	}

	// sort each slice of commits by date
	for k, _ := range byAuthors {
		commits := byAuthors[k]
		sort.Slice(commits, func(i, j int) bool {
			return commits[i].After(commits[j])
		})
	}

	res := est.Estimate(byAuthors)

	var formatter estimate.Formatter
	if *json {
		formatter = estimate.JSONFormatter{}
	} else {
		formatter = estimate.StringFormatter{}
	}

	fmt.Print(formatter.String(res))
	os.Exit(0)
}
