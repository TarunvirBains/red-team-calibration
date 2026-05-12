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

Be exhaustive: do not stop after the first few obvious findings, and do not
optimize for brevity.

Use three passes:
1. Contract compliance against prompts, specs, docs, examples, and public API.
2. Edge cases, malformed or extreme inputs, runtime failure modes, numeric
   limits, state drift, concurrency, and resource exhaustion.
3. Tests, docs, exported API usability, release packaging, adoption risks, and
   workflow contradictions.

Interleave analysis with evidence gathering: read the relevant files, form a
candidate issue, then look for the exact code, test, doc, or command evidence
that proves or disproves it. Do not ask for or output private chain-of-thought;
make the work auditable through findings, areas checked, assumptions, and
severity calibration.

For each target, include:
- material findings using the issue register;
- important areas checked where you found no issue;
- assumptions that control severity;
- one final release recommendation.

Do not invent findings to fill space.
