package main

import (
	"fmt"
	"go_exercises/github"
	"log"
)

// Test it: go run . issues repo:golang/go is:open json decoder
func issues_main(args ...string) {
	result, err := github.SearchIssues(args[0:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
	}
}
