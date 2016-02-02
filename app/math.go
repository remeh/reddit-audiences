package app

// Average computes the average in the given audiences.
// Very basic algorithm.
func Average(audiences []Audience) int64 {
	var rv int64
	for _, a := range audiences {
		rv += int64(a.Audience)
	}
	return rv / int64(len(audiences))
}
