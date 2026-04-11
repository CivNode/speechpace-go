package speechpace

// PauseWeights holds per-punctuation pause durations in milliseconds.
// Each value represents the silent gap a speaker leaves at the given
// punctuation mark during ordinary delivery. The values vary by
// language because speakers of different languages cluster their
// pauses differently: German readers rest longer on commas to let
// compound words settle; Spanish readers barrel through light
// punctuation more quickly.
type PauseWeights struct {
	Period    int
	Semicolon int
	Comma     int
	EmDash    int
	Question  int
	Exclaim   int
	Colon     int
}

// pauseWeights returns the pause weights for the given language code.
// Unknown languages fall back to English values.
func pauseWeights(language string) PauseWeights {
	switch normalizeLang(language) {
	case "de":
		return PauseWeights{Period: 450, Semicolon: 350, Comma: 220, EmDash: 300, Question: 450, Exclaim: 450, Colon: 320}
	case "es":
		return PauseWeights{Period: 380, Semicolon: 280, Comma: 170, EmDash: 260, Question: 380, Exclaim: 380, Colon: 260}
	case "fr":
		return PauseWeights{Period: 400, Semicolon: 300, Comma: 180, EmDash: 280, Question: 400, Exclaim: 400, Colon: 280}
	}
	return PauseWeights{Period: 400, Semicolon: 300, Comma: 180, EmDash: 280, Question: 400, Exclaim: 400, Colon: 280}
}

// punctuationPauseMS sums the pause milliseconds charged to a single
// sentence based on its punctuation content.
func punctuationPauseMS(sentence string, w PauseWeights) int64 {
	var total int64
	for _, r := range sentence {
		switch r {
		case '.':
			total += int64(w.Period)
		case ';':
			total += int64(w.Semicolon)
		case ',':
			total += int64(w.Comma)
		case '—':
			total += int64(w.EmDash)
		case '?':
			total += int64(w.Question)
		case '!':
			total += int64(w.Exclaim)
		case ':':
			total += int64(w.Colon)
		}
	}
	return total
}
