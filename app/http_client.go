package app

import (
	"math/rand"
	"net/http"
)

var useragents []string = []string{
	"Mozilla/5.0 (X11; Linux x86_64; rv:42.0) Gecko/20100101 Firefox/42.0",
	"Mozilla/5.0 (X11; Linux x86_64; rv:43.0) Gecko/20100101 Firefox/43.0",
	"Mozilla/5.0 (Linux; Android 5.1; XT1039 Build/LPB23.13-17.6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.95 Mobile Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.103 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.103 Safari/537.36",
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
