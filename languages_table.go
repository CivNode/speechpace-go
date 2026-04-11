package speechpace

// rateTable stores per-language articulation rates in syllables per
// second, in the order [Conversational, Formal, Ceremonial].
//
// Rates are calibrated against comparable articulation-rate research
// (the 2011 Université de Lyon study on syllable rate across languages
// and related linguistics literature). Conversational rates approximate
// natural speech between two people; formal rates approximate podium
// delivery; ceremonial rates approximate unhurried occasions like
// eulogies and weddings where gravity slows the pace further.
//
// Adjusting these values recalibrates every Pace instance globally.
// User overrides should be applied at the application layer by
// subtracting or adding a per-user offset, not by editing this table.
var rateTable = map[string][3]float64{
	// [conversational, formal, ceremonial]
	"de": {4.8, 3.8, 3.2},
	"en": {5.5, 4.5, 3.8},
	"fr": {6.0, 5.0, 4.2},
	"es": {6.5, 5.5, 4.6},
}

// breathTable stores per-language syllable ceilings for one-breath
// sentences. Speakers can comfortably deliver about this many
// syllables in a single breath at formal cadence. Sentences exceeding
// the ceiling are flagged as breath violations.
//
// Values scale with articulation rate: languages that speak faster
// can fit more syllables in a breath. German sits lowest because
// compound words and heavier consonant clusters shorten the effective
// breath capacity.
var breathTable = map[string]int{
	"de": 35,
	"en": 40,
	"fr": 45,
	"es": 50,
}
