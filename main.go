package main

import (
	"fmt"
	"ghIssueCLI/cli"
	"os"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
