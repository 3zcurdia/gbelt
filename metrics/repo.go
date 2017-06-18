package metrics

import (
	"github.com/google/go-github/github"
)

// FetchContributorsCount : fetch form github the contributors count
func (r *RepoMetrics) FetchContributorsCount() (int, error) {
	contribs, _, err := r.client.Repositories.ListCollaborators(r.ctx, r.Owner, r.Name, nil)
	if err != nil {
		return r.ContributorsCount, err
	}
	r.ContributorsCount = len(contribs)
	return r.ContributorsCount, nil
}

// FetchSpeedPer : fetch speed of the project
func (r *RepoMetrics) FetchSpeedPer(opt *github.IssueListByRepoOptions) (float64, error) {
	issues, _, err := r.client.Issues.ListByRepo(r.ctx, r.Owner, r.Name, opt)
	if err != nil {
		return 0, err
	}
	if opt.State == "closed" && len(opt.Labels) == 0 {
		r.IssuesClosed = len(issues)
	}
	for _, issue := range issues {
		elapsed := issue.ClosedAt.Sub(*issue.CreatedAt)
		// year, week := issue.ClosedAt.ISOWeek()
		// fmt.Printf("%v (%v-%v) : %v \n", issue.GetNumber(), year, week, elapsed)
		r.trends.Add(elapsed.Hours())
	}
	return r.trends.Avg(), nil
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

// FetchOpenIssues : fetch all open issues and count total
func (r *RepoMetrics) FetchOpenIssues() (int, error) {
	opt := &github.IssueListByRepoOptions{
		State:       "open",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	r.IssuesOpen = 0
	for {
		issuesOpen, resp, err := r.client.Issues.ListByRepo(r.ctx, r.Owner, r.Name, opt)
		if err != nil {
			return 0, err
		}
		r.IssuesOpen += len(issuesOpen)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return r.IssuesOpen, nil
}
