package speechpace

import "strings"

// rateFor returns the articulation rate in syllables per second for the
// given language and preset. Values are published in languages_table.go
// as a single source of truth for both the rate table and the README
// documentation.
func rateFor(language string, preset Preset) float64 {
	code := normalizeLang(language)
	row, ok := rateTable[code]
	if !ok {
		row = rateTable["en"]
	}
	switch preset {
	case Conversational:
		return row[0]
	case Formal:
		return row[1]
	case Ceremonial:
		return row[2]
	}
	return row[1]
}

// breathCeiling returns the per-language syllable ceiling for a
// one-breath sentence. Sentences exceeding this count are flagged as
// breath violations in the Result struct.
func breathCeiling(language string) int {
	code := normalizeLang(language)
	if v, ok := breathTable[code]; ok {
		return v
	}
	return breathTable["en"]
}

// normalizeLang collapses regional variants into base codes for rate
// and breath-ceiling lookup.
func normalizeLang(language string) string {
	l := strings.ToLower(language)
	switch {
	case strings.HasPrefix(l, "en"):
		return "en"
	case strings.HasPrefix(l, "de"):
		return "de"
	case strings.HasPrefix(l, "es"):
		return "es"
	case strings.HasPrefix(l, "fr"):
		return "fr"
	}
	return "en"
}
