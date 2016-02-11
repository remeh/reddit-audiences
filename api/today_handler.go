// Reddit audiences crawler
// Rémy Mathieu © 2016
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/remeh/reddit-audiences/api/object"
	"github.com/remeh/reddit-audiences/app"

	"github.com/gorilla/mux"
)

type TodayHandler struct {
	App *app.App
}

type todayHandlerResp struct {
	Audiences       []object.Audience `json:"audiences"`
	Average         int64             `json:"average"`
	LowestAudience  object.Audience   `json:"lowest_audience"`
	HighestAudience object.Audience   `json:"highest_audience"`
	Articles        []object.Article  `json:"articles"`
}

func (c TodayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subreddit := vars["subreddit"]
	if len(subreddit) == 0 {
		w.WriteHeader(400)
		return
	}

	dataAudiences, dataRankings, dataArticles, err := c.getData(subreddit)
	if err != nil {
		log.Println("err:", err.Error())
		w.WriteHeader(500)
		return
	}

	audiences := object.AudiencesFromApp(dataAudiences)
	lowest, highest := app.LowestHighest(dataAudiences)
	articles := object.ArticlesFromApp(dataArticles, dataRankings)

	buff, err := json.Marshal(todayHandlerResp{
		Audiences:       audiences,
		Average:         app.Average(dataAudiences),
		Articles:        articles,
		LowestAudience:  object.AudienceFromApp(lowest),
		HighestAudience: object.AudienceFromApp(highest),
	})
	if err != nil {
		log.Println("err:", err.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(buff)

}

func (c TodayHandler) getData(subreddit string) ([]app.Audience, map[string][]app.Ranking, []app.Article, error) {
	var start, end time.Time

	end = time.Now()
	start = time.Now().Add(-time.Hour * 36)

	audiences, err := c.App.DB().FindAudiencesInterval(subreddit, start, end)
	if err != nil {
		return nil, nil, nil, err
	}

	rankings, err := c.App.DB().FindArticlesRanking(subreddit, start, end)
	if err != nil {
		return nil, nil, nil, err
	}

	articles, err := c.App.DB().FindArticles(subreddit, start, end)
	if err != nil {
		return nil, nil, nil, err
	}

	return audiences, rankings, articles, nil
}
