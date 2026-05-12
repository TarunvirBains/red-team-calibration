Findings below are only new items, sharper evidence, or no-action checks I think matter. I did not edit files or use delegation.

**Case D**

area: public API surface
finding: `MOD` is exported even though the contract exposes only `EngagementLoadScore`; this creates an accidental Go API.
evidence: [solution.go](solver/case-d/solution.go:3) exports `MOD`; the prompt names only the function signature and says implementation only in [prompt.md](solver/case-d/prompt.md:10).
repro/test idea: From another package, `cased.MOD` is addressable via `go doc`/compile.
impact: Low correctness impact, but public cleanup later becomes a breaking API change.
likelihood: High if published as a Go package.
fix size: Trivial: rename to `mod`.
publish decision: FIX_BEFORE_TARGET

**Case E**

area: scalability beyond staff count alone
finding: The exponential issue is worse than “large staff counts”: even the notes’ “m <= 20” assumption can be unusable when session count is realistic, because every grant subset reruns an `O(n*m)` matcher and allocates per subset.
evidence: subset loop at [solution.go](solver/case-e/solution.go:25), per-subset allocations at [solution.go](solver/case-e/solution.go:37), nested session/staff scan at [solution.go](solver/case-e/solution.go:49); notes claim `m <= 20` is practical in [notes.md](solver/case-e/notes.md:49).
repro/test idea: `20` staff, `10_000` sessions, `grantCount=20`, `grantSize=0`; expected answer is at most `20`, but the implementation attempts about `2^20 * 10_000 * 20` matching checks.
impact: Realistic scheduling sizes can hang even without “large” staff arrays.
likelihood: High under the unconstrained prompt contract.
fix size: Medium: binary search answer + greedy feasibility check.
publish decision: BLOCK

area: fixed-assignment matching
finding: No new correctness issue found once a grant assignment is fixed; sorted sessions with smallest sufficient unused capacity is a defensible greedy matcher.
evidence: [solution.go](solver/case-e/solution.go:49).
repro/test idea: Keep the known recommendation for small-reference cross-checks.
impact: The main blockers remain enumeration/resource limits, negative grant count, and overflow rather than this inner greedy step.
likelihood: NO_ACTION check.
fix size: None.
publish decision: NO_ACTION

**Case F**

area: delayed 2x2 test expectations
finding: Sharper evidence for the known bad tests: the two “unreachable” large-ready 2x2 cases have exact reachable answers, while the `readyAt=5` test expecting `6` is actually correct because destination parity is even.
evidence: tests expect `-1` at [solution_test.go](solver/case-f/solution_test.go:105) and [solution_test.go](solver/case-f/solution_test.go:114). In a 2x2 grid with non-destination cells ready at `0`, bouncing reaches destination at any even minute >= 2, so `readyAt=10` returns `10` and `readyAt=1_000_000_000` returns `1_000_000_000`.
repro/test idea: After replacing the BFS, assert `{{0,0},{0,10}} == 10`; do not run the billion case against the current implementation.
impact: Current tests would reject a correct parity-aware implementation.
likelihood: High.
fix size: Small for tests; medium for algorithm.
publish decision: BLOCK

area: no-wait edge coverage
finding: Add a single-row/single-column delayed-neighbor test so the eventual fix does not overgeneralize “bounce to wait” when no move is initially possible.
evidence: prompt forbids waiting in place at [prompt.md](solver/case-f/prompt.md:9); existing delayed tests use grids with cycles.
repro/test idea: `EarliestChecklistArrival([][]int{{0, 2}})` should be `-1`; `[][]int{{0,0,3}}` should be `4`.
impact: Prevents a common parity-fix regression.
likelihood: Medium during rewrite.
fix size: Small.
publish decision: FIX_BEFORE_TARGET