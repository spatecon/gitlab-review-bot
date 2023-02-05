package templating

import "strings"

type Locale string

const (
	LocaleRuRu = "ru_ru"
	LocaleEnUs = "en_us"
	LocaleEnEn = "en_en"
)

// ParseLocale parses locale from string and returns Locale type and true if found, or Default locale and false.
func ParseLocale(locale string) (Locale, bool) {
	locale = strings.ReplaceAll(locale, "-", "_")
	locale = strings.ToLower(locale)
	locale = strings.Trim(locale, " ")

	switch locale {
	case LocaleRuRu:
		return LocaleRuRu, true
	case LocaleEnUs: // converts en_US to en_EN
		fallthrough
	case "":
		fallthrough
	case LocaleEnEn:
		return LocaleEnEn, true
	}

	return LocaleEnEn, false
}
