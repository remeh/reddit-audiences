// Reddit audiences crawler
// Rémy Mathieu © 2016
package api

import (
	"net/http"
	"sort"
	"time"

	"github.com/remeh/reddit-audiences/api/object"
	"github.com/remeh/reddit-audiences/app"
	"github.com/remeh/reddit-audiences/db"

	"github.com/gorilla/mux"
)

type ArticleHandler struct {
	App *app.App
}

type articleHandlerResp struct {
	// time at which this articles appeared on the front page
	FirstSeen time.Time `json:"first_seen"`
	// time at which the article finally exited from the front page
	LastSeen time.Time `json:"last_seen"`

	CurrentScore    int `json:"current_score"`
	CurrentRank     int `json:"current_rank"`
	CurrentComments int `json:"current_comments"`

	Scores   object.Indicators `json:"scores"`
	Ranks    object.Indicators `json:"ranks"`
	Comments object.Indicators `json:"comments"`

	DemoModeMessage bool `json:"demo_mode_message"`
}

func (c ArticleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subreddit := vars["subreddit"]
	if len(subreddit) == 0 {
		w.WriteHeader(400)
		return
	}
	articleId := vars["articleId"]
	if len(articleId) == 0 {
		w.WriteHeader(400)
		return
	}

	// no extra article information in demo mode
	// ----------------------
	user := app.GetUser(c.App.DB(), r)

	if len(user.Email) == 0 {
		// demo mode
		render(w, 200, articleHandlerResp{
			DemoModeMessage: true,
		})
		return
	}

	// get the data
	// ----------------------

	rankings, err := c.App.DB().FindArticleRanking(subreddit, articleId)
	if err != nil {
		render(w, 500, ErrorResponse{
			Code:    500,
			Message: "can't retrieve data",
		})
		return
	}

	resp := articleHandlerResp{}
	computeRankings(&resp, rankings)
	render(w, 200, resp)
}

func computeRankings(r *articleHandlerResp, rankings db.Rankings) {
	var currentScore, currentRank, currentComments int
	scores, ranks, comments := make(object.Indicators, len(rankings)), make(object.Indicators, len(rankings)), make(object.Indicators, len(rankings))

	var lastSeen time.Time
	var firstSeen time.Time = time.Now()

	// compute on sorted rankings
	sort.Sort(&rankings)

	for i, ranking := range rankings {
		if ranking.CrawlTime.Before(firstSeen) {
			firstSeen = ranking.CrawlTime
		}
		if ranking.CrawlTime.After(lastSeen) {
			lastSeen = ranking.CrawlTime
			currentScore = ranking.Score
			currentComments = ranking.Comments
			currentRank = ranking.Rank
		}

		scores[i] = object.Indicator{
			Time:  ranking.CrawlTime,
			Value: ranking.Score,
		}

		comments[i] = object.Indicator{
			Time:  ranking.CrawlTime,
			Value: ranking.Comments,
		}

		ranks[i] = object.Indicator{
			Time:  ranking.CrawlTime,
			Value: ranking.Rank,
		}
	}

	r.CurrentComments = currentComments
	r.CurrentRank = currentRank
	r.CurrentScore = currentScore
	r.Comments = comments
	r.Ranks = ranks
	r.Scores = scores
	r.FirstSeen = firstSeen
	r.LastSeen = lastSeen
}
