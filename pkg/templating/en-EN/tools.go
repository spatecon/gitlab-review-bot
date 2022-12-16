package en_EN

import (
	"fmt"
	"time"

	"github.com/gertd/go-pluralize"
)

type Tools struct {
	pluralizer *pluralize.Client
}

func NewTools() *Tools {
	return &Tools{
		pluralizer: pluralize.NewClient(),
	}
}

func (t *Tools) Plural(n int, wordForms ...string) string {
	word := ""
	if len(wordForms) > 0 {
		word = wordForms[0]
	}

	if n == 1 {
		return t.pluralizer.Singular(word)
	}

	return t.pluralizer.Plural(word)
}

func (t *Tools) Since(tm time.Time) string {
	dur := time.Since(tm)

	if dur.Hours() < 2 {
		return "just now"
	}

	if dur.Hours()/24 > 14 {
		return "more than 2 weeks ago"
	}

	var (
		val    int
		metric string
	)

	dur = dur.Truncate(time.Minute)
	if dur.Hours() > 24 {
		val = int(dur.Hours() / 24)
		metric = t.Plural(val, "day")
	} else {
		val = int(dur.Hours())
		metric = t.Plural(val, "hour")
	}

	return fmt.Sprintf("%d %s ago", val, metric)
}
