package speechpace

import "testing"

func TestRateForGermanSlowerThanSpanish(t *testing.T) {
	de := rateFor("de", Formal)
	es := rateFor("es", Formal)
	if de >= es {
		t.Errorf("expected German formal rate (%f) slower than Spanish (%f)", de, es)
	}
}

func TestRateForConversationalExceedsFormalExceedsCeremonial(t *testing.T) {
	langs := []string{"en", "de", "es", "fr"}
	for _, lang := range langs {
		t.Run(lang, func(t *testing.T) {
			conv := rateFor(lang, Conversational)
			form := rateFor(lang, Formal)
			cerem := rateFor(lang, Ceremonial)
			if !(conv > form && form > cerem) {
				t.Errorf("expected conv > form > ceremonial, got %f, %f, %f", conv, form, cerem)
			}
		})
	}
}

func TestRateForUnknownLanguageFallsBackToEnglish(t *testing.T) {
	r := rateFor("klingon", Formal)
	enRate := rateFor("en", Formal)
	if r != enRate {
		t.Errorf("unknown language should fall back to English rate %f, got %f", enRate, r)
	}
}

func TestRateForRegionalEnglishVariants(t *testing.T) {
	r1 := rateFor("en", Formal)
	r2 := rateFor("en-US", Formal)
	r3 := rateFor("en-GB", Formal)
	if r1 != r2 || r2 != r3 {
		t.Errorf("regional en variants should share the same rate, got %f, %f, %f", r1, r2, r3)
	}
}

func TestBreathCeilingGermanStricterThanSpanish(t *testing.T) {
	if breathCeiling("de") >= breathCeiling("es") {
		t.Error("German breath ceiling should be stricter than Spanish")
	}
}

func TestBreathCeilingKnownLanguages(t *testing.T) {
	cases := map[string]int{"de": 35, "en": 40, "fr": 45, "es": 50}
	for lang, want := range cases {
		if got := breathCeiling(lang); got != want {
			t.Errorf("breath ceiling for %s: want %d, got %d", lang, want, got)
		}
	}
}

func TestNormalizeLang(t *testing.T) {
	cases := map[string]string{
		"en":    "en",
		"en-US": "en",
		"en-GB": "en",
		"EN-US": "en",
		"de":    "de",
		"de-AT": "de",
		"es":    "es",
		"es-MX": "es",
		"fr":    "fr",
		"fr-CA": "fr",
		"foo":   "en",
		"":      "en",
	}
	for in, want := range cases {
		if got := normalizeLang(in); got != want {
			t.Errorf("normalizeLang(%q): want %q, got %q", in, want, got)
		}
	}
}
