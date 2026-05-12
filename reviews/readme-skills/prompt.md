# README And Skills Review Prompt

Review these files only:

- `README.md`
- `skills/red-team-release-gate/SKILL.md`
- `skills/simplify-with-review/SKILL.md`

Review goals:

1. Natural prose: clarity, grammar, awkward phrasing, over-specific wording,
   missing context for a public reader.
2. Accuracy: verify the README matches this repository's actual contents,
   especially number of cases/contracts, reviewer outputs, workflows, and test
   caveats.
3. Language and stack neutrality: flag wording that implies this is Rust-only,
   Go-only, Codex-only, or tied to a private/local workflow when the idea could
   apply broadly.
4. Skill usability: check that each skill is actionable, not bloated, and does
   not contain private repo assumptions or local path assumptions.
5. Public repo readiness: flag confusing artifacts, missing caveats, missing
   license/readme context, or anything that would make the repo less useful to
   outside readers.

Do not edit files. Do not use delegation or nested review. Return concise
findings ordered by severity. For each finding include:

```text
severity:
file/section:
finding:
evidence:
suggested fix:
```

If there are no blocking issues, say so explicitly. Include a short
"positive signal" section only for useful context, not praise.
