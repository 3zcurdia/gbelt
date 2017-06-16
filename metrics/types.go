package metrics

import (
	"context"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// ChannelError : channel to handle errors async
type ChannelError chan error

// UserMetrics : struct for github user metrics
type UserMetrics struct {
	Username     string         `json:"username"`
	Stars        int            `json:"stars"`
	AutoredRepos int            `json:"authored_repos"`
	Languages    map[string]int `json:"languages"`
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
	return m
}
