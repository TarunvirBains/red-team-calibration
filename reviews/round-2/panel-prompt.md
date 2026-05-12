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
