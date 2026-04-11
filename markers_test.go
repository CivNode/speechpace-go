package speechpace

import "testing"

func TestParsePauseMarkersBasic(t *testing.T) {
	text := "Hello. [pause 2s] World. [pause 1.5s] End."
	got := ParsePauseMarkers(text)
	if len(got) != 2 {
		t.Fatalf("want 2 markers, got %d", len(got))
	}
	if got[0].Seconds != 2.0 {
		t.Errorf("first marker: want 2.0s, got %f", got[0].Seconds)
	}
	if got[1].Seconds != 1.5 {
		t.Errorf("second marker: want 1.5s, got %f", got[1].Seconds)
	}
}

func TestParsePauseMarkersDecimal(t *testing.T) {
	cases := []struct {
		text string
		want float64
	}{
		{"[pause 0.5s]", 0.5},
		{"[pause .5s]", 0.5},
		{"[pause 2s]", 2.0},
		{"[pause 10s]", 10.0},
		{"[pause 1.25s]", 1.25},
	}
	for _, c := range cases {
		t.Run(c.text, func(t *testing.T) {
			m := ParsePauseMarkers(c.text)
			if len(m) != 1 {
				t.Fatalf("want 1 marker, got %d", len(m))
			}
			if m[0].Seconds != c.want {
				t.Errorf("want %f, got %f", c.want, m[0].Seconds)
			}
		})
	}
}

func TestParsePauseMarkersNoMatches(t *testing.T) {
	text := "Hello world. No markers here."
	got := ParsePauseMarkers(text)
	if len(got) != 0 {
		t.Errorf("want 0 markers, got %d", len(got))
	}
}

func TestParsePauseMarkersMalformed(t *testing.T) {
	cases := []string{
		"[pause]",
		"[pause abc]",
		"[pause -1s]",
		"(pause 2s)",
		"{pause 2s}",
	}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			got := ParsePauseMarkers(c)
			if len(got) != 0 {
				t.Errorf("malformed marker should not parse: %q -> %v", c, got)
			}
		})
	}
}

func TestStripPauseMarkersRemovesAll(t *testing.T) {
	text := "Hello. [pause 2s] World. [pause 1.5s] End."
	markers := ParsePauseMarkers(text)
	clean, ms := stripPauseMarkers(text, markers)
	if contains2(clean, "pause") {
		t.Errorf("stripped text should not contain 'pause', got %q", clean)
	}
	if ms != 3500 {
		t.Errorf("want 3500ms total, got %d", ms)
	}
}

func TestStripPauseMarkersNoMarkersReturnsOriginal(t *testing.T) {
	text := "No markers in this text."
	clean, ms := stripPauseMarkers(text, nil)
	if clean != text {
		t.Errorf("want original text, got %q", clean)
	}
	if ms != 0 {
		t.Errorf("want 0ms, got %d", ms)
	}
}

func TestStripPauseMarkersEmpty(t *testing.T) {
	clean, ms := stripPauseMarkers("", nil)
	if clean != "" || ms != 0 {
		t.Errorf("want empty result, got %q / %d", clean, ms)
	}
}

func contains2(s, needle string) bool {
	for i := 0; i+len(needle) <= len(s); i++ {
		if s[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
