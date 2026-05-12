---
name: red-team-release-gate
description: Use when preparing an alpha, beta, stable release, red-team review, security review, penetration-test-style code assessment, or publish gate for a crate, package, service, CLI, SDK, library, or app where correctness, isolation, persistence, auth, generated code, docs, or public API misuse could create serious release risk.
---

# Red-Team Release Gate

## Purpose

Release-gate review for code and docs. Make the release harder to misuse,
harder to corrupt, and honest about its contracts. This is review/triage, not
live exploitation. Stay read-only unless the user explicitly assigns fixes; do
not run destructive commands.

Use before publishing when the project handles user data; crosses process,
runtime, file, network, database, tenant, auth, or trust boundaries; depends on
cache/refresh/TTL/LRU/tombstones/recovery/identity/query isolation/generated
code/concurrency/wire formats; or has docs/examples that could teach unsafe
use. A plausible correctness or security footgun blocks until fixed, disproven,
documented as a deliberate limitation, or rescoped by the user.

## Core Prompt

Bind the target before review: `ALPHA`, `BETA`, `STABLE`, `PR_GATE`,
`SECURITY_REVIEW`, or a named custom target.

```text
This is a release-gate red-team review: find the risks I failed to name AND
check these risks. The checklist is required but not exhaustive; do not let it
anchor your whole review.

Target: <target>.

Do not use delegation, multi-agent, cloud review, or nested reviewer tools;
perform direct analysis only.

For each finding, include area, evidence, repro/test idea, impact, likelihood,
fix size, and publish decision.
```

Issue shape:

```text
area:
finding:
evidence:
repro/test idea:
impact:
likelihood:
fix size:
publish decision:
```

Publish decisions:

- `BLOCK`: do not publish/merge/approve until fixed, disproven, documented as
  user-approved non-goal, or consciously rescoped.
- `FIX_BEFORE_TARGET`: blocks the named target; lower-risk scope is allowed
  only if the user changes target/scope.
- `DOC_BEFORE_TARGET`: behavior is acceptable, but docs/API wording can cause
  unsafe misuse.
- `POST_TARGET`: track after this target.
- `NO_ACTION`: checked and safe enough with current evidence.

Confirmed bugs outrank speculative hardening. Suspicion becomes a blocker when
there is a plausible repro path and meaningful impact.

## Scope Setup

Before panel dispatch, perform local triage and confirm with the user:

- Scope, release target, and owner assumptions.
- Known active fixes/findings for de-duplication.
- High-risk areas to check.
- Planned reviewer order and exact tools/models.

Headless/CI may proceed only when target, scope, assumptions, known findings,
high-risk areas, and reviewer order are supplied by user instruction or trusted
repo config. Otherwise stop with `MISSING_CONFIRMATION` and print the proposed
plan.

Put repo-specific assumptions in the prompt, issue register, or final report,
not in this skill. Use owner assumptions to reduce false positives, but verify
that code, docs, defaults, and examples match the stated contract.

## Risk Checklist

Adapt to the repo, then check:

- Identity, stable IDs, cache keys, collisions, sorting.
- Tenant/auth/RLS/capability checks, query isolation, confused deputy paths.
- Stale data: refresh, invalidation, TTL, LRU, tombstones, recovery, retries.
- Backend keyspace, filesystem, Redis, database, network, wire formats.
- Feature gates, runtime/platform assumptions, optional deps, unsafe defaults.
- Generated code, migrations, serialization, public APIs, docs/API mismatch.
- Packaging: included files, feature flags, versioning, changelog, migration
  notes, credentials, secrets, fixtures.

## Panel Discipline

Use the Standing Reviewer Panel below when available and allowed. Run serially
by default so each reviewer can find different issues. Treat reviewers as a
panel, not a vote: one well-evidenced blocker stops the release.

Model/provider pins are release controls. They fail closed if unavailable. Do
not silently substitute lower or adjacent variants such as `gpt-5.4 high`,
`Claude Sonnet`, or unpinned provider defaults. If a slot cannot be pinned, mark
it `partial` or `skipped`, record why, and ask before substituting.

The orchestrator owns all delegation. Reviewers must do direct analysis only.
Forbid or disable reviewer CLI tools that spawn tasks, explorers, ultrareview,
subagents, model panels, edits, or writes. If a reviewer delegates anyway,
preserve usable evidence and rerun once with a tighter direct-review prompt. If
it still fails, mark `partial`/`skipped` and continue.

Reviewer timeout rule: use 1800s for routine slots and up to 3600s for
max-effort/deep-reasoning slots or complex targets. Do not kill a simply slow
reviewer by default; if progress is visible, mark it `pending`, start the next
reviewer with the current register, and merge late findings. Note the
de-duplication gap in the final report.

After each reviewer, add confirmed/plausible findings to the running register.
Pass concise known issues to the next reviewer; do not paste full transcripts
unless the task is explicit de-duplication.

Handoff shape:

```text
Known issues so far. Do not repeat these unless adding new evidence, a sharper
repro, a severity disagreement, or a better fix:

- [AREA] claim
  evidence:
  repro/test idea:
  current publish decision:

Your task:
1. Find novel issues and blind spots first.
2. Check whether any known issue is materially wrong or mis-severitized.
3. Report only new findings, material disagreements, better evidence, and
   important areas checked.
```

## Standing Reviewer Panel

Order:

1. Gemini 3.1 Pro via Gemini CLI, `gemini-3.1-pro-preview`, highest available
   reasoning, as Product Manager & Senior Engineer: broad correctness, adopter
   impact, feasibility, packaging, workflow contradictions, blind spots.
2. Kimi K2.6 via OpenRouter/opencode, max variant, as Stateful Logic Analyst:
   long-chain data flow, async/multi-step state consistency, design blind spots,
   alternative-reasoning review.
3. Claude Opus 4.7 via Claude Code CLI, `careful-coder`, max effort, read-only,
   as Architectural Critic & Security Specialist: docs/API misuse, architecture,
   vulnerabilities, unsafe examples, release risk.
4. DeepSeek V4 via OpenRouter/opencode, max variant, as Edge-Case Provocateur:
   subsystem bug hunts, stress paths, performance cliffs, implementation risk.
5. Codex 5.5 at `xhigh` as Logic & Security Auditor: root cause, mathematical
   correctness, file-line evidence, focused repros, TDD fixes, final integration
   judgment.
6. GLM 5.1 via OpenRouter/opencode, max variant, as Integration Skeptic:
   boundaries, type/API mappings, malformed dependency behavior, adversarial
   consensus check.

Codex is not the final slot; the orchestrator adjudicates all outputs regardless
of reviewer order. If GLM or any late reviewer raises a logic/security issue
after Codex, validate it with Codex-level rigor during triage.

## Pin And Routing Checks

Before each slot, record exact model, provider/routing, and reasoning controls.

- Codex: `codex exec -m gpt-5.5 -c 'model_reasoning_effort="xhigh"'`.
- Claude: `claude -p --agent careful-coder --model opus --effort max
  --permission-mode plan`; allow `Read,Grep,Glob`; disallow
  `Task,Bash,Edit,Write,MultiEdit,NotebookEdit,WebFetch`; no Sonnet/default
  fallback.
- Gemini: `gemini --model gemini-3.1-pro-preview`; verify local config applies
  `thinkingConfig.thinkingLevel = "HIGH"` or mark `partial`.
- OpenRouter/opencode IDs:
  `openrouter/moonshotai/kimi-k2.6`,
  `openrouter/deepseek/deepseek-v4-pro`,
  `openrouter/z-ai/glm-5.1`.

For OpenRouter reviewers, always run a provider preflight before dispatch:

- endpoints: `GET /api/v1/models/$OR_MODEL/endpoints`
- canary each allowed/currently suspect provider with
  `POST /api/v1/chat/completions` and
  `provider.only: ["$OR_PROVIDER"]`, `allow_fallbacks: false`, small
  `max_tokens`
- metadata: `GET /api/v1/generation?id=$GENERATION_ID` and record
  `provider_name`

Then dispatch with deterministic routing to a healthy provider using
`provider.only` and `allow_fallbacks:false`. Exclude only failed or
user-blacklisted providers, and only for that run. Set the selected
`options.provider.only` in opencode config before the reviewer command.

Do not permanently blacklist a provider because it failed once; routers change.
Do not infer the provider from opencode silence. Run `opencode debug config`
from the repo root and record exact model IDs plus provider policy. If opencode
exposes serving provider metadata, record it. Otherwise say only what routing
policy was requested; do not claim which backend served the completion. A
passing preflight is routing/availability evidence only.

Wrap reviewer calls with `rtk timeout`: 1800s default, 3600s for max-effort or
known-deep slots.

OpenRouter/opencode shape:

```bash
rtk timeout 3600s opencode run "$REVIEW_MESSAGE" --dir "$REPO" \
  -m openrouter/moonshotai/kimi-k2.6 \
  --variant max --dangerously-skip-permissions \
  --file "$PROMPT_FILE" > "$OUT_FILE" 2> "$ERR_FILE"
```

Swap only the pinned model ID for DeepSeek or GLM. Add to the reviewer prompt:

```text
You have full local tool autonomy for this review: use read, search, shell,
tests, and temporary scratch harnesses when useful. Do not modify tracked repo
files. Do not stop after a tool_use event; after tool checks, produce the final
markdown review.
```

If an OpenRouter reviewer reads files or starts a tool and then produces no
final text, inspect stderr, active subprocesses, opencode logs, and any
available opencode export/session recovery before declaring output lost. If
still streaming or running a bounded test, wait. If idle/wedged, preserve the
transcript and retry once with the autonomy prompt, same provider policy, and
narrower scope. If it still fails, mark `partial`.

## Prompt Template And Commands

Codex command:

```bash
codex exec -C "$REPO" -m gpt-5.5 \
  --sandbox read-only --ephemeral \
  -c 'model_reasoning_effort="xhigh"' \
  -o "$OUT_FILE" \
  - < "$PROMPT_FILE" > "$LOG_FILE" 2>&1
```

Use stdin for long prompts and capture final response separately from logs. Do
not use unsupported approval flags. Include project context, target, known
findings, high-risk areas, and the running register.

Append this shared depth template to every reviewer prompt, followed by that
slot's specialization and the running known-issue handoff:

```text
Review as a senior software engineer at a large multinational technology
company. Be exhaustive, structured, and correctness-oriented within your assigned
specialization. Do not stop at obvious findings.

Use three passes:
1. Contract compliance against prompts, specs, docs, examples, and public API.
2. Edge cases, malformed/extreme inputs, runtime failure modes, numeric limits,
   state drift, concurrency, and resource exhaustion.
3. Tests, docs, exported API usability, packaging, adoption risks, and workflow
   contradictions.

Interleave analysis with evidence gathering. Do not expose private
chain-of-thought; make work auditable through findings, areas checked,
assumptions, and severity.

Find novel issues first. Repeat known issues only for better evidence, repro,
severity, or fix. Include material findings, important no-issue areas checked,
severity assumptions, and one release recommendation. Do not invent findings.
```

Gemini command: `gemini --model gemini-3.1-pro-preview --approval-mode plan
--skip-trust --output-format text -p "$(cat "$PROMPT_FILE")" > "$OUT_FILE"
2> "$ERR_FILE"`. Verify local CLI flags before first use; if a flag is
unsupported, update the command and record the change.

Claude command: `claude -p --agent careful-coder --model opus --effort max
--permission-mode plan --tools Read,Grep,Glob --disallowedTools
Task,Bash,Edit,Write,MultiEdit,NotebookEdit,WebFetch --debug-file
"$DEBUG_FILE" --output-format text --no-session-persistence <
"$PROMPT_FILE"`. Verify local CLI flags before first use.

## Orchestration And Triage

Track each reviewer: start time, first visible output, completion time, elapsed
time, status (`completed`, `pending`, `partial`, `skipped`). Include timing for
slow, partial, skipped, or late-merged reviewers.

For every finding:

1. Restate the technical claim.
2. Verify against code, tests, docs, packaging, or runtime behavior.
3. Classify as confirmed bug, plausible hazard, documentation hazard, or false
   positive.
4. Do not blindly implement reviewer suggestions; verify fit with API and
   release goals.
5. For code fixes, use TDD: failing test, observed failure, fix.
6. For docs fixes, remove overclaims and state boundaries plainly.
7. Re-run targeted tests, then relevant broader tests.

After triage, create/update GitHub issues for unresolved `BLOCK`,
`FIX_BEFORE_TARGET`, `DOC_BEFORE_TARGET`, and `POST_TARGET` findings when access
is available. Check for duplicates first. Preserve evidence, repro/test idea,
impact, likelihood, fix size, and decision. Add labels such as `release-gate`,
`red-team`, `security`, `correctness`, `docs`, `post-target` when available. If
GitHub access is unavailable, include ready-to-create issue drafts.

## Final Gate

Before publishing, require:

- No unresolved `BLOCK`.
- No unresolved `FIX_BEFORE_TARGET` unless target/scope is delayed or rescoped.
- `DOC_BEFORE_TARGET` documented as a deliberate limitation or promoted to code
  fixes.
- Issue register covering fixed, documented, created/updated issues, and
  remaining post-target work.
- Targeted tests for every fixed bug.
- Full release verification after the final patch set.

Common mistakes:

- Only checking named risks: start with blind-spot discovery.
- Treating consensus as safety: evidence decides.
- Patching before verifying: prove the claim first.
- Downgrading docs hazards: unsafe examples can block adopters.
- Using release stage to avoid rigor: alpha can have sharp edges, but they must
  be explicit; beta/stable need stronger defaults.
