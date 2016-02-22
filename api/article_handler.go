// Reddit audiences crawler
// Rémy Mathieu © 2016
package api

import (
	"net/http"
	"time"

	"github.com/remeh/reddit-audiences/api/object"
	"github.com/remeh/reddit-audiences/app"
)

type ArticleHandler struct {
	App *app.App
}

type articleHandlerResp struct {
	// time at which this articles appeared on the front page
	Appearance time.Time `json:"appearance"`
	// time at which the article finally exited from the front page
	Disappearance time.Time `json:"disappearance"`

	CurrentScore    int `json:"current_score"`
	CurrentRank     int `json:"current_rank"`
	CurrentComments int `json:"current_comments"`

	Scores   object.Indicators `json:"scores"`
	Ranks    object.Indicators `json:"ranks"`
	Comments object.Indicators `json:"comments"`

	DemoModeMessage bool `json:"demo_mode_message"`
}

func (c ArticleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// TODO(remy):
}
