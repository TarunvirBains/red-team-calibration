# Red-Team Calibration Harness

This repository calibrates reviewer roles for the release-gate red-team panel.

The experiment uses one intentionally constrained Rust solver pass, then asks
three reviewers to review all three hard algorithm problems. The goal is not to
prove who can solve LeetCode problems. The goal is to test whether each
specialist role produces distinct, useful review signal across different
correctness surfaces.

## Reviewers

- DeepSeek V4: Edge-Case Provocateur. Expected signal: boundary failures,
  stress cases, overflow, and performance cliffs.
- Kimi K2.6: Stateful Logic Analyst. Expected signal: invariant drift,
  long-chain state consistency, mutation errors, and failure propagation.
- GLM 5.1: Integration Skeptic. Expected signal: input contracts, parser/API
  boundaries, type mapping, malformed inputs, and external-interface behavior.

## Problems

- 327 Count of Range Sum: prefix sums, inclusive bounds, overflow, and
  performance-sensitive counting.
- 715 Range Module: mutable interval state across add/remove/query operations.
- 726 Number of Atoms: parser stack, token mapping, nested groups, and canonical
  output ordering.

## Workflow

1. Run the solver prompt in `solver/prompt.md`.
2. Save the combined answer to `solver/combined-solution.md`.
3. Split the relevant solution into each problem's `solution.md`.
4. Ask each reviewer to review all three solutions using its specialist lens.
5. Save reviewer outputs under `reviews/<reviewer>/<problem>.md`.
6. Fill in `synthesis/cross-review-matrix.md`.
7. Tune role descriptions in `synthesis/reviewer-role-tuning.md`.

## Evaluation Rule

Reviewer quality is judged by differentiated signal:

- Did the reviewer find true issues others missed?
- Did the issue match the role's intended specialty?
- Did the reviewer avoid hallucinated blockers?
- Did the reviewer provide repro cases or precise counterexamples?
- Did the reviewer improve the final release-gate role definition?
