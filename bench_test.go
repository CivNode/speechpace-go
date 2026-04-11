package speechpace

import (
	"strings"
	"testing"
)

// BenchmarkComputeDemonstrationSpeech measures orchestration cost on
// the peaceful-themed demonstration speech (~400 words).
func BenchmarkComputeDemonstrationSpeech(b *testing.B) {
	p := New(Config{Language: "en", Preset: Formal})
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Compute(demonstrationSpeech)
	}
}

// BenchmarkComputeLargeSpeech measures scaling on a 10,000-word
// synthetic speech built by repeating the demonstration fixture.
// Target: sub-millisecond per call.
func BenchmarkComputeLargeSpeech(b *testing.B) {
	p := New(Config{Language: "en", Preset: Formal})
	large := strings.Repeat(demonstrationSpeech, 25)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Compute(large)
	}
}

// BenchmarkComputeSingleSentence measures minimum per-call overhead.
func BenchmarkComputeSingleSentence(b *testing.B) {
	p := New(Config{Language: "en", Preset: Formal})
	text := `The garden asks for patience, for care, for time.`
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.Compute(text)
	}
}

// BenchmarkComputeAllLanguages compares per-language cost on the
// same fixture to surface any per-language overhead.
func BenchmarkComputeAllLanguages(b *testing.B) {
	langs := []string{"en", "de", "es", "fr"}
	for _, lang := range langs {
		b.Run(lang, func(b *testing.B) {
			p := New(Config{Language: lang, Preset: Formal})
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = p.Compute(demonstrationSpeech)
			}
		})
	}
}
