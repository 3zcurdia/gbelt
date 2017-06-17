package metrics

// func (r *RepoMetrics) GetContributorsCount() {
//
// }
//
// func (r *RepoMetrics) GetSpeed() {
//
// }
//
// func (r *RepoMetrics) GetSpeedPer(author, label string) {
//
// }

func (r *RepoMetrics) fetchLenguages() error {
	r.MainLanguage = r.repo.GetLanguage()
	langs, _, err := r.client.Repositories.ListLanguages(r.ctx, r.Owner, r.Name)
	if err != nil {
		return err
	}
	r.Languages = langs
	return nil
}
