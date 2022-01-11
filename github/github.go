package github

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
)

const SET_REPO_NAME string = "Set Repository Name: "
const SET_REPO_VISIBILITY string = "Set Repository Visibility [public/private]: "
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

func CreateRepository(client *github.Client) {
	repoName := getAnswer(SET_REPO_NAME)
	repoVisibility := (getAnswer(SET_REPO_VISIBILITY) == PRIVATE)

	repo := &github.Repository{
		Name:    github.String(repoName),
		Private: github.Bool(repoVisibility),
	}

	newRepo, _, err := client.Repositories.Create(context.Background(), "", repo)
	checkErr(err)

	fmt.Println(CREATED_REPO, newRepo.GetGitURL())
}

func Heatmap(id string) {
	datasets := getHeatMap(id)
	printHeatmapCalendar(datasets)
}

func getHeatMap(id string) []heatmapSets {
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

	return datasets
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

func getAnswer(question string) (ans string) {
	fmt.Print(question)
	fmt.Scanln(&ans)

	return
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
