package api

type context struct {
	Name string
	URL  string
}

func NewContext(name string, url string) Context {
	return &context{
		Name: name,
		URL:  url,
	}
}

func (c *context) GetName() string {
	return c.Name
}
func (c *context) GetURL() string {
	return c.URL
}
