package github

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

func login(accessToken string) (context.Context, *http.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)

	return ctx, oauth2.NewClient(ctx, ts)
}

func GithubClient(accessToken string) *github.Client {
	_, tc := login(accessToken)
	return github.NewClient(tc)
}

func PrintRepositories(client *github.Client) {
	repos, _, err := client.Repositories.List(context.Background(), "", nil)
	checkErr(err)

	for idx, repo := range repos {
		fmt.Printf("%.3d %s\n", idx, repo.GetName())
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
