package main

import (
	"fmt"
	"os"

	"main.com/auth"
	"main.com/github"
)

// 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

const HELP_STRING string = `
usage: gh <command>

Here are few commands <3:
	repo			Print all repositories on the github
	create repo		Create a new repository
`

func run(args []string) {

	if len(args) == 0 {
		args = append(args, "help")
	}

	if args[0] == "help" {
		printHelp()
		return
	}

	accessToken, _ := auth.GetAuth()
	client := github.GithubClient(accessToken)

	switch args[0] {
	case "repo":
		github.PrintRepositories(client)

	case "create":
		if args[1] == "repo" {
			github.CreateRepository(client)
		}

	case "heatmap":
		github.Heatmap()
	}
}

func printHelp() {
	fmt.Println(HELP_STRING)
}

func main() {

	run(os.Args[1:])
}
