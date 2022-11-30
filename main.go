package main

import (
	"fmt"
	"github.com/xiaxyi/Search-GithubIssues/cmd"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
