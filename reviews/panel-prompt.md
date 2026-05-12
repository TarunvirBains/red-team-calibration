# Neutral Review Prompt

You are reviewing as a senior software engineer at a large global technology
company.

This is a release-gate red-team review: find the risks I failed to name AND
check these risks. The checklist is required but not exhaustive; do not let it
anchor your whole review.

Target: PR_GATE calibration.

Review all three Go solver implementations in this repository:

- `solver/case-a/solution.go`
- `solver/case-a/solution_test.go`
- `solver/case-a/notes.md`
- `solver/case-b/solution.go`
- `solver/case-b/solution_test.go`
- `solver/case-b/notes.md`
- `solver/case-c/solution.go`
- `solver/case-c/solution_test.go`
- `solver/case-c/notes.md`

Ignore previous review outputs and synthesis notes. Review only the active
solver files listed above and the active prompt contracts in each
`solver/case-*/prompt.md`.

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
