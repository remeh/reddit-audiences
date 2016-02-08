package object

import (
	"sort"
	"strings"
	"time"

	"github.com/remeh/reddit-audiences/app"
)

type Article struct {
	ArticleId    string           `json:"id"`
	ArticleTitle string           `json:"title"`
	ArticleLink  string           `json:"link"`
	State        app.ArticleState `json:"state"`
	Author       string           `json:"author"`
	Promoted     bool             `json:"promoted"`
	Sticky       bool             `json:"sticky"`
	MinRank      int              `json:"min_rank"`
	CurrentRank  int              `json:"current_rank"`
	MaxRank      int              `json:"max_rank"`
	//Ranking      []Ranking `json:"ranking"`
}

type ByRank []Article

func (r ByRank) Len() int           { return len(r) }
func (r ByRank) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByRank) Less(i, j int) bool { return r[i].CurrentRank < r[j].CurrentRank }

func ArticlesFromApp(articles []app.Article, rankings map[string][]app.Ranking) []Article {
	rv := make([]Article, len(articles))
	for i, a := range articles {
		rv[i] = ArticleFromApp(a, rankings[a.ArticleId])
	}

	byRank := ByRank(rv)
	sort.Sort(&byRank)
	return byRank
}

func ArticleFromApp(article app.Article, ranking []app.Ranking) Article {
	if ranking == nil {
		return Article{}
	}

	var min, max, current int
	var latest time.Time
	min = 10E6

	for _, r := range ranking {
		if r.Rank > max {
			max = r.Rank
		}
		if r.Rank < min {
			min = r.Rank
		}

		if r.CrawlTime.After(latest) {
			current = r.Rank
		}
	}

	// rebuild the http link for self posts
	link := article.ArticleLink
	if strings.HasPrefix(link, "/r/") {
		link = "https://reddit.com" + link
	}

	state := app.ComputeArticleState(article, ranking)

	return Article{
		ArticleId:    article.ArticleId,
		ArticleTitle: article.ArticleTitle,
		ArticleLink:  link,
		State:        state,
		Author:       article.Author,
		Promoted:     article.Promoted,
		Sticky:       article.Sticky,
		CurrentRank:  current,
		MinRank:      min,
		MaxRank:      max,
		//Ranking:      ranking,
	}
}
