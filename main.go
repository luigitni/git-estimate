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

func main() {
	path := flag.String("repo", ".", "git repository path. If no flag is specified the current folder is assumed")

	estMethod := flag.String("estimate", estSession,
		fmt.Sprintf("estimation method. Accepted values are %q and %q.", estSession, estDay))

	json := flag.Bool("json", false, "if true will output estimates in JSON format")

	baseline := flag.Float64("baseline", 2.0, "baseline value for session estimate")

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
	iter, err := repo.Log(&git.LogOptions{All: true, Order: git.LogOrderCommitterTime})
	if err != nil {
		fmt.Printf("error reading log of repository: %s", err.Error())
		os.Exit(1)
	}

	defer iter.Close()

	byAuthors := make(map[string][]time.Time)

	// first group commits by authors, then for each author count the working days
	if err := iter.ForEach(func(commit *object.Commit) error {

		sl, ok := byAuthors[commit.Author.Email]
		if !ok {
			sl = make([]time.Time, 0)
		}
		sl = append(sl, commit.Author.When)
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
