I've read the three target files and verified the README's specific claims against the repository contents. A constraint I should surface up front: plan mode tells me to write the plan to `$HOME/.claude/plans/readme-and-skills-zippy-lampson.md` and to call `ExitPlanMode`, but the tools available to me in this session are read-only (`Read`, `Glob`, `Grep` only — no `Write`, `Edit`, or `ExitPlanMode`). I therefore cannot create the plan file or formally exit plan mode. Since the user's request is itself a read-only review and explicitly forbids edits and delegation, I'm delivering the review findings directly in this response.

# Review: README and Skills

Ordered by severity. No findings block public publication; Findings 1–3 are the ones I'd fix before treating this repo as a polished public reference.

---

## HIGH

```text
severity: high
file/section: README.md — Workflow step 5
finding: Documented reviewer output path does not match the repository.
evidence: README says "Save reviewer outputs as Markdown under `reviews/<reviewer>/cases/<case>.md`." Actual layout is one combined file per reviewer per round: `reviews/<reviewer>/cases/full-panel.md` and `reviews/<reviewer>/round-2/full-panel.md`. No `case-a.md` / `case-b.md` style files exist anywhere under `reviews/`.
suggested fix: Replace step 5 with: "Save reviewer outputs as `reviews/<reviewer>/cases/full-panel.md` (one file containing all cases for that round; round 2 outputs live under `reviews/<reviewer>/round-2/full-panel.md`)."
```

```text
severity: high
file/section: README.md — Workflow step 4 vs. "Active Contracts"
finding: "Active Contracts" lists six cases (A–F) but Workflow step 4 says "Ask each reviewer to review all three solver implementations." The repo actually runs two rounds — round 1 over A/B/C, round 2 over D/E/F — and the README never mentions rounds. A public reader sees six contracts above a workflow that names three implementations and has no signpost for the gap.
evidence: README:38 ("all three solver implementations"); `reviews/panel-prompt.md` lists only A/B/C; `reviews/round-2/panel-prompt.md` lists only D/E/F; `synthesis/round-2-review-matrix.md` records round-2 reviewer status. The README contains no occurrence of "round 2" or "rounds".
suggested fix: Either (a) describe the workflow as two rounds of three cases each and link `reviews/round-2/panel-prompt.md` and `synthesis/round-2-review-matrix.md`, or (b) reword step 4 to "all solver implementations in scope for the current panel run" so the count is not pinned.
```

```text
severity: high
file/section: skills/red-team-release-gate/SKILL.md (lines ~154, 175, 208, 214, 281); skills/simplify-with-review/SKILL.md (lines ~22–32, 106, 113)
finding: Skills leak author-private tools and agent names, presenting them as canonical workflow rather than examples. `rtk timeout` (a personal CLI proxy), `careful-coder` (a user-private Claude Code agent under `~/.claude/agents/`), and Codex-specific identifiers `multi_tool_use.parallel`, `apply_patch`, and `spawn_agent` all appear in the main prose. The Standing Reviewer Panel pins six specific models (Gemini 3.1 Pro, Kimi K2.6, Claude Opus 4.7 with `careful-coder`, DeepSeek V4, Codex 5.5, GLM 5.1) as the panel, not as an example. A public reader cannot install `rtk` or `careful-coder`, and Codex tool names do not apply to other CLIs.
evidence: red-team SKILL.md:208 ("Wrap reviewer calls with `rtk timeout`"); :175 and :281 (`claude -p --agent careful-coder`); :214 (`rtk timeout 3600s opencode run ...`); simplify SKILL.md "Codex Rules" section names `multi_tool_use.parallel`, `apply_patch`, `spawn_agent`; :106 names `careful-coder`.
suggested fix: Move tool- and model-specific shapes into an "Example commands" appendix at the bottom of each skill. In the body, name the role and capability (e.g., "max-effort, read-only, no nested delegation") so the skill works with any CLI. Replace `rtk timeout` with neutral wording like "wrap reviewer calls with a per-process timeout, e.g. `timeout 1800s ...`". Replace `careful-coder` mentions with "a read-only review agent (e.g., …)" or note explicitly that this is the author's local agent.
```

---

## MEDIUM

```text
severity: medium
file/section: README.md — intro paragraph
finding: The "low-reasoning Haiku implementer passes" claim is not verifiable from anything checked into the repo. No solver prompt names a model, and there is no implementer manifest, command log, or config recording which model produced each `solution.go`. "Haiku" also assumes the reader knows it refers to Anthropic Claude Haiku.
evidence: `solver/case-a/prompt.md` and the other prompt.md files contain only "You are a constrained implementation agent" with no model name; no implementer-side log exists under `solver/` or `reviews/`.
suggested fix: Either (a) check in a brief "How the solver outputs were produced" note listing model, provider, reasoning effort, and dispatch command, or (b) reword to "intentionally constrained low-reasoning implementer passes (a small-tier LLM at low reasoning effort)" so the claim is generic and reproducible.
```

```text
severity: medium
file/section: README.md — "Reviewers" section
finding: README says reviewers compare "the distinct review signal each fixed model slot produces" but never names the slots. A new reader must glob `reviews/` to discover the panel. The skill files list the panel explicitly; the README, which is the project's front door, does not.
evidence: README:14–17; `reviews/` contains directories for `gemini-3.1-pro`, `kimi-k2.6`, `claude-opus-4.7`, `deepseek-v4`, `codex-5.5`, `glm-5.1`; `synthesis/round-2-review-matrix.md` lists the canonical names.
suggested fix: Add a short bulleted list under "Reviewers" enumerating the six slots (model + provider/CLI), or link to `synthesis/round-2-review-matrix.md` for the authoritative table.
```

```text
severity: medium
file/section: README.md — Workflow step 5
finding: The README cautions that solver outputs may be intentionally flawed, but gives no equivalent caveat for reviewer outputs. The repo ships partial/empty reviewer artifacts: `reviews/readme-skills/opus.md` is 1 line; `reviews/deepseek-v4/round-2/full-panel.err`, `full-panel-retry.err`, `case-f-narrow.err`, `case-f-narrow-no-novita.err`, `case-f-narrow-no-novita-autonomous.err` reflect a 502 broad-run failure followed by narrowed retries. A reader doing a naive size or "first finding" comparison across reviewers will be misled.
evidence: As cited; `synthesis/round-2-review-matrix.md` already records the DeepSeek 502/narrow-retry chain, so the caveat is already true in practice.
suggested fix: Add one line to step 5 or a new "Caveats" subsection: "Some reviewer artifacts are empty or partial (provider 502s, timeouts, narrow retries). See `synthesis/round-2-review-matrix.md` for the canonical reviewer status table."
```

```text
severity: medium
file/section: README.md — intro paragraph
finding: "Go library implementations" overstates what each case is. Each case is a single exported function plus a small `_test.go` in its own package (e.g., `solver/case-a/solution.go` is a 30-line single function). "Library" is grander than the artifact.
evidence: `solver/case-a/solution.go` (30 lines, one function); other cases are similar.
suggested fix: Replace "Go library implementations" with "small single-function Go packages" or "small Go solver packages". (Note: the solver prompts themselves use "small Go library", so consider aligning both.)
```

```text
severity: medium
file/section: skills/red-team-release-gate/SKILL.md — frontmatter description
finding: The description's artifact list leads with `crate`, which is Rust-specific. The skill body is otherwise mostly language-neutral, but front-loading a Rust term signals (incorrectly) that the skill is Rust-only.
evidence: SKILL.md:3 — "publish gate for a crate, package, service, CLI, SDK, library, or app".
suggested fix: Drop `crate` (it is already covered by "package" / "library"), or replace with "module".
```

---

## LOW

```text
severity: low
file/section: README.md — Workflow steps 6, 7
finding: The synthesis templates the workflow points at are stale relative to the current case/reviewer counts. `synthesis/cross-review-matrix.md` only has rows for cases A/B/C (no D/E/F), and `synthesis/panel-instruction-tuning.md` only has sections for DeepSeek V4, Kimi K2.6, and GLM 5.1 — missing Gemini, Claude Opus, and Codex. A reader following the workflow finds templates that no longer cover the data.
evidence: `synthesis/cross-review-matrix.md` lines 5–22; `synthesis/panel-instruction-tuning.md` lists three reviewer sections; `synthesis/round-2-review-matrix.md` covers cases D–F separately.
suggested fix: Expand the templates to cover all six cases and all six reviewers, or replace steps 6–7 with a single pointer to `synthesis/round-2-review-matrix.md` plus a "fill in your own matrix" template stub.
```

```text
severity: low
file/section: README.md — "Skills" section
finding: "Scrubbed copies" is opaque to a public reader and inconsistent with what's in the files: several private identifiers (`rtk`, `careful-coder`, Codex-specific tool names) remain (see Finding 3).
evidence: README:57; SKILL.md content as cited in Finding 3.
suggested fix: Drop "scrubbed" and say what they are: "checked-in copies of two process skills used during this calibration." If further scrubbing is intended, do it (per Finding 3) before relanding this wording.
```

```text
severity: low
file/section: README.md
finding: The README does not mention `reviews/readme-skills/`, which contains the meta-review prompt and outputs reviewing the README and skills themselves. A public reader has no way to discover this self-referential artifact.
evidence: `reviews/readme-skills/{prompt.md, opus.md, gemini.md, *.err}` exist.
suggested fix: One line in the "Skills" section (or a new "Meta-review" line): "The README and skills themselves were reviewed; see `reviews/readme-skills/prompt.md` and the reviewer outputs alongside it."
```

```text
severity: low
file/section: README.md — Workflow step 3
finding: The caveat about test failures is vague where it could be specific. `synthesis/round-2-review-matrix.md` already records a concrete failure: Case F's package tests exit with code 124 under an 8s timeout. The README hedges with "may fail or time out" instead of naming the case.
evidence: `synthesis/round-2-review-matrix.md` lines 6–9.
suggested fix: Replace the hedge with the concrete pointer, e.g.: "Some checked-in solver outputs are intentionally flawed; in particular, `solver/case-f` does not finish its package tests under an 8s timeout. See `synthesis/round-2-review-matrix.md`."
```

```text
severity: low
file/section: README.md "Reviewers"; skills/red-team-release-gate/SKILL.md depth template (line ~256)
finding: The phrasing "senior software engineer at a large global/multinational technology company" reads as the author's house style rather than a generic prompt frame. Not wrong; just over-specific for a public skill.
evidence: README:16–17; red-team SKILL.md:256 (and the same phrasing recurs in simplify SKILL.md:132–133).
suggested fix: Optional — generalize to "an experienced senior software engineer".
```

```text
severity: low
file/section: skills/red-team-release-gate/SKILL.md — Standing Reviewer Panel (lines 146–168)
finding: The numbered panel ends with "Codex is not the final slot; the orchestrator adjudicates all outputs regardless of reviewer order." This contradicts (or at least rephrases away) the numbered list immediately above, where Codex is #5 and GLM is #6. Reads as backfilling.
evidence: SKILL.md:146–168.
suggested fix: Either reorder the list so Codex is not numbered as a fixed slot, or replace the panel with a role-based list (Architectural Critic, Edge-Case Provocateur, Logic & Security Auditor, etc.) and put specific model pins in an example block. This also addresses Finding 3.
```

```text
severity: low
file/section: skills/red-team-release-gate/SKILL.md — line 26
finding: "Bind the target before review" — "bind" is jargon in this context.
suggested fix: "Set the target before review" or "Choose the release target before review".
```

---

## Positive signal (context, not praise)

- The README's "Evaluation Rule" bullets are the most reusable part of the public artifact — they define what makes one reviewer's output more valuable than another's in concrete terms.
- `synthesis/round-2-review-matrix.md` is an honest status artifact (records provider 502s, narrow retries, timeouts). Surfacing it from the README is a low-cost win and naturally fixes Findings 2, 5, 6, and 9 at once.
- `simplify-with-review`'s workflow shape (Scope → Chunk → Read → Verify → Review → Fix, with explicit "Stop Conditions" and "Anti-patterns") generalizes well beyond Codex once the tool-specific names move into an examples appendix.

---

## Blockers

None. Nothing in these three files prevents public release. The high-severity findings are accuracy and private-identifier hygiene; the medium findings are reader-orientation gaps; the low findings are polish.
