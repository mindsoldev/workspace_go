package main

import (
	"fmt"
	"go_exercises/github"
	"io"
	"log"
	"os"
	"text/template"
	"time"
)

var stdoutForIssues io.Writer = os.Stdout

const templ = `{{.TotalCount}} issues:
{{range .Items}}----------------------------------------
Number:	{{.Number}}
User:	{{.User.Login}}
Title:	{{.Title | printf "%d.64s"}}
Age:	{{.CreatedAt | daysAgo}} days
{{end}}`

var report = template.Must(template.New("issuelist").
	Funcs(template.FuncMap{"daysAgo": daysAgo}).
	Parse(templ))

// Test it: go run . issuesByAges repo:golang/go is:open json decoder
func issuesByAges_main(args ...string) {
	result, err := github.SearchIssues(args[0:])
	if err != nil {
		log.Fatal(err)
	}
	printByDaysAgo(result)
}

func daysAgo(t time.Time) int {
	return int(time.Since(t).Hours() / 24)
}

func printByDaysAgo(result *github.IssuesSearcResult) {
	now := time.Now()

	day := now.AddDate(0, 0, -1)
	month := now.AddDate(0, -1, 0)
	year := now.AddDate(-1, 0, 0)

	var dayIssues github.IssuesSearcResult
	var monthIssues github.IssuesSearcResult
	var yearIssues github.IssuesSearcResult
	var otherIssues github.IssuesSearcResult

	for _, issue := range result.Items {
		switch {
		case issue.CreatedAt.After(day):
			(dayIssues.TotalCount)++
			dayIssues.Items = append(dayIssues.Items, issue)
		case issue.CreatedAt.After(month):
			(monthIssues.TotalCount)++
			monthIssues.Items = append(monthIssues.Items, issue)
		case issue.CreatedAt.After(year):
			(yearIssues.TotalCount)++
			yearIssues.Items = append(yearIssues.Items, issue)
		default:
			(otherIssues.TotalCount)++
			otherIssues.Items = append(otherIssues.Items, issue)
		}
	}
	fmt.Fprint(stdoutForIssues, "\n----------------------------------------\n")
	fmt.Fprint(stdoutForIssues, "Last day:\t")
	if err := report.Execute(stdoutForIssues, dayIssues); err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(stdoutForIssues, "\n----------------------------------------\n")
	fmt.Fprint(stdoutForIssues, "Last month:\t")
	if err := report.Execute(stdoutForIssues, monthIssues); err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(stdoutForIssues, "\n----------------------------------------\n")
	fmt.Fprint(stdoutForIssues, "Last year:\t")
	if err := report.Execute(stdoutForIssues, yearIssues); err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(stdoutForIssues, "\n----------------------------------------\n")
	fmt.Fprint(stdoutForIssues, "Long long ago:\t")
	if err := report.Execute(stdoutForIssues, otherIssues); err != nil {
		log.Fatal(err)
	}
}
