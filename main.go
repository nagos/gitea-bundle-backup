package main

import (
	"errors"
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

func parseArgs() error {
	flag.Usage = myUsage
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		return errors.New("bad arguments")
	}

	gitea.SetUrl(flag.Arg(0))
	gitea.SetApiKey(flag.Arg(1))
	return nil
}

func shellRun(cmd string, dir string, args ...string) error {
	c := exec.Command(cmd, args...)
	if len(dir) != 0 {
		c.Dir = dir
	}
	_, err := c.Output()

	return err
}

func gitBundle(repo string) {
	url := strings.Split(repo, "/")
	repo_name := url[len(url)-1]
	repo_user := url[len(url)-2]
	fmt.Printf("Bundling %s/%s\n", repo_user, repo_name)
	tmp_dir, err := ioutil.TempDir("", "gitea_backup")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmp_dir)

	err = shellRun("git", "", "clone", "--mirror", repo, tmp_dir)
	if err != nil {
		log.Println(err.Error())
		return
	}

	cwd, _ := os.Getwd()
	err = shellRun("git", cwd, "bundle", "create", fmt.Sprintf("%s/%s_%s", cwd, repo_user, repo_name), "--all")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func main() {
	if e := parseArgs(); e != nil {
		log.Fatal("Failed to run program: ", e.Error())
	}

	repos, e := gitea.Repos()
	if e != nil {
		log.Fatal("Cannot bundle repos: ", e.Error())
	}
	for _, i := range repos {
		gitBundle(i)
	}
}
