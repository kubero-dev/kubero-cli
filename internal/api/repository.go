package api

type repository struct {
	Name string
	URL  string
}

func NewRepository(name string, url string) Repository {
	return &repository{
		Name: name,
		URL:  url,
	}
}

func (r *repository) GetName() string {
	return r.Name
}
func (r *repository) GetURL() string {
	return r.URL
}
func (r *repository) GetRepositories() ([]Repository, error) {
	return []Repository{r}, nil
}
func (r *repository) GetContexts() ([]Context, error) {
	return []Context{r}, nil
}
