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
}

func (c TodayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	subreddit := vars["subreddit"]
	if len(subreddit) == 0 {
		w.WriteHeader(400)
		return
	}

	data, err := c.getData(subreddit)
	if err != nil {
		log.Println("err:", err.Error())
		w.WriteHeader(500)
		return
	}

	audiences := object.AudiencesFromApp(data)

	buff, err := json.Marshal(todayHandlerResp{
		Audiences: audiences,
		Average:   app.Average(data),
	})
	if err != nil {
		log.Println("err:", err.Error())
		w.WriteHeader(500)
		return
	}
	w.Write(buff)

}

// lowestHighest is a quick implementation retrieving the
// lowest and the highest audience for today.
func (c TodayHandler) lowestHighest(audiences []app.Audience) (app.Audience, app.Audience) {
	var lowest, highest app.Audience
	lowest.Audience = 10E10

	for _, a := range audiences {
		if a.Audience > highest.Audience {
			highest = a
			continue
		}

		if a.Audience < lowest.Audience {
			lowest = a
			continue
		}
	}

	if lowest.Audience == 10E10 {
		lowest.Audience = 0
	}

	return lowest, highest
}

func (c TodayHandler) getData(subreddit string) ([]app.Audience, error) {
	var start, end time.Time

	end = time.Now()
	start = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	return c.App.DB().FindAudiencesInterval(subreddit, start, end)
}
