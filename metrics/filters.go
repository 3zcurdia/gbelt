package metrics

import (
	"strings"

	"github.com/google/go-github/github"
)

// IssuesFilter : filter for issues, leaving empty a value will be default as all
type IssuesFilter struct {
	State   string   `json:"state"`
	Authors []string `json:"authors"`
	Labels  []string `json:"labels"`
}

// Match : returns boolean for the match parametes
func (f *IssuesFilter) Match(issue *github.Issue) bool {
	if issue.GetState() == f.State {
		if len(f.Authors) > 0 {
			if len(f.Labels) > 0 {
				// authors and labels
				return f.containAuthor(issue.User.GetLogin()) && f.containLabels(issue.Labels)
			}
			// authors no labels
			return f.containAuthor(issue.User.GetLogin())
		}
		if len(f.Labels) > 0 {
			// no authors but labels
			return f.containLabels(issue.Labels)
		}
		// no authors no labels
		return true
	}
	return false
}

func (f *IssuesFilter) containAuthor(e string) bool {
	for _, a := range f.Authors {
		if a == e {
			return true
		}
	}
	return false
}

func (f *IssuesFilter) containLabels(labels []github.Label) bool {
	for _, a := range f.Labels {
		for _, lbl := range labels {
			if strings.ToLower(a) == strings.ToLower(lbl.GetName()) {
				return true
			}
		}
	}
	return false
}

func containLabel(labels []*github.Label, name string) bool {
	for _, label := range labels {
		if strings.ToLower(label.GetName()) == strings.ToLower(name) {
			return true
		}
	}
	return false
}
