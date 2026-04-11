package speechpace

import (
	"regexp"
	"strconv"
)

// PauseMarker is an inline pause directive parsed from source text,
// typically of the form [pause 2s] or [pause 1.5s]. Speech writers use
// markers to insert explicit silence beyond what punctuation alone
// would produce.
type PauseMarker struct {
	Start   int     // byte offset in the source where the marker begins
	End     int     // byte offset where it ends (exclusive)
	Seconds float64 // duration in seconds
}

// Regex accepts [pause 2s], [pause 1.5s], [pause 0.5s], [pause .5s],
// and [pause 2]. Whitespace is tolerated around the number and unit.
var pauseRegex = regexp.MustCompile(`\[pause\s+(\d+(?:\.\d+)?|\.\d+)s?\s*\]`)

// ParsePauseMarkers returns every pause marker present in the input
// text, preserving source order. Markers overlapping each other are
// unusual but handled in first-match order.
func ParsePauseMarkers(text string) []PauseMarker {
	var out []PauseMarker
	for _, m := range pauseRegex.FindAllStringSubmatchIndex(text, -1) {
		secStr := text[m[2]:m[3]]
		sec, err := strconv.ParseFloat(secStr, 64)
		if err != nil {
			continue
		}
		out = append(out, PauseMarker{Start: m[0], End: m[1], Seconds: sec})
	}
	return out
}

// stripPauseMarkers removes pause marker substrings from the input,
// returning a clean text string plus the total milliseconds of pauses
// extracted. Callers use the cleaned text for sentence splitting and
// syllable counting, then add the marker milliseconds to the final
// Result.TotalMS.
func stripPauseMarkers(text string, markers []PauseMarker) (string, int64) {
	if len(markers) == 0 {
		return text, 0
	}
	var total int64
	var b []byte
	cursor := 0
	for _, m := range markers {
		b = append(b, text[cursor:m.Start]...)
		cursor = m.End
		total += int64(m.Seconds * 1000)
	}
	b = append(b, text[cursor:]...)
	return string(b), total
}
