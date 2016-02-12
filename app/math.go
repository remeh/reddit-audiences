package app

import (
	"github.com/remeh/reddit-audiences/db"
)

// Average computes the average in the given audiences.
// Very basic algorithm.
func Average(audiences []db.Audience) int64 {
	var rv int64
	for _, a := range audiences {
		rv += int64(a.Audience)
	}

	if len(audiences) == 0 {
		return 0
	}

	return rv / int64(len(audiences))
}

// LowestHighest is a quick implementation retrieving the
// lowest and the highest audience for the given slice of audiences.
func LowestHighest(audiences []db.Audience) (db.Audience, db.Audience) {
	var lowest, highest db.Audience
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

func ComputeArticleState(article db.Article, ranking []db.Ranking) db.ArticleState {
	if len(ranking) <= 2 {
		return db.New
	}

	increasing, decreasing := 0, 0

	prev := ranking[0]
	for i, r := range ranking {
		if i == 0 {
			continue
		}

		if prev.Rank > r.Rank {
			increasing += 1
		} else if prev.Rank < r.Rank {
			decreasing += 1
		}

		prev = r
	}

	if increasing > decreasing {
		return db.Rising
	} else if decreasing > increasing {
		return db.Falling
	}

	return db.Stagnant
}
