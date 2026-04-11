package speechpace

import (
	"strings"
	"testing"
)

func TestNewDefaultsToNaiveCounter(t *testing.T) {
	p := New(Config{Language: "en", Preset: Formal})
	if p == nil {
		t.Fatal("New returned nil")
	}
	if p.cfg.SyllableCounter == nil {
		t.Fatal("expected default NaiveCounter, got nil")
	}
}

func TestComputeEmptyTextReturnsZero(t *testing.T) {
	p := New(Config{Language: "en", Preset: Formal})
	r := p.Compute("")
	if r.TotalMS != 0 {
		t.Errorf("expected 0ms for empty text, got %d", r.TotalMS)
	}
}

func TestComputeRealisticShortText(t *testing.T) {
	p := New(Config{Language: "en", Preset: Formal})
	text := `The garden asks for patience, for care, for time. The seasons turn slowly, and the gardener waits.`
	r := p.Compute(text)
	if r.TotalMS < 3000 {
		t.Errorf("expected at least 3000ms for ~20-word speech, got %d", r.TotalMS)
	}
	if r.TotalMS > 20000 {
		t.Errorf("expected at most 20000ms for ~20-word speech, got %d", r.TotalMS)
	}
	if len(r.BySentence) != 2 {
		t.Errorf("expected 2 sentence timings, got %d", len(r.BySentence))
	}
}

func TestComputePauseMarkersAddTime(t *testing.T) {
	p := New(Config{Language: "en", Preset: Formal})
	base := `Hello world. And then the day began.`
	with := `Hello world. [pause 5s] And then the day began.`
	a := p.Compute(base)
	b := p.Compute(with)
	delta := b.TotalMS - a.TotalMS
	if delta < 4500 || delta > 5500 {
		t.Errorf("expected ~5000ms delta from pause marker, got %d", delta)
	}
}

func TestComputePauseMarkerVariations(t *testing.T) {
	cases := []struct {
		text   string
		wantMS int64
	}{
		{`Hello. [pause 2s] World.`, 2000},
		{`Hello. [pause 1.5s] World.`, 1500},
		{`Hello. [pause 0.5s] World.`, 500},
		{`Hello. [pause .5s] World.`, 500},
	}
	p := New(Config{Language: "en", Preset: Formal})
	baseline := p.Compute(`Hello. World.`).TotalMS
	for _, c := range cases {
		t.Run(c.text, func(t *testing.T) {
			got := p.Compute(c.text).TotalMS - baseline
			if got < c.wantMS-100 || got > c.wantMS+100 {
				t.Errorf("want ~%dms from marker, got %d", c.wantMS, got)
			}
		})
	}
}

func TestComputeBreathViolations(t *testing.T) {
	p := New(Config{Language: "en", Preset: Formal})
	// A very long single sentence with many syllables should trip the
	// breath ceiling for English (40 syllables).
	longSentence := strings.Repeat("magnificent ", 20) + "afternoon."
	r := p.Compute(longSentence)
	if len(r.BreathViolations) == 0 {
		t.Errorf("expected breath violation for very long sentence, got none")
	}
}

func TestComputeNoBreathViolationShortSentence(t *testing.T) {
	p := New(Config{Language: "en", Preset: Formal})
	r := p.Compute(`The orchard is wide. The vineyard is narrow.`)
	if len(r.BreathViolations) != 0 {
		t.Errorf("expected no breath violation for short sentences, got %v", r.BreathViolations)
	}
}

func TestComputeResultLanguageNormalized(t *testing.T) {
	cases := []string{"en", "en-US", "en-GB"}
	for _, lang := range cases {
		p := New(Config{Language: lang, Preset: Formal})
		r := p.Compute("The orchard is wide.")
		if r.Language != "en" {
			t.Errorf("language %s normalized to %q, expected en", lang, r.Language)
		}
	}
}
