package speechpace

import "strings"

// Compute runs the pacing calculation over the given text and returns
// a Result containing total delivery time, per-sentence breakdowns,
// and any breath-limit violations. An empty text returns a zero Result
// with no error.
func (p *Pace) Compute(text string) Result {
	if text == "" {
		return Result{Language: normalizeLang(p.cfg.Language)}
	}

	markers := ParsePauseMarkers(text)
	cleanText, markerMS := stripPauseMarkers(text, markers)

	rate := rateFor(p.cfg.Language, p.cfg.Preset)
	if rate <= 0 {
		rate = 4.5 // guard against misconfiguration
	}
	msPerSyllable := int64(1000.0 / rate)
	weights := pauseWeights(p.cfg.Language)
	ceiling := breathCeiling(p.cfg.Language)

	sentences := splitSentences(cleanText)
	result := Result{Language: normalizeLang(p.cfg.Language)}
	for i, s := range sentences {
		syls := p.cfg.SyllableCounter.CountSyllables(s, p.cfg.Language)
		speak := int64(syls) * msPerSyllable
		pause := punctuationPauseMS(s, weights)
		result.BySentence = append(result.BySentence, SentenceTiming{
			Index:     i,
			Syllables: syls,
			SpeakMS:   speak,
			PauseMS:   pause,
		})
		result.TotalMS += speak + pause
		if syls > ceiling {
			result.BreathViolations = append(result.BreathViolations, i)
		}
	}
	result.TotalMS += markerMS
	return result
}

// splitSentences is a minimal tokenizer shared with the test suite.
// speechpace keeps its own implementation so downstream consumers
// don't need to depend on rhetoric-go just to use the pacing
// calculator.
func splitSentences(text string) []string {
	if strings.TrimSpace(text) == "" {
		return nil
	}
	var result []string
	var current strings.Builder
	runes := []rune(text)
	for i, r := range runes {
		current.WriteRune(r)
		if r == '.' || r == '!' || r == '?' {
			// Look ahead to see if the next non-space rune starts a new sentence.
			atEnd := i == len(runes)-1
			nextIsBoundary := atEnd
			if !atEnd {
				for j := i + 1; j < len(runes); j++ {
					if runes[j] == ' ' || runes[j] == '\t' || runes[j] == '\n' {
						continue
					}
					if isUpper(runes[j]) {
						nextIsBoundary = true
					}
					break
				}
			}
			if nextIsBoundary {
				s := strings.TrimSpace(current.String())
				if s != "" {
					result = append(result, s)
				}
				current.Reset()
			}
		}
	}
	if s := strings.TrimSpace(current.String()); s != "" {
		result = append(result, s)
	}
	return result
}

func isUpper(r rune) bool {
	if r >= 'A' && r <= 'Z' {
		return true
	}
	if r >= 'À' && r <= 'Ý' {
		return true
	}
	return false
}
