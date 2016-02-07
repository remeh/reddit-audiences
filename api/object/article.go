package object

import (
	"strings"

	"github.com/remeh/reddit-audiences/app"
)

type Article struct {
	ArticleId    string `json:"article_id"`
	ArticleTitle string `json:"article_title"`
	ArticleLink  string `json:"article_link"`
	Author       string `json:"author"`
	Promoted     bool   `json:"promoted"`
	Sticky       bool   `json:"sticky"`
	MinRank      int    `json:"min_rank"`
	MaxRank      int    `json:"max_rank"`
	//Ranking      []Ranking `json:"ranking"`
}

func ArticlesFromApp(articles []app.Article, rankings Rankings) []Article {
	rv := make([]Article, len(articles))
	for i, a := range articles {
		rv[i] = ArticleFromApp(a, rankings[a.ArticleId])
	}
	return rv
}

func ArticleFromApp(article app.Article, ranking []Ranking) Article {
	var min, max int
	min = 10E6

	for _, r := range ranking {
		if r.Rank > max {
			max = r.Rank
		}
		if r.Rank < min {
			min = r.Rank
		}
	}

	// rebuild the http link for self posts
	link := article.ArticleLink
	if strings.HasPrefix(link, "/r/") {
		link = "https://reddit.com" + link
	}

	return Article{
		ArticleId:    article.ArticleId,
		ArticleTitle: article.ArticleTitle,
		ArticleLink:  link,
		Author:       article.Author,
		Promoted:     article.Promoted,
		Sticky:       article.Sticky,
		MinRank:      min,
		MaxRank:      max,
		//Ranking:      ranking,
	}
}
