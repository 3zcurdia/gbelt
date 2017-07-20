package metrics

import (
	"log"
	"time"
)

// FetchProfile : load user profile info
func (um *UserMetrics) FetchProfile() error {
	user, _, err := um.client.Users.Get(um.ctx, um.Name)
	if err != nil {
		return err
	}
	um.user = user
	um.Email = user.GetEmail()
	um.Name = user.GetName()
	um.Location = user.GetLocation()
	um.Followers = user.GetFollowers()
	return nil
}

// FetchLanguagesCount : count languages lines of code
func (um *UserMetrics) FetchLanguagesCount(detail bool) (map[string]int, error) {
	um.Languages = make(map[string]int)
	errc := make(ChannelError)
	lngc := make(chan map[string]int)
	for _, repo := range um.repos {
		if detail {
			go um.fetchLanguages(repo, lngc, errc)
		} else {
			um.addCount(repo.MainLanguage, 1)
		}
	}
	if detail {
		err := um.listenLenguageLines(lngc, errc)
		if err != nil {
			return um.Languages, err
		}
	}
	return um.Languages, nil
}

func (um *UserMetrics) fetchLanguages(rm *RepoMetrics, lngc chan map[string]int, errc ChannelError) {
	langs, err := rm.fetchLanguages()
	if err != nil {
		errc <- err
	} else {
		lngc <- langs
	}
	return
}

// Listen all go routines for each repo language url
func (um *UserMetrics) listenLenguageLines(lngc chan map[string]int, errc ChannelError) error {
	reposLeft := um.AutoredRepos
	var err error
	for {
		select {
		case res := <-lngc:
			um.addCountHash(res)
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

func (um *UserMetrics) addCount(lang string, value int) {
	lang = normalizeLang(lang)
	_, ok := um.Languages[lang]
	if ok {
		um.Languages[lang] += value
	} else {
		um.Languages[lang] = value
	}
}

func (um *UserMetrics) addCountHash(langLines map[string]int) {
	if len(langLines) == 0 {
		return
	}
	for lang, value := range langLines {
		um.addCount(lang, value)
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
