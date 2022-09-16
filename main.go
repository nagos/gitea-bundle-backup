package main

import (
	"flag"
	"fmt"
	"os"
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

func main() {
	if parseArgs() {
		os.Exit(1)
	}

	fmt.Println(gitea.Repos())
}
