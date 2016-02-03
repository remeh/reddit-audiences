package app

// Average computes the average in the given audiences.
// Very basic algorithm.
func Average(audiences []Audience) int64 {
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
func LowestHighest(audiences []Audience) (Audience, Audience) {
	var lowest, highest Audience
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
