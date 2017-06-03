package metrics

// UserMetrics : struct for github user metrics
type UserMetrics struct {
	Username     string         `json:"username"`
	Stars        int            `json:"stars"`
	AutoredRepos int            `json:"authored_repos"`
	Languages    map[string]int `json:"languages"`
}

// NewUserMetrics : initialize user metrics given a username
func NewUserMetrics(name string) UserMetrics {
	m := UserMetrics{Username: name, Stars: 0, AutoredRepos: 0}
	m.Languages = make(map[string]int)
	return m
}
