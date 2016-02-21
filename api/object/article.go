package object

import (
	"sort"
	"strings"
	"time"

	"github.com/remeh/reddit-audiences/app"
	"github.com/remeh/reddit-audiences/db"
)

type Article struct {
	ArticleId    string          `json:"id"`
	ArticleTitle string          `json:"title"`
	ArticleLink  string          `json:"link"`
	State        db.ArticleState `json:"state"`
	FirstSeen    time.Time       `json:"first_seen,omitempty"`
	lastSeen     time.Time
	Score        int    `json:"score"`
	Comments     int    `json:"comments"`
	Author       string `json:"author"`
	Promoted     bool   `json:"promoted"`
	Sticky       bool   `json:"sticky"`
	MinRank      int    `json:"min_rank"`
	CurrentRank  int    `json:"current_rank"`
	MaxRank      int    `json:"max_rank"`
	//Ranking      []Ranking `json:"ranking"`
}

type ByRank []Article

func (r ByRank) Len() int           { return len(r) }
func (r ByRank) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ByRank) Less(i, j int) bool { return r[i].CurrentRank < r[j].CurrentRank }

func ArticlesFromApp(articles []db.Article, rankings map[string][]db.Ranking) []Article {
	rv := make([]Article, len(articles))
	for i, a := range articles {
		rv[i] = ArticleFromApp(a, rankings[a.ArticleId])
	}

	// sort by current rank
	byRank := ByRank(rv)
	sort.Sort(&byRank)

	// compute disappeared articles
	// it's done on already sorted by rank articles
	// to be quicker (no need to test other ranks).
	byRank.computeRemoved()

	return byRank
}

func (r ByRank) computeRemoved() {
	for i, tested := range r {
		// ignore sticky and promoted ones.
		if tested.Sticky || tested.Promoted {
			continue
		}

		for _, a := range r {
			if a.CurrentRank > tested.CurrentRank {
				// no need to test with ones having an higher rank
				break
			}

			if tested.CurrentRank == a.CurrentRank {
				if tested.lastSeen.Before(a.lastSeen) {
					tested.State = db.Removed
					r[i] = tested
					break
				}
			}
		}
	}
}

// NOTE(remy): ranking is ORDER BY crawl_time, for first seen
// and last seen, I could simply use [0] and [len()-1].
func ArticleFromApp(article db.Article, ranking []db.Ranking) Article {
	if ranking == nil {
		return Article{}
	}

	var min, max, current int
	var firstSeen time.Time = time.Now()
	var lastSeen time.Time
	min = 10E6

	for _, r := range ranking {
		if r.Rank > max {
			max = r.Rank
		}
		if r.Rank < min {
			min = r.Rank
		}

		if r.CrawlTime.After(lastSeen) {
			current = r.Rank
			lastSeen = r.CrawlTime
		}

		if r.CrawlTime.Before(firstSeen) {
			firstSeen = r.CrawlTime
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
		Score:        article.Score,
		Comments:     article.Comments,
		Author:       article.Author,
		Promoted:     article.Promoted,
		Sticky:       article.Sticky,
		CurrentRank:  current,
		FirstSeen:    firstSeen,
		lastSeen:     lastSeen,
		MinRank:      min,
		MaxRank:      max,
		//Ranking:      ranking,
	}
}
