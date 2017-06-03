package metrics

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func (m *UserMetrics) GetLanguagesCount(detail bool) map[string]int {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "e4065f79d98e4c4c345e29215e089f8f2f88718c"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	opt := &github.RepositoryListOptions{Type: "owner", Sort: "updated", Direction: "desc"}

	repos, _, err := client.Repositories.List(ctx, m.Username, opt)
	if err != nil {
		panic(err)
	}
	m.AutoredRepos = len(repos)
	return m.languagesCount(repos, detail)
}

func (m *UserMetrics) languagesCount(repos []*github.Repository, detail bool) map[string]int {
	m.Stars = 0
	for _, repo := range repos {
		if repo.GetFork() {
			m.AutoredRepos--
			continue
		}
		m.addStar(repo)
		if repo.GetLanguage() != "" {
			if detail {
				fmt.Println(repo.GetLanguagesURL())
			} else {
				m.addCount(repo.GetLanguage())
			}
		}
	}
	return m.Languages
}

func (m *UserMetrics) addStar(repo *github.Repository) {
	stars := repo.GetStargazersCount()
	if stars > 1 {
		m.Stars += stars
	}
}

func (m *UserMetrics) addCount(lang string) {
	_, ok := m.Languages[lang]
	if ok {
		m.Languages[lang]++
	} else {
		m.Languages[lang] = 1
	}
}

func (m *UserMetrics) addCountHash(langLines map[string]int) {
}
