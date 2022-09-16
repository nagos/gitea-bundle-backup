package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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

func (g *Gitea) getOrgs() []string {
	var ret []string

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/orgs?access_token=%s", g.url, g.key))
	if err != nil {
		fmt.Println("API error", err)
		return ret
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var res []org
	json.Unmarshal([]byte(body), &res)

	for _, i := range res {
		ret = append(ret, i.Username)
	}
	return ret
}

func (g *Gitea) getUsers() []string {
	var ret []string

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/admin/users?access_token=%s", g.url, g.key))
	if err != nil {
		fmt.Println("API error", err)
		return ret
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var res []user
	json.Unmarshal([]byte(body), &res)

	for _, i := range res {
		ret = append(ret, i.Login)
	}
	return ret
}

func (g *Gitea) getOrgReposPage(org string, page int) []string {
	var ret []string

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/orgs/%s/repos?page=%d&access_token=%s", g.url, org, page, g.key))
	if err != nil {
		fmt.Println("API error", err)
		return ret
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var res []repo
	json.Unmarshal([]byte(body), &res)

	for _, i := range res {
		ret = append(ret, i.CloneUrl)
	}

	return ret
}

func (g *Gitea) getOrgRepos(org string) []string {
	var ret []string

	for i := 0; ; i++ {
		r := g.getOrgReposPage(org, i)
		if len(r) == 0 {
			break
		}
		ret = append(ret, r...)
	}

	return ret
}

func (g *Gitea) getUserReposPage(user string, page int) []string {
	var ret []string

	resp, err := http.Get(fmt.Sprintf("%s/api/v1/users/%s/repos?page=%d&access_token=%s", g.url, user, page, g.key))
	if err != nil {
		fmt.Println("API error", err)
		return ret
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var res []repo
	json.Unmarshal([]byte(body), &res)

	for _, i := range res {
		ret = append(ret, i.CloneUrl)
	}

	return ret
}

func (g *Gitea) getUserRepos(user string) []string {
	var ret []string

	for i := 0; ; i++ {
		r := g.getUserReposPage(user, i)
		if len(r) == 0 {
			break
		}
		ret = append(ret, r...)
	}

	return ret
}

func (g *Gitea) Repos() []string {
	var ret []string

	for _, org := range g.getOrgs() {
		ret = append(ret, g.getOrgRepos(org)...)
	}

	for _, org := range g.getUsers() {
		ret = append(ret, g.getUserRepos(org)...)
	}

	return ret
}

func (g *Gitea) SetUrl(url string) {
	g.url = url
}

func (g *Gitea) SetApiKey(key string) {
	g.key = key
}
