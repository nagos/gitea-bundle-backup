package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var gitea Gitea

func myUsage() {
	fmt.Printf("Usage: %s [OPTIONS] URL APIKEY ...\n", os.Args[0])
	flag.PrintDefaults()
}

func parseArgs() bool {
	flag.Usage = myUsage
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		return true
	}

	gitea.SetUrl(flag.Arg(0))
	gitea.SetApiKey(flag.Arg(1))
	return false
}

func gitBundle(repo string) {
	fmt.Println(repo)
	url := strings.Split(repo, "/")
	repo_name := url[len(url)-1]
	repo_user := url[len(url)-2]
	tmp_dir, err := ioutil.TempDir("", "gitea_backup")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmp_dir)

	cmd := exec.Command("git", "clone", "--mirror", repo, tmp_dir)
	_, err = cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cwd, _ := os.Getwd()
	cmd = exec.Command("git", "bundle", "create", fmt.Sprintf("%s/%s_%s", cwd, repo_user, repo_name), "--all")
	cmd.Dir = tmp_dir
	_, err = cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	if parseArgs() {
		os.Exit(1)
	}

	repos := gitea.Repos()
	for _, i := range repos {
		gitBundle(i)
	}
}
