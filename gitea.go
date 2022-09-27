package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrorApi = errors.New("API error")

type org struct {
	Username string `json:"username"`
}

type repo struct {
	CloneUrl string `json:"clone_url"`
}

type user struct {
	Login string `json:"login"`
}

type Gitea struct {
	url string
	key string
}

func (g *Gitea) getOrgs() ([]string, error) {
	var ret []string

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/orgs?access_token=%s", g.url, g.key))
	if err != nil {
		return ret, ErrorApi
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ret, ErrorApi
	}

	body, _ := io.ReadAll(resp.Body)

	var res []org
	json.Unmarshal([]byte(body), &res)

	for _, i := range res {
		ret = append(ret, i.Username)
	}
	return ret, nil
}

func (g *Gitea) getUsers() ([]string, error) {
	var ret []string

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/admin/users?access_token=%s", g.url, g.key))
	if err != nil {
		return ret, ErrorApi
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ret, ErrorApi
	}

	body, _ := io.ReadAll(resp.Body)

	var res []user
	json.Unmarshal([]byte(body), &res)

	for _, i := range res {
		ret = append(ret, i.Login)
	}
	return ret, nil
}

func (g *Gitea) getOrgReposPage(org string, page int) ([]string, error) {
	var ret []string

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/orgs/%s/repos?page=%d&access_token=%s", g.url, org, page, g.key))
	if err != nil {
		return ret, ErrorApi
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ret, ErrorApi
	}

	body, _ := io.ReadAll(resp.Body)

	var res []repo
	json.Unmarshal([]byte(body), &res)

	for _, i := range res {
		ret = append(ret, i.CloneUrl)
	}

	return ret, nil
}

func (g *Gitea) getOrgRepos(org string) ([]string, error) {
	var ret []string

	for i := 0; ; i++ {
		r, e := g.getOrgReposPage(org, i)
		if e != nil {
			return ret, e
		}
		if len(r) == 0 {
			break
		}
		ret = append(ret, r...)
	}

	return ret, nil
}

func (g *Gitea) getUserReposPage(user string, page int) ([]string, error) {
	var ret []string

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/users/%s/repos?page=%d&access_token=%s", g.url, user, page, g.key))
	if err != nil {
		return ret, ErrorApi
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return ret, ErrorApi
	}

	body, _ := io.ReadAll(resp.Body)

	var res []repo
	json.Unmarshal([]byte(body), &res)

	for _, i := range res {
		ret = append(ret, i.CloneUrl)
	}

	return ret, nil
}

func (g *Gitea) getUserRepos(user string) ([]string, error) {
	var ret []string

	for i := 0; ; i++ {
		r, e := g.getUserReposPage(user, i)
		if e != nil {
			return ret, e
		}
		if len(r) == 0 {
			break
		}
		ret = append(ret, r...)
	}

	return ret, nil
}

func (g *Gitea) Repos() ([]string, error) {
	var ret []string

	orgs, e := g.getOrgs()
	if e != nil {
		return ret, e
	}
	users, e := g.getUsers()
	if e != nil {
		return ret, e
	}

	// org repos
	for _, org := range orgs {
		r, e := g.getOrgRepos(org)
		if e != nil {
			return ret, e
		}
		ret = append(ret, r...)
	}

	// user repos
	for _, org := range users {
		r, e := g.getUserRepos(org)
		if e != nil {
			return ret, e
		}
		ret = append(ret, r...)
	}

	return ret, nil
}

func (g *Gitea) SetUrl(url string) {
	g.url = url
}

func (g *Gitea) SetApiKey(key string) {
	g.key = key
}
