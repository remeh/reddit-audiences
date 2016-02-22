// Reddit audiences crawler
// Rémy Mathieu © 2016
package object

import (
	"time"
)

type Indicators []Indicator

func (i Indicators) Len() int           { return len(i) }
func (i Indicators) Swap(a, b int)      { i[a], i[b] = i[b], i[a] }
func (i Indicators) Less(a, b int) bool { return i[a].Time.Before(i[b].Time) }

// Indicator is an int value
// evolving in time
type Indicator struct {
	Time  time.Time `json:"time"`
	Value int       `json:"value"`
}
