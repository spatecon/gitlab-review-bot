package templating

import (
	"time"

	en_EN "github.com/spatecon/gitlab-review-bot/pkg/templating/en-EN"
	ru_RU "github.com/spatecon/gitlab-review-bot/pkg/templating/ru-RU"
)

type TimeDiff interface {
	// Since returns the time elapsed since t in human-readable format.
	// For example, "3 days ago" or "just now".
	Since(t time.Time) string
}

type Tools interface {
	TimeDiff

	Plural(n int, wordForms ...string) string
}

func NewTools(locale Locale) Tools {
	switch locale {
	case LocaleRuRu:
		return ru_RU.NewTools()
	default:
		return en_EN.NewTools()
	}
}
