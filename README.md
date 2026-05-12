# Red-Team Calibration Harness

This repository calibrates release-gate review output across a fixed model
panel.

The experiment uses intentionally constrained low-reasoning implementer passes
in parallel, then asks reviewers to review the resulting small Go solver
packages. The active calibration tasks are product-style contracts written for
this repository, so memorized public benchmark answers are less useful. The
goal is to compare the distinct review signal each fixed model slot produces
from the same neutral prompt.

## Reviewers

The active run does not give reviewers specialty labels or specialty cues. Each
reviewer is asked to review as an experienced senior software engineer using
the same issue-register shape.

The checked-in reviewer slots are:

- Gemini 3.1 Pro
- Kimi K2.6
- Claude Opus 4.7
- DeepSeek V4
- Codex 5.5
- GLM 5.1

## Active Contracts

- Case A: utility-meter audit helper for reviewable periods in reading data.
- Case B: customer entitlement helper for account-id eligibility spans.
- Case C: lab inventory reconciliation helper for delivered and withdrawn
  compounds.
- Case D: care-plan fatigue scoring across continuous day runs.
- Case E: coaching coverage helper with limited support grants.
- Case F: readiness-gated checklist routing helper.

## Workflow

1. Run each batch of case prompts under `solver/case-*/` in parallel with a
   small-tier LLM at low reasoning effort.
2. Each solver writes real Go files: `solution.go`, `solution_test.go`, and
   `notes.md`.
3. Run targeted Go tests before review. Some checked-in solver outputs are
   intentionally flawed calibration artifacts; in particular, `solver/case-f`
   does not complete the full package test gate under a short timeout. See
   `synthesis/round-2-review-matrix.md`.
4. Ask each reviewer to review the solver implementations in scope for the
   current panel run with the same neutral prompt. This repo has two checked-in
   rounds: cases A-C under `reviews/*/cases/`, and cases D-F under
   `reviews/*/round-2/`.
5. Save reviewer outputs as consolidated Markdown files such as
   `reviews/<reviewer>/cases/full-panel.md` or
   `reviews/<reviewer>/round-2/full-panel.md`.
6. Record reviewer status and signal in synthesis files such as
   `synthesis/cross-review-matrix.md` and
   `synthesis/round-2-review-matrix.md`.
7. Tune panel instructions in `synthesis/panel-instruction-tuning.md`.

Some reviewer artifacts are partial because provider failures, timeouts, and
narrow retries are part of the calibration data. Treat the synthesis files as
the canonical status summaries.

## Evaluation Rule

Reviewer quality is judged by differentiated signal:

- Did the reviewer find true issues others missed?
- Did the issue match a useful recurring strength for that model slot?
- Did the reviewer avoid hallucinated blockers?
- Did the reviewer provide repro cases or precise counterexamples?
- Did the reviewer improve the final release-gate panel definition?

## Skills

The `skills/` directory contains checked-in copies of the two process skills
that came out of this calibration work:

- `red-team-release-gate`: a heavier serial release-gate review workflow.
- `simplify-with-review`: a lighter parallel simplification review workflow.

They are Codex-oriented skill files with command examples from this calibration
setup. The review principles are language-agnostic; adapt model names and CLI
commands to your own tooling.

## License

MIT. See `LICENSE`.
