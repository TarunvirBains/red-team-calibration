**Case A**

area: Case A arithmetic correctness
finding: Running sums can silently overflow `int64`, producing false counts for valid `int` inputs near the 64-bit boundary.
evidence: [solution.go](solver/case-a/solution.go:20) does unchecked `sum += int64(readings[j])`; [notes.md](solver/case-a/notes.md:36) claims 64-bit arithmetic avoids overflow.
repro/test idea: `readings := []int{math.MaxInt64, 1}`, band `[math.MinInt64, math.MinInt64]`, `maxReadings := 2`; wrapped arithmetic can count a period whose mathematical sum is outside the band.
impact: Incorrect audit counts on boundary inputs; silent false positives are worse than rejection.
likelihood: Low to medium, depending on input bounds.
fix size: Small to medium.
publish decision: FIX_BEFORE_TARGET

area: Case A notes contract
finding: `notes.md` adds an extra heading not allowed by the prompt’s required heading set.
evidence: Prompt requires headings `Design`, `Correctness`, `Complexity`, and `Assumptions` at [prompt.md](solver/case-a/prompt.md:32); [notes.md](solver/case-a/notes.md:1) starts with `Notes: Case A`.
repro/test idea: Parse Markdown headings and compare exact heading text to the required set.
impact: Contract-format failure for release/grading even if code is correct.
likelihood: High.
fix size: Tiny.
publish decision: DOC_BEFORE_TARGET

**Case B**

area: Case B scalability / resource exhaustion
finding: The ledger materializes every account ID in a span, so large valid spans can consume unbounded CPU and memory.
evidence: [solution.go](solver/case-b/solution.go:4) stores per-ID credits, and `Grant`, `Revoke`, and `Eligible` loop over every ID at [solution.go](solver/case-b/solution.go:17), [solution.go](solver/case-b/solution.go:26), and [solution.go](solver/case-b/solution.go:40).
repro/test idea: `e.Grant(0, 1_000_000_000)` or `e.Eligible(0, 1_000_000_000, 1)` should demonstrate impractical runtime/memory.
impact: Trivial denial of service for a span-based customer platform API.
likelihood: Medium; the contract gives no maximum span size.
fix size: Medium to large, likely interval/segment representation.
publish decision: FIX_BEFORE_TARGET

area: Case B API safety
finding: The exported zero value panics on `Grant` because the map is only initialized by the constructor.
evidence: Constructor initializes `credits` at [solution.go](solver/case-b/solution.go:7), but `Grant` writes without a nil-map guard at [solution.go](solver/case-b/solution.go:18).
repro/test idea: `var e EntitlementLedger; e.Grant(0, 1)` panics with assignment to nil map.
impact: Common Go zero-value use crashes instead of behaving as an empty ledger.
likelihood: Medium.
fix size: Tiny.
publish decision: FIX_BEFORE_TARGET

area: Case B notes contract
finding: `notes.md` violates the required heading contract and adds an unsupported non-negative account-ID assumption.
evidence: Prompt requires exact notes headings at [prompt.md](solver/case-b/prompt.md:37); [notes.md](solver/case-b/notes.md:1) adds `Case B: Entitlement Ledger`, and [notes.md](solver/case-b/notes.md:36) assumes non-negative IDs though the API accepts `int` spans.
repro/test idea: Heading parser plus a behavior test such as `Grant(-2, 1)` to show negative IDs are accepted.
impact: Contract-format failure and ambiguous API boundary.
likelihood: High for heading failure; medium for API confusion.
fix size: Tiny.
publish decision: DOC_BEFORE_TARGET

**Case C**

area: Case C numeric correctness
finding: Formula counts and multipliers use unchecked `int` arithmetic, so large multi-digit counts can overflow into wrong output or false `"INVALID"`.
evidence: Counts are parsed into `int` at [solution.go](solver/case-c/solution.go:137), updated unchecked at [solution.go](solver/case-c/solution.go:139), and multiplied unchecked at [solution.go](solver/case-c/solution.go:95); the prompt allows multi-digit counts at [prompt.md](solver/case-c/prompt.md:30) without a bound.
repro/test idea: `NetInventory("H18446744073709551617")` on a 64-bit platform can wrap to a tiny count instead of returning the canonical huge count.
impact: Silent inventory corruption on large but syntactically valid quantities.
likelihood: Medium for imported/untrusted data.
fix size: Medium.
publish decision: FIX_BEFORE_TARGET

area: Case C zero-count semantics
finding: Explicit `0` counts or group multipliers are treated as missing counts and converted to `1`.
evidence: Both group and atom parsing replace parsed `0` with `1` at [solution.go](solver/case-c/solution.go:85) and [solution.go](solver/case-c/solution.go:100).
repro/test idea: If zero counts are valid, `NetInventory("H0")` should return `""`, but this implementation returns `"H"`.
impact: Silent overcounting unless the contract explicitly excludes zero multipliers/counts.
likelihood: Medium because the prompt does not state counts must be positive.
fix size: Small.
publish decision: FIX_BEFORE_TARGET

I did not edit files or run tests; this was direct static review of the listed active solver files and prompt contracts only.