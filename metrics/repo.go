package metrics

import (
	"github.com/3zcurdia/fastrends"
	"github.com/google/go-github/github"
)

// FetchLastCommit : Fetch last commit from master
func (r *RepoMetrics) FetchLastCommit() (*github.Commit, error) {
	branch, _, err := r.client.Repositories.GetBranch(r.ctx, r.Owner, r.Name, "master")
	if err != nil {
		return r.LastCommit, err
	}
	r.LastCommit = branch.Commit.Commit
	return r.LastCommit, nil
}

// FetchContributorsCount : fetch form github the contributors count
func (r *RepoMetrics) FetchContributorsCount() (int, error) {
	r.ContributorsCount = 0
	opt := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		contribs, resp, err := r.client.Repositories.ListContributors(r.ctx, r.Owner, r.Name, opt)
		if err != nil {
			return r.ContributorsCount, err
		}
		r.ContributorsCount += len(contribs)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return r.ContributorsCount, nil
}

// Issues : lazy load all repo issues
func (r *RepoMetrics) Issues() []*github.Issue {
	if len(r.issues) > 0 {
		return r.issues
	}
	opt := &github.IssueListByRepoOptions{
		State:       "all",
		Sort:        "created_at",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	r.issues, _ = r.FetchAllIssues(opt)
	return r.issues
}

func (r *RepoMetrics) IssuesFiltered(filter *IssuesFilter) []*github.Issue {
	filtered := make([]*github.Issue, 0)
	for _, issue := range r.Issues() {
		if filter.Match(issue) {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

// IssuesClosed : fetch all closed issues return the trends map per year and week
func (r *RepoMetrics) IssuesClosed() []*github.Issue {
	issues := r.IssuesFiltered(&IssuesFilter{State: "closed"})
	r.IssuesClosedCount = len(issues)
	return issues
}

// IssuesOpen : fetch all open issues and count total
func (r *RepoMetrics) IssuesOpen() []*github.Issue {
	issues := r.IssuesFiltered(&IssuesFilter{State: "open"})
	r.IssuesOpenCount = len(issues)
	return issues
}

// FetchAllIssues : fetch all issues
func (r *RepoMetrics) FetchAllIssues(opt *github.IssueListByRepoOptions) ([]*github.Issue, error) {
	allIssues := make([]*github.Issue, 0)
	for {
		issues, resp, err := r.client.Issues.ListByRepo(r.ctx, r.Owner, r.Name, opt)
		if err != nil {
			return allIssues, err
		}
		for _, issue := range issues {
			allIssues = append(allIssues, issue)
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allIssues, nil
}

// FetchStats : fetch stats from closed issues in the project
func (r *RepoMetrics) FetchStats() map[int]map[int]*fastrends.TrendFloat64 {
	r.Speed = 0
	stats := make(map[int]map[int]*fastrends.TrendFloat64)
	for _, issue := range r.IssuesClosed() {
		elapsed := issue.ClosedAt.Sub(*issue.CreatedAt)
		year, week := issue.ClosedAt.ISOWeek()
		if _, oky := stats[year]; oky {
			if _, okw := stats[year][week]; !okw {
				stats[year][week] = fastrends.NewTrendFloat64()
			}
		} else {
			stats[year] = make(map[int]*fastrends.TrendFloat64)
			stats[year][week] = fastrends.NewTrendFloat64()
		}
		stats[year][week].Add(elapsed.Hours())
		r.trends.Add(elapsed.Hours())
	}
	r.Speed = r.trends.WeightedAvg()
	return stats
}

// FetchStatsBy : fetch stats from the filter pointer
func (r *RepoMetrics) FetchStatsBy(filter *IssuesFilter) map[int]map[int]*fastrends.TrendFloat64 {
	stats := make(map[int]map[int]*fastrends.TrendFloat64)
	for _, issue := range r.IssuesFiltered(filter) {
		elapsed := issue.ClosedAt.Sub(*issue.CreatedAt)
		year, week := issue.ClosedAt.ISOWeek()
		if _, oky := stats[year]; oky {
			if _, okw := stats[year][week]; !okw {
				stats[year][week] = fastrends.NewTrendFloat64()
			}
		} else {
			stats[year] = make(map[int]*fastrends.TrendFloat64)
			stats[year][week] = fastrends.NewTrendFloat64()
		}
		stats[year][week].Add(elapsed.Hours())
	}
	return stats
}

func (r *RepoMetrics) fetchLanguages() (map[string]int, error) {
	r.MainLanguage = r.repo.GetLanguage()
	langs, _, err := r.client.Repositories.ListLanguages(r.ctx, r.Owner, r.Name)
	if err != nil {
		return r.Languages, err
	}
	r.Languages = langs
	return r.Languages, nil
}
