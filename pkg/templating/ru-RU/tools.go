package ru_RU

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spatecon/gitlab-review-bot/pkg/motivational"
)

type Tools struct {
}

func NewTools() *Tools {
	return &Tools{}
}

func (t *Tools) Plural(n int, wordForms ...string) string {
	forms := []string{"", "", ""}

	if len(wordForms) == 1 {
		forms[0] = wordForms[0]
		forms[1] = wordForms[0]
		forms[2] = wordForms[0]
	}

	if len(wordForms) == 3 {
		forms = wordForms
	}

	if n >= 11 && n <= 19 {
		return forms[2]
	}

	n = n % 10

	if n == 1 {
		return forms[0]
	}

	if n >= 2 && n <= 4 {
		return forms[1]
	}

	return forms[2]
}

func (t *Tools) Since(tm time.Time) string {
	dur := time.Since(tm)

	if dur.Hours() < 2 {
		return "совсем недавно"
	}

	if dur.Hours()/24 > 14 {
		return "больше двух недель назад"
	}

	var (
		val    int
		metric string
	)

	dur = dur.Truncate(time.Minute)
	if dur.Hours() > 24 {
		val = int(dur.Hours() / 24)
		metric = t.Plural(val, "день", "дня", "дней")
	} else {
		val = int(dur.Hours())
		metric = t.Plural(val, "час", "часа", "часов")
	}

	return fmt.Sprintf("%d %s назад", val, metric)
}

func (t *Tools) Motivation() string {
	// pick random motivational phrase from motivational.ReviewMotivationEnPhrases
	i := rand.Int() % len(motivational.ReviewMotivationEnPhrases)

	return motivational.ReviewMotivationEnPhrases[i]
}
