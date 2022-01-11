package main

import (
	"os"

	"main.com/auth"
	"main.com/github"
)

// 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

func main() {
	accessToken, _ := auth.GetAuth()
	client := github.GithubClient(accessToken)

	if len(os.Args) == 1 {
		return
	}

	if os.Args[1] == "repo" {
		github.PrintRepositories(client)
	}
}
