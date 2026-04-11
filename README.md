# speechpace-go

**How long will this speech actually take to deliver?**

A Go package that answers the question every speechwriter asks on the day of the talk. Given a draft and a cadence, `speechpace-go` returns the expected delivery time in milliseconds, a per-sentence breakdown for editor rulers and waveform visualizations, and a list of sentences that exceed the breath limit a speaker can sustain in one pass.

The calculation is deterministic: syllables times seconds-per-syllable at the chosen cadence, plus per-language punctuation pauses, plus any explicit `[pause 2s]` markers in the text. No network calls, no machine learning, no surprises. Pick a language, pick a preset, hand over your draft.

This package was extracted from [CivNode](https://civnode.com), a European writing platform where speechwriters watch their delivery time update live as they type. Publishing it as a standalone library means anyone building a teleprompter, an audiobook platform, a podcast editor, or an accessibility tool can use the same timing engine.

---

## A live example

```go
package main

import (
	"fmt"

	"github.com/civnode/speechpace-go"
)

func main() {
	draft := `The garden asks for patience, for care, for time. ` +
		`[pause 1s] Not every season is generous. Not every harvest is abundant.`

	p := speechpace.New(speechpace.Config{
		Language: "en",
		Preset:   speechpace.Formal,
	})
	r := p.Compute(draft)
	fmt.Printf("Delivery: %d ms\n", r.TotalMS)
	fmt.Printf("Sentences: %d\n", len(r.BySentence))
	fmt.Printf("Breath violations: %d\n", len(r.BreathViolations))
}
```

Output:

```
Delivery: 9240 ms
Sentences: 3
Breath violations: 0
```

A little over nine seconds at formal cadence, accounting for the pause marker and the natural rest at every comma and period.

---

## Per-language articulation rates

The cadence presets are calibrated against comparable articulation-rate research on real speech (the 2011 Université de Lyon study on syllables-per-second across languages, and related linguistics work). Values are in syllables per second.

| Language | Conversational | Formal | Ceremonial | Breath ceiling |
|---|---|---|---|---|
| German | 4.8 | 3.8 | 3.2 | 35 syllables |
| English | 5.5 | 4.5 | 3.8 | 40 syllables |
| French | 6.0 | 5.0 | 4.2 | 45 syllables |
| Spanish | 6.5 | 5.5 | 4.6 | 50 syllables |

Spanish runs fastest, German runs slowest. That's why a translated speech that takes six minutes in English can take seven in German at the same preset. The breath ceiling scales with the rate: Spanish speakers can fit more syllables in one breath because the syllables themselves are faster.

**Conversational** is the natural rate of unhurried speech between two people in a quiet room. **Formal** is the podium-address rate used for keynotes, business presentations, and most professional delivery. **Ceremonial** is the slowest preset, used for eulogies, weddings, religious speech, and other occasions where gravity calls for unhurried pacing.

---

## What Compute returns

```go
type Result struct {
    TotalMS          int64              // full delivery time in milliseconds
    BySentence       []SentenceTiming   // per-sentence breakdown for rulers
    BreathViolations []int              // sentence indices exceeding breath ceiling
    Language         string             // normalized language code
}

type SentenceTiming struct {
    Index     int   // 0-based sentence index in the source
    Syllables int   // content syllables in the sentence
    SpeakMS   int64 // speaking time only
    PauseMS   int64 // punctuation + pause-marker pauses charged to this sentence
}
```

The `BySentence` slice gives editor rulers exactly what they need to render time labels next to sentences and waveform visualizations of pacing. `BreathViolations` flags sentences a speaker can't comfortably deliver in one breath — usually a prompt to split the sentence in two.

---

## Pause markers

Inline markers in the form `[pause 2s]` add literal silence to the delivery time. The marker is parsed and removed before sentence splitting, so it doesn't interfere with syllable counting or punctuation pauses.

```
The garden is ready. [pause 2s] Shall we begin?
```

Accepted forms:

```
[pause 2s]     -> 2.0 seconds
[pause 1.5s]   -> 1.5 seconds
[pause 0.5s]   -> 0.5 seconds
[pause .5s]    -> 0.5 seconds
[pause 10s]    -> 10.0 seconds
```

Malformed markers (`[pause]`, `[pause abc]`, `[pause -1s]`) are ignored silently — the detector treats them as literal text and moves on.

---

## Syllable counting

`speechpace-go` ships with `NaiveCounter`, a minimal vowel-group counter that approximates syllable counts well enough for tests and small speeches. Production systems should plug in a richer counter through the `SyllableCounter` interface:

```go
type SyllableCounter interface {
    CountSyllables(text string, language string) int
}
```

CivNode wires a [speedata/hyphenation](https://pkg.go.dev/github.com/speedata/hyphenation)-backed counter at the application layer, using embedded Liang hyphenation patterns for en-us, en-gb, de-1996, es, and fr. That implementation is significantly more accurate on compound words and accented vocabulary than the naive fallback, and it's a fifty-line adapter to write if you're building your own.

The package is deliberately zero-dependency at the module level, so downstream consumers don't pay for pattern files they might not need. Bring your own counter.

---

## Who this is for

- **Speechwriters** who need to know whether a draft fits the allotted time before rehearsal
- **Teleprompter and prompter software** needing a reliable per-language scroll speed
- **Audiobook platforms** estimating chapter length before recording
- **Accessibility tool builders** producing text-to-speech timing metadata
- **Podcast editors** estimating episode length from script
- **Conference organizers** checking whether a submitted keynote fits the session slot
- **Go developers** who appreciate small focused libraries with good test coverage

If you fall into any of those categories and something's missing, open an issue.

---

## Quick start

```sh
go get github.com/civnode/speechpace-go
```

```go
package main

import (
	"fmt"

	"github.com/civnode/speechpace-go"
)

func main() {
	p := speechpace.New(speechpace.Config{
		Language: "es",
		Preset:   speechpace.Ceremonial,
	})
	r := p.Compute(`Amigos, gracias por estar aquí esta tarde en el jardín.`)
	fmt.Printf("Delivery: %d ms at ceremonial Spanish cadence\n", r.TotalMS)
}
```

That's the entire API. One constructor, one method, a result struct. No async, no channels, no clients.

---

## Performance

Benchmarks on an AMD Ryzen 9 3950X with Go 1.22, using the naive built-in syllable counter:

| Input | Speed | Notes |
|---|---|---|
| One sentence | ~940 ns | Fast enough for on-keystroke detection in an editor |
| 400-word speech (a typical 3-minute talk) | ~25 µs | Essentially free |
| 10,000-word speech (a 60-minute keynote) | ~635 µs | Sub-millisecond at keynote length |

All four supported languages perform identically — the per-language table is a map lookup, not a pipeline stage. Switching between English, German, Spanish, and French costs nothing.

With a richer syllable counter plugged in (Liang patterns via speedata/hyphenation), expect the benchmarks to slow by a factor of two to three, which is still well below perceptible latency for editor feedback.

---

## Test coverage

54 test cases across the configuration, pacing table, punctuation weights, pause-marker parser, compute orchestrator, and a full integration test against a peaceful-themed demonstration speech. Every code path is exercised, every language is covered, every preset ordering is verified, and every marker edge case is pinned down. All test fixtures are synthesized on peaceful themes so nothing in the suite echoes historical political speeches.

---

## Contributing

Wanted, in rough priority order:

- **Calibration data from more languages.** Portuguese, Italian, Dutch, Polish, Turkish, Japanese, Mandarin. Each new language needs an articulation-rate row (conversational, formal, ceremonial) and a breath ceiling. Published linguistics research on comparable syllable rates is the best source.
- **Better published rate citations** in the README. Right now the existing table cites the Université de Lyon 2011 study as an anchor; pull requests that add primary-source citations for each rate are welcome.
- **Alternative syllable counters** — a Liang-pattern adapter that embeds TeX patterns for zero-configuration use, a CMU-dict adapter for English, a phonetic counter for languages with strict syllable rules (Spanish).
- **A streaming `Compute` variant** for very long texts where the caller wants to process one chunk at a time.

No contributor license agreement. Match the existing code style. Add tests for anything you change. The CI workflow includes a guard that fails if any LLM client, grammar service, or AI library creeps into `go.mod`, so don't add those.

---

## See also

Looking for rhetorical device detection — anaphora, tricolon, chiasmus, callback, and the rest of the classical canon? See [rhetoric-go](https://github.com/civnode/rhetoric-go), the companion package extracted from the same CivNode feature.

---

## License

MIT. See [LICENSE](LICENSE).

---

## Built for CivNode

Built for [CivNode](https://civnode.com) — a European writing platform with professional editing tools. All data stored in Germany, no US ties, GDPR by design.

CivNode uses `speechpace-go` inside its speech-writing feature: the delivery time at the top of the speech stats panel, the time labels on the left-margin reading ruler, the audience preset chips, the waveform timeline, the breath violation warnings — all of it driven by this package. The entire speech-writing feature runs deterministically on the German server, no tokens, no quotas, no cloud round-trip, no text leaving European jurisdiction. Speechwriters working on sensitive material — political speeches, eulogies, legal arguments — know their words stay put.

If you're writing a speech — a wedding toast, a eulogy, a TED talk, a political keynote, a commencement address — come try CivNode.

---

Built with care in Germany.
