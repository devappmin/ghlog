package github

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

const SET_REPO_NAME string = "Set Repository Name: "
const SET_REPO_VISIBILITY string = "Set Repository Visibility [default:public/private]: "
const CREATED_REPO string = "Created Repository: "
const PRIVATE string = "private"

const BASE_URL string = "https://github.com/"

type heatmapSets struct {
	date  time.Time
	level int
	count int
}

func login(accessToken string) (context.Context, *http.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)

	return ctx, oauth2.NewClient(ctx, ts)
}

func GithubClient(accessToken string) (*github.Client, context.Context) {
	ctx, tc := login(accessToken)
	return github.NewClient(tc), ctx
}

func PrintRepositories(client *github.Client, ctx context.Context) {
	repos, _, err := client.Repositories.List(ctx, "", nil)
	checkErr(err)

	for idx, repo := range repos {
		fmt.Printf("%.3d %s\n", idx, repo.GetName())
	}
}

func PrintOrganizations(client *github.Client, ctx context.Context) {
	orgs, _, err := client.Organizations.List(ctx, "", nil)
	checkErr(err)

	for idx, org := range orgs {
		fmt.Printf("%.3d %s\n", idx, org.GetLogin())
	}
}

func Search(client *github.Client, ctx context.Context, title string, from int, to int) {
	results, _, err := client.Search.Repositories(ctx, title, nil)
	checkErr(err)

	resultsLen := len(results.Repositories)

	if to >= resultsLen {
		to = resultsLen - 1
	}

	if to == 0 || from > to {
		log.Fatalln("No repositories found")
	}

	for i := from; i <= to; i++ {
		description := results.Repositories[i].GetDescription()
		if len(description) > 50 {
			description = description[:50] + "..."
		}

		fmt.Printf("%.3d────┬ name:        %s\n", i, results.Repositories[i].GetName())
		fmt.Printf("       ├ description: %s\n", description)
		fmt.Printf("       ├ author:      %s\n", results.Repositories[i].GetOwner().GetLogin())
		fmt.Printf("       ├ link:        %s\n", results.Repositories[i].GetHTMLURL())
		fmt.Printf("       ├ star:        %d Stars\n", results.Repositories[i].GetStargazersCount())
		fmt.Printf("       └ language:    %s\n", results.Repositories[i].GetLanguage())
	}
	fmt.Printf("Total %d repositories found, show matches from %d to %d", resultsLen, from, to)
}

func CreateRepository(client *github.Client, ctx context.Context) {
	repoName := getAnswer(SET_REPO_NAME)
	repoVisibility := (strings.ToLower(getAnswer(SET_REPO_VISIBILITY)) == PRIVATE)

	repo := &github.Repository{
		Name:    github.String(repoName),
		Private: github.Bool(repoVisibility),
	}

	newRepo, _, err := client.Repositories.Create(ctx, "", repo)
	checkErr(err)

	fmt.Println(CREATED_REPO, newRepo.GetHTMLURL()+".git")
}

func Heatmap(id string) {
	c := make(chan []heatmapSets)
	go getHeatMap(id, c)

	datasets := <-c

	printHeatmapCalendar(datasets)
}

func getHeatMap(id string, cc chan<- []heatmapSets) {
	datasets := []heatmapSets{}
	c := make(chan heatmapSets)

	userUrl := BASE_URL + "/" + id

	res, err := http.Get(userUrl)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	heatmapContainers := doc.Find(".ContributionCalendar-day")

	heatmapContainers.Each(func(i int, container *goquery.Selection) {
		go extractContainer(container, c)
	})

	for i := 0; i < heatmapContainers.Length(); i++ {
		container := <-c
		if !container.date.IsZero() {
			datasets = append(datasets, container)
		}
	}

	// Sort structs
	sort.Slice(datasets, func(i int, j int) bool {
		if datasets[i].date.After(datasets[j].date) {
			return false
		} else {
			return true
		}
	})

	cc <- datasets
}

func printHeatmapCalendar(datasets []heatmapSets) {
	for i := 0; i < 7; i++ {
		printHeatmapRow(datasets, i)
		fmt.Println()
	}
}

func printHeatmapRow(datasets []heatmapSets, weekday int) {
	for i := weekday; i < len(datasets); i += 7 {
		printHeatmapContainer(datasets[i].count, datasets[i].level)
	}
}

func printHeatmapContainer(count int, level int) {
	var color string

	if level != 0 {
		color = "\033[32m"
	} else {
		color = "\033[0m"
	}

	fmt.Printf("%s %2d", color, count)
}

func extractContainer(container *goquery.Selection, c chan<- heatmapSets) {
	dateStr, _ := container.Attr("data-date")
	date, _ := time.Parse("2006-01-02", dateStr)

	countStr, _ := container.Attr("data-count")
	count, _ := strconv.Atoi(countStr)

	levelStr, _ := container.Attr("data-level")
	level, _ := strconv.Atoi(levelStr)

	c <- heatmapSets{
		date:  date,
		count: count,
		level: level,
	}
}

func getAnswer(question string) string {
	in := bufio.NewReader(os.Stdin)

	fmt.Print(question)
	ans, _ := in.ReadString('\n')
	ans = strings.Replace(strings.TrimSpace(ans), " ", "-", -1)
	return ans
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed:", res.Status)
	}
}
