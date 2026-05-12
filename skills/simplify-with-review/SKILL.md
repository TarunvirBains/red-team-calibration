---
name: simplify-with-review
description: Use when the user asks to simplify, simplify-with-review, simplify-clean-clusters, clean up a non-trivial diff/range, deduplicate code or docs, reduce complexity, or run a simplification gate.
---

# Simplify With Review

## Principle

Simplification is behavior-preserving cleanup that makes a system easier to
understand, change, verify, scale, or operate. It must not narrow requirements,
hide invariants, delete difficult tests, or turn a real product/design question
into a silent deferral.

For non-trivial work, completion requires real file/commit grounding, local
verification, and independent review or an explicit human stop.

## Codex Rules

- The main Codex agent owns scope, edits, verification, reviewer reconciliation,
  and final judgment.
- Use `multi_tool_use.parallel` for independent reads/searches and
  `apply_patch` for manual edits.
- Use `spawn_agent` only when the user authorized subagents, delegation, or
  independent review. Prefer read-only explorers for review.
- Do not set model overrides unless requested or technically justified. When
  authorized for final adversarial review, prefer high/xhigh reasoning.
- Give reviewers compact current context: scope, diff/files, invariants,
  verification, non-goals, gaps, and required verdict line.
- Ask reviewers for novel issues and sibling-pattern gaps before named risks.
- Reviewer model pins fail closed. Do not silently substitute lower or adjacent
  variants; mark the slot `partial`/`skipped` and ask before substituting.

## Workflow

### 1. Scope

- Identify base, tip, branch/PR/issue, commit range, workstream, and files.
- Check worktree status and preserve user changes.
- If the boundary is ambiguous, inspect local evidence and propose the scope
  before editing.
- If the user asked only for review/planning, do not edit without confirmation.

### 2. Chunk Large Ranges

Split milestone-scale work before editing. A good chunk has one subsystem or
bug class, one primary verification strategy, auditably small read/write
surfaces, and minimal overlap with other chunks.

Useful chunk families: core model, public API, generated output, storage/state
transfer, user workflow, tests/docs/examples, automation/tooling.

For each chunk, capture:

```text
name:
purpose:
included commits/files:
write/read surface:
sibling audit targets:
verification:
order/parallelism:
risks:
```

### 3. Read And Simplify

- Read changed files plus adjacent sibling sites.
- Search for repeated names, copied helpers, stale comments/docs, old API
  shapes, duplicate fixtures, equivalent branches, and repeated setup.
- Consolidate real duplication, clarify invariants, align docs/tests/errors
  with behavior, and remove avoidable work without weakening correctness.
- Prefer small behavior-preserving commits. Keep unrelated cleanup out.
- Do not delete tests because behavior is hard to express; rewrite them or
  escalate the gap.
- Do not change public semantics unless API/product cleanup is explicitly in
  scope and tests/docs move with it.
- Preserve the expected growth path; cleanup should make later work easier, not
  merely make today's diff smaller.

For tiny one-file cleanup, stop after this step plus verification unless the
user explicitly asked for independent review.

### 4. Verify

- Run the smallest command that proves each cleanup, then the repo's expected
  gate for the touched surface.
- The orchestrator owns long or fragile checks: full suites, local CI emulation,
  external services, browser/integration checks, or release checks.
- Do not claim completion without fresh output and exit status.

### 5. Review

Dispatch reviewers only after a coherent diff exists and verification is green,
or the remaining verification gap is explicit.

Reviewer constraints: read-only; no edits; no nested delegation; no long
verification unless assigned; audit sibling sites; cite concrete evidence. Use
1800s reviewer timeouts by default and up to 3600s for max-effort/deep-reasoning
slots or complex chunks.

When an independent panel is authorized, use non-opencode reviewers first:

- Gemini 3.1 Pro via Gemini CLI, `gemini-3.1-pro-preview`, highest available
  reasoning: product/contract fit, broad blind spots, docs/tests/workflow risk.
- Claude Opus 4.7 via Claude Code CLI, `careful-coder`, max effort, read-only:
  architecture, API misuse, security, maintainability, docs/examples.
- Codex 5.5 at `xhigh`: correctness, root cause, file-line evidence, focused
  repros, final integration judgment.

Run this review panel in parallel by default for quick broad-stroke coverage.
Verify local CLI flags before first use and wrap reviewer calls with
`rtk timeout` using the budgets above. If a reviewer is slow but visibly
progressing, mark it `pending`, finish triage with available outputs, and merge
late findings if they arrive before final judgment.

Reviewer prompt shape:

```text
Repo/path:
Base/tip/commit range:
Chunk scope:
Simplification commits:
Files touched:
Behavior that must remain unchanged:
Duplication / risk family:
Sibling sites to audit:
Non-goals:
Verification already run:
Known gaps:

Review as a senior software engineer at a large multinational technology
company. Be exhaustive, structured, and correctness-oriented within your assigned
specialization. Use three passes:
1. Contract and behavior preservation.
2. Edge cases, sibling patterns, state/resource/concurrency risks, and tests.
3. Docs/examples/API usability, maintainability, and adoption risk.

Find novel issues and sibling-pattern gaps first, then check named risks and
known findings. Repeat known findings only for better evidence, repro, severity,
or fix. Include important no-issue areas checked.
For each finding include severity, area, evidence, impact, and suggested fix.
End with exactly one line: VERDICT: ALLOW or VERDICT: BLOCK.
If BLOCK, enumerate BLOCK-1, BLOCK-2, ... above the verdict.
```

If a follow-up round is needed, pass a concise issue register instead of full
transcripts:

```text
Known issues so far. Do not repeat these unless adding new evidence, a sharper
repro, a severity disagreement, or a better fix:

- [AREA] claim
  evidence:
  suggested fix:

Your task: find novel issues first; then report only material disagreements,
better evidence, and important areas checked.
```

### 6. Fix And Re-review

- The implementer/orchestrator reviews every `BLOCK` finding against code and
  evidence before it becomes blocking.
- Reviewer feedback is advisory until validated. Promote only concrete
  correctness, safety, public API, or behavior-preservation issues to blockers.
- Verify findings against code before changing anything.
- Fix valid blockers in atomic commits, re-run relevant verification, and
  re-review affected chunks until the latest round allows or a stop condition
  applies.

## Stop Conditions

Stop for human judgment when scope is materially ambiguous, simplification
would change public behavior, reviewers find a broader product/API/design
issue, review rounds do not converge, evidence is ambiguous, required
verification cannot run, or worktree state blocks safe cleanup.

## Anti-patterns

- Reviewer majority vote.
- Shortening code by hiding invariants.
- Leaving sibling sites on the old pattern.
- Mixing unrelated style churn into behavior-sensitive cleanup.
