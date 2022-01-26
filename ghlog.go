package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"main.com/auth"
	"main.com/github"
)

const GHLOG_VERSION string = "v0.1.0"
const HELP_STRING string = `
usage: ghlog <command>

Here are few commands:
	repo                Print all repositories on your github
	org                 Print all organizations you joined
	create repo         Create a new repository
	heatmap	[user]      Print all contributions as heatmap calendar

You can also search repositories using:
	search <query> [from] [to]          Search repositories from [from] to [to]
	                                    Default value
	                                    [from] = 0
	                                    [to]   = 5`

func run(args []string) {

	if len(args) == 0 {
		args = append(args, "help")
	}

	if args[0] == "help" {
		printHelp()
		return
	}

	accessToken, id := auth.GetAuth()
	client, ctx := github.GithubClient(accessToken)

	switch args[0] {
	case "repo":
		github.PrintRepositories(client, ctx)

	case "org":
		github.PrintOrganizations(client, ctx)

	case "create":
		if len(args) < 2 {
			log.Fatalln(errors.New("Can't find command."))
		}
		if args[1] == "repo" {
			github.CreateRepository(client, ctx)
		}

	case "heatmap":
		if len(args) < 2 {
			args = append(args, id)
		}
		github.Heatmap(args[1])

	case "search":
		if len(args) < 2 {
			log.Fatalln(errors.New("Missing argument."))
		}

		from := 0
		to := 5

		if len(args) >= 3 {
			from, _ = strconv.Atoi(args[2])
		}

		if len(args) >= 4 {
			to, _ = strconv.Atoi(args[3])
		}

		github.Search(client, ctx, args[1], from, to)
	}
}

func printHelp() {
	fmt.Printf("Unofficial Github Cli v%s\n", GHLOG_VERSION)
	fmt.Println(HELP_STRING)
}

func main() {
	run(os.Args[1:])
}
