package main

import (
	"fmt"

	"github.com/3zcurdia/gbelt/metrics"
	"github.com/3zcurdia/gbelt/search"
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

	issues := rm.Issues()
	fmt.Printf("All issues: %v\n", len(issues))

	rm.IssuesOpen()
	fmt.Printf("Open issues: %v\n", rm.IssuesOpenCount)

	rm.IssuesClosed()
	fmt.Printf("Closed issues: %v\n", rm.IssuesClosedCount)

	trendsMap := rm.FetchStats()
	fmt.Println("Weekly Speeds in 2017")
	for week, trends := range trendsMap[2017] {
		fmt.Printf("  * %2d => %.3f [h/issue]\n", week, trends.Avg())
	}
	fmt.Printf("Project speed: %.3f [h/issue]\n", rm.Speed)

	filter := &metrics.IssuesFilter{
		State:  "closed",
		Labels: []string{"bug"},
	}
	trendsMapFilter := rm.FetchStatsBy(filter)
	fmt.Println("Weekly closed bugs in 2017")
	for week, trends := range trendsMapFilter[2017] {
		fmt.Printf("  * %2d => %d [bug]\n", week, trends.Count)
	}
}
