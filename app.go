package main

import (
	"fmt"

	"github.com/3zcurdia/gbelt/metrics"
	"github.com/3zcurdia/gbelt/search"
	"github.com/google/go-github/github"
)

func main() {
	fmt.Println("Start searching github user")
	users, _ := search.UserByName("ezcurdia")

	um := metrics.NewUserMetrics(users[0].Username)
	fmt.Printf("Github User: %s\n", um.Username)
	fmt.Printf("Name: %s\n", um.Name)
	fmt.Printf("Email: %s\n", um.Email)
	fmt.Printf("Location: %s\n", um.Location)
	fmt.Printf("Followers: %d\n", um.Followers)
	fmt.Printf("Autored Repos: %d\n", um.AutoredRepos)

	um.FetchLanguagesCount(true)
	fmt.Printf("Stars: %d\n", um.Stars)
	fmt.Printf("Langauges: %v\n", um.Languages)

	fmt.Println("======================================")

	rm := metrics.NewRepoMetrics("stretchr", "testify")
	fmt.Printf("Github Repository: '%s/%s'\n", rm.Owner, rm.Name)
	fmt.Printf("Stars: %d  Forks: %d\n", rm.Stars, rm.Forks)
	fmt.Printf("Language: %s\n", rm.MainLanguage)

	commit, _ := rm.FetchLastCommit()
	fmt.Printf("Last commit date: %v\n", commit.Author.GetDate())

	count, _ := rm.FetchContributorsCount()
	fmt.Printf("Contributors: %d\n", count)

	rm.FetchOpenIssues()
	fmt.Printf("Open issues: %v\n", rm.IssuesOpen)

	rm.FetchClosedIssues()
	fmt.Printf("Closed issues: %v\n", rm.IssuesClosed)

	opt := &github.IssueListByRepoOptions{
		State:       "closed",
		Sort:        "closed_at",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	trendsMap, _ := rm.FetchStatsPer(opt)
	fmt.Println("Weekly Speeds in 2017")
	for week, trends := range trendsMap[2017] {
		fmt.Printf("  * %2d => %.3f [h/issue]\n", week, trends.Avg())
	}
	fmt.Printf("Project speed: %.3f [h/issue]\n", rm.Speed)
}
