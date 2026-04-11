package speechpace

import "testing"

func TestPauseWeightsGermanHeavierOnComma(t *testing.T) {
	de := pauseWeights("de")
	es := pauseWeights("es")
	if de.Comma <= es.Comma {
		t.Errorf("German comma pause (%d) should be heavier than Spanish (%d)", de.Comma, es.Comma)
	}
}

func TestPauseWeightsPeriodHeavierThanComma(t *testing.T) {
	for _, lang := range []string{"en", "de", "es", "fr"} {
		w := pauseWeights(lang)
		if w.Period <= w.Comma {
			t.Errorf("%s: period pause %d should exceed comma %d", lang, w.Period, w.Comma)
		}
	}
}

func TestPunctuationPauseMSCommaAndPeriod(t *testing.T) {
	w := pauseWeights("en")
	text := "The garden asks for patience, for care, for time."
	got := punctuationPauseMS(text, w)
	// Two commas + one period = 2*180 + 400 = 760
	want := int64(2*w.Comma + w.Period)
	if got != want {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestPunctuationPauseMSQuestionAndExclaim(t *testing.T) {
	w := pauseWeights("en")
	text := "Is this the orchard? Yes, it is! Truly."
	got := punctuationPauseMS(text, w)
	// 1 question + 1 comma + 1 exclaim + 1 period
	want := int64(w.Question + w.Comma + w.Exclaim + w.Period)
	if got != want {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestPunctuationPauseMSEmDash(t *testing.T) {
	w := pauseWeights("en")
	text := "The orchard — always the orchard — waits."
	got := punctuationPauseMS(text, w)
	want := int64(2*w.EmDash + w.Period)
	if got != want {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestPunctuationPauseMSEmpty(t *testing.T) {
	w := pauseWeights("en")
	if got := punctuationPauseMS("", w); got != 0 {
		t.Errorf("want 0 for empty string, got %d", got)
	}
}

func TestPauseWeightsUnknownLanguageFallsBackToEnglish(t *testing.T) {
	en := pauseWeights("en")
	other := pauseWeights("klingon")
	if en != other {
		t.Errorf("unknown language should fall back to English weights")
	}
}
