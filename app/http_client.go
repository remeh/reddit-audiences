package app

import (
	"math/rand"
	"net/http"
)

var useragents []string = []string{
	"Mozilla/5.0 (X11; Linux x86_64; rv:42.0) Gecko/20100101 Firefox/42.0",
	"Mozilla/5.0 (X11; Linux x86_64; rv:43.0) Gecko/20100101 Firefox/43.0",
}

func NewRequest(url string) (*http.Request, error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	ua := randomUseragent()
	r.Header.Set("User-Agent", ua)
	return r, nil
}

func randomUseragent() string {
	return useragents[rand.Int()%len(useragents)]
}
