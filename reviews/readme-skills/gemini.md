Positive signal:
The README clearly explains the intent of the calibration harness, the cases are well-documented, and the directory structure matches the stated goals. The skills present detailed, rigorous review methodologies that are conceptually strong for complex code reviews.

```text
severity: high
file/section: skills/simplify-with-review/SKILL.md / Codex Rules
finding: The skill is hardcoded to assume a "Codex agent" orchestrator and relies on Codex-specific proprietary tool names, breaking agent neutrality and public usability.
evidence: Mentions of "The main Codex agent", `multi_tool_use.parallel`, `spawn_agent`, and `apply_patch`.
suggested fix: Generalize the language to be agent-agnostic (e.g., "The orchestrating agent") and describe the required actions conceptually rather than using specific tool API names.

severity: high
file/section: skills/red-team-release-gate/SKILL.md & skills/simplify-with-review/SKILL.md / Pin And Routing Checks & Panel Discipline
finding: Both skills heavily depend on private/local workflows, proprietary orchestration CLIs, and specific model versions, making them unusable as drop-in skills for public readers.
evidence: Hardcoded commands using `rtk timeout`, `opencode run`, `claude -p`, `gemini --model`, and `codex exec`. The strict "Standing Reviewer Panel" mandates specific models like "gpt-5.5", "Gemini 3.1 Pro", and "Kimi K2.6".
suggested fix: Abstract the review panel into required roles or capabilities. Remove private CLI tool references (`rtk`, `opencode`) and provide the exact CLI invocation commands as optional examples rather than mandatory steps.

severity: medium
file/section: README.md / Workflow
finding: The documented path for saved reviewer outputs is inaccurate and does not match the actual repository structure.
evidence: The README states "Save reviewer outputs as Markdown under `reviews/<reviewer>/cases/<case>.md`." However, the repository structure saves outputs into consolidated files like `reviews/<reviewer>/cases/full-panel.md`.
suggested fix: Update step 5 in the Workflow section to reflect the actual file naming convention: "Save reviewer outputs as Markdown under `reviews/<reviewer>/cases/full-panel.md` (or similar consolidated files)."

severity: low
file/section: skills/red-team-release-gate/SKILL.md / Pin And Routing Checks
finding: The instructions for OpenRouter API preflight checks and routing are bloated and over-specific for a general skill definition.
evidence: Detailed HTTP request sequences like `GET /api/v1/models/$OR_MODEL/endpoints`, canary POST requests, and `GET /api/v1/generation?id=$GENERATION_ID`.
suggested fix: Simplify the routing checks to focus on the core principle of "verify model availability and fallback behavior before dispatching" rather than dictating exact internal API payloads and shell workflows.
```
