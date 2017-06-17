package metrics

import (
	"context"
	"os"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// ChannelError : channel to handle errors async
type ChannelError chan error

// UserMetrics : struct for github user metrics
type UserMetrics struct {
	client       *github.Client
	ctx          context.Context
	Username     string         `json:"username"`
	Stars        int            `json:"stars"`
	AutoredRepos int            `json:"authored_repos"`
	Languages    map[string]int `json:"languages"`
}

// RepoMetrics : struct for repository metrics
type RepoMetrics struct {
	client     *github.Client
	ctx        context.Context
	Name       string    `json:"name"`
	Owner      string    `json:"owner"`
	Stars      int       `json:"stars"`
	LastCommit time.Time `json:"date"`
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
	return m
}

// NewRepoMetrics : initialize user metrics given a owner and name
func NewRepoMetrics(owner, name string) RepoMetrics {
	m := RepoMetrics{Name: name, Owner: owner}
	m.client, m.ctx = InitGithubClient()
	return m
}
