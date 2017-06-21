package search

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// UserByEmail search in github by email
func UserByEmail(email string) ([]User, error) {
	return UserByTerm(email, "email")
}

// UserByName search in github by name
func UserByName(name string) ([]User, error) {
	return UserByTerm(name, "name")
}

// UserByTerm search in github by given term
func UserByTerm(q, term string) ([]User, error) {
	return By(q, term, "Users")
}

// By search in github by given term and type
func By(q, term, searchType string) ([]User, error) {
	params := make(map[string]string)
	ghURL, err := buildURLWithQuery("https://github.com/search", params)
	if err != nil {
		return nil, err
	}

	query := ghURL.Query()
	query.Set("q", q+" in:"+term)
	query.Set("type", searchType)
	ghURL.RawQuery = query.Encode()

	root, err := getHTML(ghURL)
	if err != nil {
		return nil, err
	}
	return scrapeUsers(root)
}

func userItemDiv(n *html.Node) bool {
	if n.DataAtom == atom.Div {
		return strings.Contains(scrape.Attr(n, "class"), "user-list-item")
	}
	return false
}

func usernameA(n *html.Node) bool {
	if n.DataAtom == atom.A && n.Parent != nil && n.Parent.DataAtom == atom.Div {
		return strings.Contains(scrape.Attr(n.Parent, "class"), "user-list-info")
	}
	return false
}

func userEmailA(n *html.Node) bool {
	return n.DataAtom == atom.A && n.Parent != nil && n.Parent.DataAtom == atom.Li
}

func parseUser(n *html.Node) User {
	user := User{}
	usernameNode, ok := scrape.Find(n, usernameA)
	if ok {
		user.Username = scrape.Text(usernameNode)
	}
	userEmailNode, ok := scrape.Find(n, userEmailA)
	if ok {
		user.Email = scrape.Text(userEmailNode)
	}
	return user
}

func scrapeUsers(root *html.Node) ([]User, error) {
	nodes := scrape.FindAll(root, userItemDiv)
	contacts := []User{}
	for _, node := range nodes {
		contacts = append(contacts, parseUser(node))
	}
	return contacts, nil
}

func getHTML(uri *url.URL) (*html.Node, error) {
	resp, err := http.Get(uri.String())
	if err != nil {
		return nil, err
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	return root, nil
}

func buildURLWithQuery(s string, params map[string]string) (*url.URL, error) {
	ghURL, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	query := ghURL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	ghURL.RawQuery = query.Encode()
	return ghURL, nil
}
