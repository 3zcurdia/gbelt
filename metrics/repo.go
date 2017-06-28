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

// FetchOpenIssues : fetch all open issues and count total
func (r *RepoMetrics) FetchOpenIssues() ([]*github.Issue, error) {
	opt := &github.IssueListByRepoOptions{
		State:       "open",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	issues, err := r.FetchAllIssues(opt)
	if err != nil {
		return issues, err
	}
	r.IssuesOpen = len(issues)
	return issues, nil
}

// FetchClosedIssues : fetch all closed issues return the trends map per year and week
func (r *RepoMetrics) FetchClosedIssues() ([]*github.Issue, error) {
	opt := &github.IssueListByRepoOptions{
		State:       "closed",
		Sort:        "closed_at",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	issues, err := r.FetchAllIssues(opt)
	if err != nil {
		return issues, err
	}
	r.IssuesClosed = len(issues)
	return issues, nil
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

// FetchStatsPer : fetch speed of the project
func (r *RepoMetrics) FetchStatsPer(opt *github.IssueListByRepoOptions) (map[int]map[int]*fastrends.TrendFloat64, error) {
	allClosed := opt.State == "closed" && len(opt.Labels) == 0
	if allClosed {
		r.Speed = 0
	}
	stats := make(map[int]map[int]*fastrends.TrendFloat64)
	for {
		issues, resp, err := r.client.Issues.ListByRepo(r.ctx, r.Owner, r.Name, opt)
		if err != nil {
			return stats, err
		}
		for _, issue := range issues {
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
			if allClosed {
				r.trends.Add(elapsed.Hours())
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	if allClosed {
		r.Speed = r.trends.Avg()
	}
	return stats, nil
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
