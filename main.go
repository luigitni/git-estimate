package main

import (
	"flag"
	"fmt"
	"git-estimate/estimate"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	estSession   = "session"
	estDay       = "day"
	jiraPattern  = "jira"
	typePattern  = "type"
	scopePattern = "scope"
)

const timeArgsLayout = "2006-01-02T15-04"

func main() {
	path := flag.String("repo", ".", "git repository path. If no flag is specified the current folder is assumed")

	estMethod := flag.String("estimate", estSession,
		fmt.Sprintf("estimation method. Accepted values are %q and %q.", estSession, estDay))

	json := flag.Bool("json", false, "if true will output estimates in JSON format")

	group := flag.String("group", "", "group estimates based on comment message content using a predefined or custom pattern. "+
		"Custom patterns should identify exactly 1 capturing group. See https://github.com/google/re2/wiki/Syntax for syntax.\n"+
		"Predefined patterns available:\n"+
		"\n\tjira - Captures the first Jira issue key based on the smart commit format (https://support.atlassian.com/bitbucket-cloud/docs/use-smart-commits/)"+
		"\n\ttype - Captures the type component of conventional commit messages (https://www.conventionalcommits.org/en/v1.0.0/)"+
		"\n\tscope - Captures the scope component of conventional commit messages (https://www.conventionalcommits.org/en/v1.0.0/)"+
		"\n")

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
		fmt.Printf("Invalid estimation method. Accepted values are %q and %q", estSession, estDay)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error opening repository at %s: %s", *path, err.Error())
		os.Exit(1)
	}

	var re *regexp.Regexp
	groupSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "group" {
			groupSet = true
		}
	})

	if groupSet {
		switch *group {
		case jiraPattern:
			re = regexp.MustCompile("(?:[\\w:]+ )?([a-zA-Z]\\w+-\\d+)\\D")
		case typePattern:
			re = regexp.MustCompile("^([a-zA-Z!]+)[\\(:]")
		case scopePattern:
			re = regexp.MustCompile("^[^\\()]+\\(([^\\)]+)\\)")
		case "":
			fmt.Printf("Invalid grouping pattern. Provide a custom pattern or specify a valid preset pattern: %q", jiraPattern)
			os.Exit(1)
		default:
			re = regexp.MustCompile(*group)
			caps := re.NumSubexp()
			if caps == 0 || caps > 1 {
				fmt.Printf("Invalid grouping pattern. Pattern must specify exactly 1 capture group")
				os.Exit(1)
			}
		}
	}

	// get the commit history
	var start time.Time
	if *from != "" {
		// parse the string
		start, err = time.Parse(timeArgsLayout, *from)
		if err != nil {
			fmt.Printf("Unable to parse 'from' %s given. %s", *from, err.Error())
			os.Exit(1)
		}
	}

	iter, err := repo.Log(&git.LogOptions{All: true, Order: git.LogOrderCommitterTime})
	if err != nil {
		fmt.Printf("Error reading log of repository: %s", err.Error())
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

		var builder strings.Builder
		builder.WriteString(commit.Author.Email)

		if *group != "" {
			matches := re.FindStringSubmatch(commit.Message)
			if len(matches) > 1 {
				builder.WriteString("@")
				builder.WriteString(matches[1])
			}
		}

		key := builder.String()

		sl, ok := byAuthors[key]
		if !ok {
			sl = make([]time.Time, 0)
		}

		sl = append(sl, when)
		byAuthors[key] = sl

		return nil
	}); err != nil {
		fmt.Printf("Error reading commits: %s", err.Error())
		os.Exit(1)
	}

	// sort each slice of commits by date
	for k := range byAuthors {
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
