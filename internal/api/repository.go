package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"reflect"
)

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
func (r *repository) GetRepositories() (repositoriesResponse *resty.Response, err error) {
	client := resty.New()
	resp, err := client.R().Get("https://api.kubero.dev/v1/repositories")
	return resp, err
}
func (r *repository) GetContexts() (contextsResponse *resty.Response, err error) {
	client := resty.New()
	resp, err := client.R().Get("https://api.kubero.dev/v1/contexts")
	return resp, err
}
func (r *repository) loadContexts() {
	cont, _ := r.GetContexts()
	var contexts []Context
	jsonUnmarshalErr := json.Unmarshal(cont.Body(), &contexts)
	if jsonUnmarshalErr != nil {
		fmt.Println("Error: Unable to load contexts")
		return
	}
	contextSimpleList := make([]string, len(contexts))
	for _, contextt := range contexts {
		contextSimpleList = append(contextSimpleList, contextt.GetName())
	}
}
func (r *repository) loadRepositories() {
	res, err := r.GetRepositories()
	if res == nil {
		fmt.Println("Error: Can't reach Kubero API. Make sure, you are logged in.")
		os.Exit(1)
	}
	if res.StatusCode() != 200 {
		fmt.Println("Error:", res.StatusCode(), "Can't reach Kubero API. Make sure, you are logged in.")
		os.Exit(1)
	}
	if err != nil {
		fmt.Println("Error: Unable to load repositories")
		fmt.Println(err)
		os.Exit(1)
	}
	var availRep []Repository
	jsonUnmarshalErr := json.Unmarshal(res.Body(), &availRep)
	if jsonUnmarshalErr != nil {
		fmt.Println("Error: Unable to load repositories")
		return
	}
	t := reflect.TypeOf(availRep)
	repoSimpleList := make([]string, t.NumField())
	for i := range repoSimpleList {
		if reflect.ValueOf(availRep).Field(i).Bool() {
			repoSimpleList[i] = t.Field(i).Name
		}
	}
}
