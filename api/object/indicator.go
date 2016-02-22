package object

import (
	"time"
)

type Indicators []Indicator

func (i Indicators) Len() int           { return len(i) }
func (i Indicators) Swap(a, b int)      { i[a], i[b] = i[b], i[a] }
func (i Indicators) Less(a, b int) bool { return i[a].Value < i[b].Value }

// Indicator is an int64 value
// evolving in time
type Indicator struct {
	Time  time.Time `json:"time"`
	Value int64     `json:"value"`
}
