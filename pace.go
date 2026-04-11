// Package speechpace computes speech delivery time from text.
//
// Given a source text and a language/preset configuration, the package
// returns the expected spoken delivery time in milliseconds, plus a
// per-sentence breakdown and a list of breath-limit violations. The
// calculation is deterministic: syllables times seconds-per-syllable
// plus punctuation pauses plus inline pause markers.
//
// speechpace has no network dependencies and no machine learning
// components. Syllable counting is delegated to a pluggable counter
// interface; the package ships a naive vowel-group counter as its
// default, and CivNode wires a speedata/hyphenation-backed counter at
// the application layer for higher accuracy.
package speechpace

// Preset identifies a delivery cadence. Per-language seconds-per-syllable
// rates are defined in languages.go and calibrated against published
// linguistics research on comparable articulation rates.
type Preset int

const (
	// Conversational is the natural rate of unhurried speech between
	// two people in a quiet room.
	Conversational Preset = iota
	// Formal is the podium-address rate used for keynotes, business
	// presentations, and most professional delivery.
	Formal
	// Ceremonial is the slowest preset, used for eulogies, weddings,
	// religious speech, and other occasions where gravity calls for
	// unhurried pacing.
	Ceremonial
)

// Config configures a pace calculator.
type Config struct {
	// Language is one of "en", "en-US", "en-GB", "de", "es", "fr".
	// Unknown values fall back to English.
	Language string
	// Preset selects the cadence (conversational, formal, ceremonial).
	Preset Preset
	// SyllableCounter is optional; if nil the package uses NaiveCounter.
	// Callers can plug in a Liang-pattern counter or any other
	// implementation of the interface for better accuracy.
	SyllableCounter SyllableCounter
}

// SyllableCounter counts syllables in text for a given language.
// Implementations must be safe for concurrent use across multiple
// Compute calls from the same or different Pace instances.
type SyllableCounter interface {
	CountSyllables(text string, language string) int
}

// Pace is a configured calculator produced by New.
type Pace struct {
	cfg Config
}

// New returns a Pace calculator configured for the given language and
// preset. If SyllableCounter is nil it falls back to NaiveCounter.
func New(cfg Config) *Pace {
	if cfg.SyllableCounter == nil {
		cfg.SyllableCounter = NaiveCounter{}
	}
	return &Pace{cfg: cfg}
}

// Result is the output of Compute. TotalMS is the full delivery time
// including punctuation pauses and inline pause markers. BySentence
// gives per-sentence breakdowns for ruler rendering and spatial
// visualization. BreathViolations lists 0-based sentence indices that
// exceed the per-language breath ceiling.
type Result struct {
	TotalMS          int64
	BySentence       []SentenceTiming
	BreathViolations []int
	Language         string
}

// SentenceTiming is the per-sentence breakdown used by editor rulers
// and waveform visualizations.
type SentenceTiming struct {
	Index     int   // 0-based sentence index in the source
	Syllables int   // content syllables in the sentence
	SpeakMS   int64 // speaking time only (syllables times rate)
	PauseMS   int64 // punctuation and inline-marker pauses charged to this sentence
}

// NaiveCounter is a minimal fallback syllable counter that approximates
// counts by counting vowel groups. Good enough for tests and a safe
// default; production systems should plug in a Liang-pattern counter
// via the SyllableCounter interface for significantly better accuracy.
type NaiveCounter struct{}

// CountSyllables returns an approximate syllable count using the
// vowel-group heuristic. It treats accented vowels as vowels so the
// same counter works across English, German, Spanish, and French.
func (NaiveCounter) CountSyllables(text string, language string) int {
	count := 0
	inVowel := false
	for _, r := range text {
		if isVowelRune(r) {
			if !inVowel {
				count++
			}
			inVowel = true
		} else {
			inVowel = false
		}
	}
	return count
}

func isVowelRune(r rune) bool {
	switch r {
	case 'a', 'e', 'i', 'o', 'u', 'y',
		'A', 'E', 'I', 'O', 'U', 'Y',
		'ä', 'ö', 'ü', 'Ä', 'Ö', 'Ü',
		'á', 'é', 'í', 'ó', 'ú', 'ý',
		'Á', 'É', 'Í', 'Ó', 'Ú', 'Ý',
		'à', 'è', 'ì', 'ò', 'ù',
		'À', 'È', 'Ì', 'Ò', 'Ù',
		'â', 'ê', 'î', 'ô', 'û',
		'Â', 'Ê', 'Î', 'Ô', 'Û',
		'ã', 'õ', 'ñ',
		'Ã', 'Õ', 'Ñ':
		return true
	}
	return false
}
