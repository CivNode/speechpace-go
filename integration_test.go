package speechpace

import (
	"strings"
	"testing"
)

// A purpose-built peaceful-themed speech of approximately 400 words used
// as the integration fixture. It contains commas, periods, question
// marks, em dashes, and pause markers so every code path in the
// orchestrator is exercised.
const demonstrationSpeech = `Friends, thank you for joining us in the garden this morning. The roses have just begun to open, the herbs are high enough to brush against the path, and the bees are already at their day's work.

I want to tell you a short story about patience. When we first planted this orchard, twelve years ago, almost nothing grew in the first season. The soil was thin and the rains did not come on time. We thought we had failed. [pause 1s]

And yet — here we are. Here we are, with the orchard in full leaf, with cherries ripening in the upper beds, with plums in the lower, with quinces in the back corner where the wall holds the morning warmth. None of this came quickly. None of this came easily. All of it came.

What did the garden teach us? Patience, first. Then attention. Then humility — the kind of humility that shows up when you realize the garden is mostly in charge. [pause 2s]

We measured, we waited, we returned each morning to see what the night had done. We learned that the smallest change in light could move a whole row of seedlings from thriving to failing, and we learned that a gardener's best tool is often just a chair in the shade and a quiet hour.

Not every season is generous. Not every harvest is abundant. But every season teaches, and every harvest, however modest, is a gift that we did not make alone. The soil made it. The rain made it. The years of tending made it.

So this morning, as you walk the rows and gather what you like, I hope you will remember three things: that the garden is older than we are, that it will outlast us, and that it will reward any hour you are willing to give it. Thank you, and welcome.`

func TestComputeDemonstrationSpeechFormal(t *testing.T) {
	p := New(Config{Language: "en", Preset: Formal})
	r := p.Compute(demonstrationSpeech)
	// Expected delivery time for ~400 words at formal English (~4.5 syl/s)
	// is roughly 2-4 minutes. Allow a wide window for fixture drift.
	if r.TotalMS < 90_000 {
		t.Errorf("want at least 90 seconds delivery time, got %d ms", r.TotalMS)
	}
	if r.TotalMS > 360_000 {
		t.Errorf("want at most 6 minutes delivery time, got %d ms", r.TotalMS)
	}
	if len(r.BySentence) < 10 {
		t.Errorf("want at least 10 sentence timings, got %d", len(r.BySentence))
	}
}

func TestComputeDemonstrationSpeechPresetOrdering(t *testing.T) {
	cases := []Preset{Conversational, Formal, Ceremonial}
	var totals []int64
	for _, preset := range cases {
		p := New(Config{Language: "en", Preset: preset})
		totals = append(totals, p.Compute(demonstrationSpeech).TotalMS)
	}
	// Conversational < Formal < Ceremonial (faster to slower)
	if !(totals[0] < totals[1] && totals[1] < totals[2]) {
		t.Errorf("expected conversational < formal < ceremonial, got %v", totals)
	}
}

func TestComputeDemonstrationSpeechCrossLanguage(t *testing.T) {
	// Build a short fixture per language that uses the local
	// cultivation vocabulary. Each should return non-zero delivery.
	cases := []struct {
		lang string
		text string
	}{
		{"de", `Wir säen im Frühling. Wir säen im Sommer. Wir säen, wenn die Sonne zurückkehrt.`},
		{"es", `Cuidamos el jardín en primavera. Cuidamos el jardín en verano. Cuidamos el jardín cuando el sol regresa.`},
		{"fr", `Nous cultivons le jardin au printemps. Nous cultivons le jardin en été. Nous cultivons le jardin quand le soleil revient.`},
	}
	for _, c := range cases {
		t.Run(c.lang, func(t *testing.T) {
			p := New(Config{Language: c.lang, Preset: Formal})
			r := p.Compute(c.text)
			if r.TotalMS == 0 {
				t.Errorf("expected non-zero delivery time for %s, got 0", c.lang)
			}
			if r.Language != c.lang {
				t.Errorf("expected language %s, got %s", c.lang, r.Language)
			}
		})
	}
}

func TestComputeHonorsMarkersInLongText(t *testing.T) {
	p := New(Config{Language: "en", Preset: Formal})
	base := p.Compute(strings.ReplaceAll(demonstrationSpeech, "[pause 1s]", ""))
	withMarkers := p.Compute(demonstrationSpeech)
	// The demonstration speech has [pause 1s] + [pause 2s] = 3000 ms
	// of marker-driven silence. Total should differ by roughly that much.
	clean := strings.ReplaceAll(demonstrationSpeech, "[pause 1s]", "")
	clean = strings.ReplaceAll(clean, "[pause 2s]", "")
	baseline := p.Compute(clean)
	delta := withMarkers.TotalMS - baseline.TotalMS
	if delta < 2800 || delta > 3200 {
		t.Errorf("want ~3000ms from markers, got %d (base=%d, withMarkers=%d)",
			delta, base.TotalMS, withMarkers.TotalMS)
	}
}
