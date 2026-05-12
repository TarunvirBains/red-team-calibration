# Round 2 Review Matrix

## Solver Status

- Case D, Case E: `go test ./solver/case-d ./solver/case-e` passed with
  46 tests across 2 packages.
- Case F: `timeout 8s go test ./solver/case-f` exited 124. The current
  implementation compiles but does not complete its own package tests in the
  bounded run.
- Prompt hygiene scan for the active round-2 prompts/templates returned no
  blocked cue terms after cleanup.

## Reviewer Status

| Reviewer | Status | Output |
| --- | --- | --- |
| Gemini 3.1 Pro | completed | `reviews/gemini-3.1-pro/round-2/full-panel.md` |
| Kimi K2.6 | completed | `reviews/kimi-k2.6/round-2/full-panel.md` |
| Claude Opus 4.7 | completed | `reviews/claude-opus-4.7/round-2/full-panel.md` |
| DeepSeek V4 | completed narrow retry | broad run hit provider 502; no-Novita autonomous Case F retry completed in `reviews/deepseek-v4/round-2/case-f-narrow-no-novita-autonomous.md` |
| Codex 5.5 | completed | `reviews/codex-5.5/round-2/full-panel.md` |
| GLM 5.1 | completed | `reviews/glm-5.1/round-2/full-panel.md` |

## Signal By Reviewer

Gemini gave the right first-pass map: it found the main Case D contract
failure, the Case E resource failure, and the Case F exact-time state blowup
without being seeded with known issues.

Kimi added strong edge and state signal: Case D large-value overflow, Case E
negative grant count and grant-size overflow, and Case F bad tests plus
near-MaxInt overflow.

Opus gave the best docs/tests release-gate pass: it identified that the Case D
tests and notes are wrong-anchored, Case E tests are weak, and Case F notes
describe a different approach from the implementation.

Codex added precise late-stage triage: accidental exported API in Case D,
realistic Case E scaling even at 20 staff with many sessions, and exact Case F
parity expectations for delayed 2x2 cases.

GLM added boundary/fix-shape value: negative `grantSize` handling in Case E
and a bounded parity-state formulation for Case F.

DeepSeek did not produce usable broad review text in this round. The broad run
hit a provider 502, and the first broad/narrow retries read the target files
but produced no final text. A no-Novita autonomous Case F retry did complete and
used tools effectively: it ran focused Go tests, confirmed the large-ready tests
are wrong/hanging, sharpened the notes issue into a provably incorrect
row/column-only algorithm claim, and showed that `rows + cols + maxReady`
overflow silently converts reachable cases into false `-1` results.

## Calibration Takeaways

- The serial order is working: Gemini first found broad blockers, then Kimi and
  Opus found different classes of follow-up issues instead of only repeating
  Gemini.
- The known-issue handoff reduced duplicate output without suppressing useful
  disagreement or sharper repros.
- Opus remains especially valuable for tests/docs/maintainer-risk depth.
- DeepSeek needs explicit tool autonomy and bounded review scope in this setup;
  with those in place, it produced useful stress evidence instead of stopping
  after file reads.
- Kimi and GLM were the most useful non-Codex sources for numeric, boundary,
  and state-transition details.
- Codex remains valuable late because it validates severity, adds file-line
  precision, and distinguishes confirmed issues from no-action checks.
