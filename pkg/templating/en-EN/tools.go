package en_EN

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gertd/go-pluralize"
	"github.com/spatecon/gitlab-review-bot/pkg/motivational"
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

func (t *Tools) Motivation() string {
	// pick random motivational phrase from motivational.ReviewMotivationEnPhrases
	i := rand.Int() % len(motivational.ReviewMotivationEnPhrases) //nolint:gosec

	return motivational.ReviewMotivationEnPhrases[i]
}
