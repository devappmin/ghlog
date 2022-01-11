package github

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
	date  string
	level string
	count string
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

func Heatmap() {
	getHeatMap("devappmin")
}

func getHeatMap(id string) {
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
		if container.date != "" {
			datasets = append(datasets, container)
		}
	}

	for idx, container := range datasets {
		fmt.Println(idx, container.date, container.level, container.count)
	}
}

func extractContainer(container *goquery.Selection, c chan<- heatmapSets) {
	date, _ := container.Attr("data-date")
	count, _ := container.Attr("data-count")
	level, _ := container.Attr("data-level")

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
