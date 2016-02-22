// Reddit audiences crawler
// Rémy Mathieu © 2016
package api

import (
	"log"
	"net/http"
	"time"

	"github.com/remeh/reddit-audiences/api/object"
	"github.com/remeh/reddit-audiences/app"
	"github.com/remeh/reddit-audiences/db"

	"github.com/gorilla/mux"
)

type TodayHandler struct {
	App *app.App
}

type todayHandlerResp struct {
	Audiences       []object.Audience   `json:"audiences"`
	Average         int64               `json:"average"`
	LowestAudience  object.Audience     `json:"lowest_audience"`
	HighestAudience object.Audience     `json:"highest_audience"`
	Articles        []object.Article    `json:"articles"`
	Annotations     []object.Annotation `json:"annotations"`
	DemoModeMessage bool                `json:"demo_mode_message"`
}

func (c TodayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subreddit := vars["subreddit"]
	if len(subreddit) == 0 {
		w.WriteHeader(400)
		return
	}

	r.ParseForm()
	t := r.Form.Get("t")
	hours := 36

	// ensure the right to retrieve more than 36h
	// ----------------------
	tmplParams := app.TmplParams(c.App, r, "ApiToday")

	demoModeMessage := false
	if len(t) > 0 && t != "36h" && len(tmplParams.User.Email) == 0 {
		// not auth user, demo message and stay on 36h
		demoModeMessage = true
	} else {
		// auth user, test if he wants 7 days of data
		if t == "7d" {
			hours = 24 * 7
		}
	}

	// retrieve data
	// ----------------------

	dataAudiences, dataRankings, dataArticles, err := c.getData(subreddit, hours)
	if err != nil {
		log.Println("err:", err.Error())
		w.WriteHeader(500)
		return
	}

	audiences := object.AudiencesFromApp(dataAudiences)
	lowest, highest := app.LowestHighest(dataAudiences)
	articles := object.ArticlesFromApp(dataArticles, dataRankings)

	// If the user is logged in, we need to
	// retrieve its annotations.
	// ----------------------
	user := app.GetUser(c.App.DB(), r)

	annotations := make([]object.Annotation, 0)
	if len(user.Email) > 0 {
		after := time.Now().Add(-time.Hour * time.Duration(hours))
		values, err := c.App.DB().FindAnnotations(subreddit, user.Uuid, after)
		if err != nil {
			log.Println("err:", err.Error())
			w.WriteHeader(500)
			return
		}
		for _, a := range values {
			annotations = append(annotations, object.Annotation{
				Message: a.Message,
				Time:    a.Time,
			})
		}
	}

	// serialize and send response
	// ----------------------

	render(w, 200, todayHandlerResp{
		Audiences:       audiences,
		Average:         app.Average(dataAudiences),
		Articles:        articles,
		DemoModeMessage: demoModeMessage,
		Annotations:     annotations,
		LowestAudience:  object.AudienceFromApp(lowest),
		HighestAudience: object.AudienceFromApp(highest),
	})
}

func (c TodayHandler) getData(subreddit string, hours int) ([]db.Audience, map[string][]db.Ranking, []db.Article, error) {
	var start, end time.Time

	end = time.Now()
	start = time.Now().Add(-time.Hour * time.Duration(hours))

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
