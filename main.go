package main

import (
	"flag"
	"fmt"
	"os"
)

var giteaUrl, apiKey string

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

	giteaUrl = flag.Arg(0)
	apiKey = flag.Arg(1)
	return false
}

func main() {
	if parseArgs() {
		os.Exit(1)
	}

	fmt.Println(giteaUrl, apiKey)
}
