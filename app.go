package main

import (
	"fmt"

	"github.com/3zcurdia/gbelt/metrics"
	"github.com/google/go-github/github"
)

func main() {
	um := metrics.NewUserMetrics("3zcurdia")
	langs, _ := um.FetchLanguagesCount(true)
	fmt.Printf("user langs : %v\n", langs)

	rm := metrics.NewRepoMetrics("hubbubhealth", "rcs")
	opened, _ := rm.FetchOpenIssues()
	fmt.Printf("open issues : %v\n", opened)

	opt := &github.IssueListByRepoOptions{
		State:       "closed",
		Sort:        "closed_at",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	avgHours, _ := rm.FetchSpeedPer(opt)
	fmt.Printf("closed issues : %v\n", rm.IssuesClosed)
	fmt.Printf("avg hour : %v\n", avgHours)
}
