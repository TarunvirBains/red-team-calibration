# Narrow Review Prompt

You are reviewing as a senior software engineer at a large global technology
company.

Target: PR_GATE calibration, round 2, Case F only.

Review only these files:

- `solver/case-f/prompt.md`
- `solver/case-f/solution.go`
- `solver/case-f/solution_test.go`
- `solver/case-f/notes.md`

Known issues so far. Do not repeat these unless adding new evidence, a sharper
repro, a severity disagreement, or a better fix:

- Case F tracks position and exact minute up to the largest readiness value.
  Large readiness values can consume unbounded time and memory. Current publish
  decision: BLOCK.
- Case F notes misstate visited-state shape and resource cost, and describe a
  different approach from the code. Current publish decision: BLOCK.
- Case F tests assert unreachable for 2x2 layouts that can be reached by moving
  back and forth until the destination is ready. Current publish decision:
  BLOCK.
- Case F can overflow when computing `rows + cols + maxReady`. Current publish
  decision: BLOCK.
- Case F queue reslicing retains backing storage and amplifies the memory
  failure. Current publish decision: POST_TARGET once the main blocker is fixed.
- Case F should add single-row/single-column delayed-neighbor tests so a future
  fix does not incorrectly allow waiting in place. Current publish decision:
  FIX_BEFORE_TARGET.

Your task:
1. Find only novel stress, boundary, runtime, resource, or correctness issues.
2. Check whether any known issue is materially wrong or mis-severitized.
3. Report only new findings, material disagreements, better evidence, and
   important areas checked.

For each finding, use this shape:

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
