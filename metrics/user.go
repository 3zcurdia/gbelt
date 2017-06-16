package metrics

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/github"
)

// GetLanguagesCount : count languages lines of code
func (m *UserMetrics) GetLanguagesCount(detail bool) map[string]int {
	client, ctx := InitGithubClient()
	languages, err := m.languagesCount(ctx, client.Repositories, detail)
	if err != nil {
		panic(err)
	}
	return languages
}

func (m *UserMetrics) languagesCount(ctx context.Context, service *github.RepositoriesService, detail bool) (map[string]int, error) {
	m.Stars = 0
	m.Languages = make(map[string]int)
	opt := &github.RepositoryListOptions{Type: "owner", Sort: "updated", Direction: "desc"}
	repos, _, err := service.List(ctx, m.Username, opt)
	if err != nil {
		return m.Languages, err
	}
	m.AutoredRepos = len(repos)
	errc := make(ChannelError)
	lngc := make(chan map[string]int)
	for _, repo := range repos {
		if repo.GetFork() || repo.GetLanguage() == "" {
			m.AutoredRepos--
			continue
		}
		m.addStar(repo)
		if repo.GetLanguage() != "" {
			if detail {
				go m.fetchLenguageLines(ctx, service, repo, lngc, errc)
			} else {
				m.addCount(repo.GetLanguage(), 1)
			}
		}
	}
	if detail {
		err := m.listenLenguageLines(lngc, errc)
		if err != nil {
			return m.Languages, err
		}
	}

	return m.Languages, nil
}

func (m *UserMetrics) fetchLenguageLines(ctx context.Context, service *github.RepositoriesService, repo *github.Repository, lngc chan map[string]int, errc ChannelError) {
	langs, _, err := service.ListLanguages(ctx, m.Username, repo.GetName())
	if err != nil {
		errc <- err
	} else {
		lngc <- langs
	}
	return
}

// Listen all go routines for each repo language url
func (m *UserMetrics) listenLenguageLines(lngc chan map[string]int, errc ChannelError) error {
	reposLeft := m.AutoredRepos
	var err error
	for {
		select {
		case res := <-lngc:
			m.addCountHash(res)
			reposLeft--
		case err = <-errc:
			log.Fatalln(err)
			reposLeft--
			break
		case <-time.After(1000 * time.Millisecond):
			log.Fatalln("Timeout")
		}
		if reposLeft <= 0 {
			break
		}
	}
	return err
}

func (m *UserMetrics) addStar(repo *github.Repository) {
	stars := repo.GetStargazersCount()
	if stars > 1 {
		m.Stars += stars
	}
}

func (m *UserMetrics) addCount(lang string, value int) {
	lang = normalizeLang(lang)
	_, ok := m.Languages[lang]
	if ok {
		m.Languages[lang] += value
	} else {
		m.Languages[lang] = value
	}
}

func (m *UserMetrics) addCountHash(langLines map[string]int) {
	if len(langLines) == 0 {
		return
	}
	for lang, value := range langLines {
		m.addCount(lang, value)
	}
}

func normalizeLang(key string) string {
	switch key {
	case "Emacs Lisp":
		return "Lisp"
	case "C++":
		return "CPlusPlus"
	case "C#":
		return "CSharp"
	default:
		return key
	}
}
