# Red-Team Calibration Harness

This repository calibrates release-gate review output across a fixed model
panel.

The experiment uses intentionally constrained low-reasoning Haiku implementer
passes in parallel, then asks reviewers to review the resulting Go library
implementations. The active calibration tasks are product-style contracts
written for this repository, so memorized public benchmark answers are less
useful. The goal is to compare the distinct review signal each fixed model slot
produces from the same neutral prompt.

## Reviewers

The active run does not give reviewers specialty labels or specialty cues. Each
reviewer is asked to review as a senior software engineer at a large global
technology company using the same issue-register shape.

## Active Contracts

- Case A: utility-meter audit helper for reviewable periods in reading data.
- Case B: customer entitlement helper for account-id eligibility spans.
- Case C: lab inventory reconciliation helper for delivered and withdrawn
  compounds.
- Case D: care-plan fatigue scoring across continuous day runs.
- Case E: coaching coverage helper with limited support grants.
- Case F: readiness-gated checklist routing helper.

## Workflow

1. Run each batch of case prompts under `solver/case-*/` in parallel with
   Haiku at low reasoning.
2. Each solver writes real Go files: `solution.go`, `solution_test.go`, and
   `notes.md`.
3. Run targeted Go tests before review. Some checked-in solver outputs are
   intentionally flawed calibration artifacts, so the full `go test ./solver/...`
   gate may fail or time out after reviewers expose a blocker.
4. Ask each reviewer to review all three solver implementations with the same
   neutral prompt.
5. Save reviewer outputs as Markdown under `reviews/<reviewer>/cases/<case>.md`.
6. Fill in `synthesis/cross-review-matrix.md`.
7. Tune panel instructions in `synthesis/panel-instruction-tuning.md`.

## Evaluation Rule

Reviewer quality is judged by differentiated signal:

- Did the reviewer find true issues others missed?
- Did the issue match a useful recurring strength for that model slot?
- Did the reviewer avoid hallucinated blockers?
- Did the reviewer provide repro cases or precise counterexamples?
- Did the reviewer improve the final release-gate panel definition?

## License

MIT. See `LICENSE`.
