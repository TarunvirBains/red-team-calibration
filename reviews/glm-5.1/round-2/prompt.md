# Neutral Review Prompt

You are reviewing as a senior software engineer at a large global technology
company.

This is a release-gate red-team review: find the risks I failed to name AND
check these risks. The checklist is required but not exhaustive; do not let it
anchor your whole review.

Target: PR_GATE calibration, round 2.

Review these Go solver implementations in this repository:

- `solver/case-d/solution.go`
- `solver/case-d/solution_test.go`
- `solver/case-d/notes.md`
- `solver/case-e/solution.go`
- `solver/case-e/solution_test.go`
- `solver/case-e/notes.md`
- `solver/case-f/solution.go`
- `solver/case-f/solution_test.go`
- `solver/case-f/notes.md`

Ignore previous review outputs and synthesis notes. Review only the solver
files listed above and the prompt contracts in each listed directory.

Known issues so far. Do not repeat these unless adding new evidence, a sharper
repro, a severity disagreement, or a better fix:

- Case D computes only same-value consecutive groups. The contract appears to
  require every continuous day span. Repro: `[]int{1, 2}` should score `8`,
  but the implementation returns `5`. Current publish decision: BLOCK.
- Case D tests and notes are anchored to the same wrong same-value grouping
  interpretation. Current publish decision: BLOCK.
- Case D can overflow before modular reduction when very large daily values
  are added. Current publish decision: BLOCK.
- Case E enumerates staff grant subsets with `1 << uint(len(staff))`. Large
  staff counts can hang, overflow, return zero silently, or become unusable in
  realistic scheduling sizes. Current publish decision: BLOCK.
- Case E negative `grantCount` causes zero coverage even when no grant is
  needed. Current publish decision: BLOCK.
- Case E can overflow when adding `grantSize` to staff capacity. Current
  publish decision: FIX_BEFORE_TARGET.
- Case E tests do not strongly validate the tricky matching behavior and lack
  a small-reference cross-check. Current publish decision: FIX_BEFORE_TARGET.
- Case F tracks position and exact minute up to the largest readiness value.
  Large readiness values can consume unbounded time and memory. Current publish
  decision: BLOCK.
- Case F notes misstate visited-state shape and resource cost, and describe a
  different approach from the code. Current publish decision: BLOCK.
- Case F tests assert unreachable for 2x2 layouts that can be reached by
  moving back and forth until the destination is ready. Current publish
  decision: BLOCK.
- Case F can overflow when computing `rows + cols + maxReady`. Current publish
  decision: BLOCK.
- Case F queue reslicing retains backing storage and amplifies the memory
  failure. Current publish decision: POST_TARGET once the main blocker is fixed.

Your task:
1. Find novel issues and blind spots first.
2. Check whether any known issue is materially wrong or mis-severitized.
3. Report only new findings, material disagreements, better evidence, and
   important areas checked.

For each case, return findings in this shape:

```text
area:
finding:
evidence:
repro/test idea:
impact:
likelihood:
fix size:
publish decision: BLOCK | FIX_BEFORE_TARGET | DOC_BEFORE_TARGET | POST_TARGET | NO_ACTION
```

Do not use delegation, multi-agent, cloud review, or nested reviewer tools.
Do not edit files. Use direct analysis only.
