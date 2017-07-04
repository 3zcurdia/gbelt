package metrics

import (
	"context"
	"log"
	"os"

	"github.com/3zcurdia/fastrends"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// ChannelError : channel to handle errors async
type ChannelError chan error

// UserMetrics : struct for github user metrics
type UserMetrics struct {
	client       *github.Client
	ctx          context.Context
	user         *github.User
	repos        []*RepoMetrics
	Username     string         `json:"username"`
	Name         string         `json:"name"`
	Email        string         `json:"email"`
	Location     string         `json:"location"`
	Stars        int            `json:"stars"`
	Followers    int            `json:"followers"`
	AutoredRepos int            `json:"authored_repos"`
	Languages    map[string]int `json:"languages"`
}

// RepoMetrics : struct for repository metrics
type RepoMetrics struct {
	client            *github.Client
	ctx               context.Context
	repo              *github.Repository
	trends            *fastrends.TrendFloat64
	issues            []*github.Issue
	Name              string         `json:"name"`
	Owner             string         `json:"owner"`
	Stars             int            `json:"stars"`
	Forks             int            `json:"forks"`
	ContributorsCount int            `json:"contributors"`
	MainLanguage      string         `json:"main_language"`
	Languages         map[string]int `json:"languages"`
	Speed             float64        `json:"speed"`
	IssuesOpenCount   int            `json:"issues_open_count"`
	IssuesClosedCount int            `json:"issues_closed_count"`
	LastCommit        *github.Commit
}

// InitGithubClient : initialize github client
func InitGithubClient() (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client, ctx
}

// NewUserMetrics : initialize user metrics given a username
func NewUserMetrics(name string) UserMetrics {
	m := UserMetrics{Username: name, Stars: 0, AutoredRepos: 0}
	m.Languages = make(map[string]int)
	m.client, m.ctx = InitGithubClient()
	user, _, err := m.client.Users.Get(m.ctx, name)
	if err != nil {
		log.Fatalln(err)
		return m
	}
	m.user = user
	m.Email = user.GetEmail()
	m.Name = user.GetName()
	m.Location = user.GetLocation()
	m.Followers = user.GetFollowers()
	_ = m.initRepoMetrics()
	return m
}

func (m *UserMetrics) initRepoMetrics() error {
	opt := &github.RepositoryListOptions{
		Type:        "owner",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	repos, _, err := m.client.Repositories.List(m.ctx, m.Username, opt)
	if err != nil {
		return err
	}
	m.AutoredRepos = len(repos)
	m.repos = make([]*RepoMetrics, 0)
	for _, repo := range repos {
		if repo.GetFork() || repo.GetLanguage() == "" {
			m.AutoredRepos--
			continue
		}
		repoMetric := NewRepoMetrics(m.Username, repo.GetName())
		m.repos = append(m.repos, &repoMetric)
		m.addStars(repoMetric.Stars)
	}
	return nil
}

// NewRepoMetrics : initialize user metrics given a owner and name
func NewRepoMetrics(owner, name string) RepoMetrics {
	m := RepoMetrics{Name: name, Owner: owner}
	m.client, m.ctx = InitGithubClient()
	repo, _, err := m.client.Repositories.Get(m.ctx, owner, name)
	if err != nil {
		log.Fatalln(err)
		return m
	}
	m.repo = repo
	m.Stars = repo.GetStargazersCount()
	m.Forks = repo.GetForksCount()
	m.MainLanguage = repo.GetLanguage()
	m.trends = fastrends.NewTrendFloat64()
	return m
}
