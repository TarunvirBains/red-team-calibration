I have read all nine target files and the active prompts. Since plan-mode is active, my tools are read-only (Read/Glob/Grep), and the caller explicitly forbade delegation, nested reviewers, file edits, and command execution — I will deliver the red-team review directly. Note: my tool set does not include `Write`, `ExitPlanMode`, or `AskUserQuestion`, so I cannot persist a plan file or formally exit plan mode; the findings below are the review deliverable.

I traced each test against the implementation and re-derived expected values for every non-trivial assertion. The findings below are ordered by severity, novel issues first.

---

## Case C — `solver/case-c/solution.go`

```
area: parser — count tokenization
finding: parseNumber conflates "no digits present" with "explicit zero digit". After parseAtom and after `)`, the code does `if count == 0 { count = 1 }`. An explicit `H0` (or `(OH)0`) parses as count=0 from parseNumber, then is silently promoted to count=1 — so `H0` behaves as `H1`, not as zero/no-op/invalid.
evidence: solution.go:101-103 (atom branch) and lines 86-88 (group-close branch); parseNumber at lines 133-143 has no "found" signal — it returns the same `0` for both "no digit char" and "the digit '0'".
repro/test idea: assert NetInventory("H0") behavior; today the implementation returns "H". A faithful reading of the spec ("multi-digit counts are allowed", input is "valid") leaves zero counts ambiguous, but silent 1-promotion is a hidden trap. A test like {"H0","H"} would document current behavior; {"H0",""} or {"H0","INVALID"} would catch it.
impact: silent semantic drift on a parser primitive; if a future input or fuzzer drives a literal zero count anywhere in a formula or group multiplier, the result becomes wrong with no exception path.
likelihood: low under current "valid input" assumption; medium-to-high if any consumer extends inputs.
fix size: ~5 lines — return (n,bool) or use -1 sentinel; treat literal `0` as either invalid or kept-as-zero per spec clarification.
publish decision: FIX_BEFORE_TARGET
```

```
area: top-level expression parser — leading sign
finding: A leading `+` or `-` is silently accepted because the term-collector loop runs zero iterations on the first operator, then the operator branch updates `sign`. The spec says "a missing leading sign means addition" — it does not say a leading `+` or `-` is allowed. `+H2O` returns "H2O"; `-H2O` returns "INVALID". Behavior is undocumented either way and not under test.
evidence: solution.go:11-34 (no rejection of empty leading term); no test in solution_test.go targets a leading sign character.
repro/test idea: NetInventory("+H2O") and NetInventory("-H2O") both produce defined but unspecified results. Add a test that pins the chosen contract.
impact: ambiguous contract surface; downstream callers may rely on differing interpretations.
likelihood: medium for malformed inputs.
fix size: a few lines — either reject empty first term explicitly or document the leniency in notes.
publish decision: DOC_BEFORE_TARGET
```

```
area: parser — accepted-but-undefined characters
finding: parseFormula has a catch-all `else { i++ }` (solution.go:107-109) that silently skips any character that is not `(`, `)`, or uppercase ASCII. Whitespace, lowercase-leading garbage ("hello"), digits at the head of a sub-formula, and any other byte are silently dropped without comment.
evidence: solution.go:71-113. parseFormula has no error return.
repro/test idea: NetInventory("H2 O") today returns "H2O"; NetInventory("h2") returns "" silently. Add a test that pins or rejects this.
impact: typos and stray whitespace produce silently wrong inventories — a "valid input" contract that fails open is high-risk in audit/inventory contexts.
likelihood: medium for human-edited inputs.
fix size: small — return an error or panic on unexpected byte, or document tolerance.
publish decision: DOC_BEFORE_TARGET
```

```
area: parser — stack robustness
finding: parseFormula assumes well-formed parens. A leading `)`, or more `)` than `(`, performs `stack = stack[:len(stack)-1]` and then indexes `stack[len(stack)-1]` — runtime panic on out-of-range slice index.
evidence: solution.go:82-96. No depth check before pop.
repro/test idea: NetInventory(")H") panics. Spec says formulas are valid so the spec excuses this, but a misclassified top-level char could still drive it via the term boundary not handling embedded `+`/`-` inside parens (it can't, but a malformed input still reaches the parser).
impact: server-side panic on malformed input. For a release gate this is borderline.
likelihood: low under contract, but non-zero in production where inputs aren't always validated upstream.
fix size: tiny — guard on `len(stack) < 2` before popping and return an error/INVALID.
publish decision: POST_TARGET (unless inputs are not pre-validated, then FIX_BEFORE_TARGET)
```

```
area: performance — string building
finding: Three O(n²) string-concat hot paths: term assembly (`term += string(expr[i])` per byte, line 15), result assembly (`result += atom; result += s`, lines 53-65), and number-to-string (`s = char + s` prepend, lines 60-63).
evidence: solution.go:15, 53-65. No use of strings.Builder or strconv.Itoa.
repro/test idea: benchmark on a 10KB expression — each concat is O(current_len). Easy to demonstrate quadratic blow-up.
impact: a 100KB inventory term spends ~5e9 ops in string copying alone. Real for batch imports.
likelihood: high at scale, low at tiny test sizes.
fix size: small — strings.Builder + strconv.Itoa.
publish decision: POST_TARGET
```

```
area: notes — heading level mismatch with sibling cases
finding: case-c/notes.md uses level-1 (`# Design`, `# Correctness`, ...). case-a and case-b nest the required headings at level 2 under a level-1 title. The spec only requires the heading names, not the level, but the inconsistency across sibling cases is a visible calibration smell and could break tooling that pattern-matches by level.
evidence: solver/case-c/notes.md:1-29 vs solver/case-a/notes.md:1-39 and solver/case-b/notes.md:1-39.
repro/test idea: grep -E "^## (Design|Correctness|Complexity|Assumptions)$" — case-c fails the same recipe that matches case-a/b.
impact: minor; tool/CI lint risk only.
likelihood: low.
fix size: trivial.
publish decision: NO_ACTION (unless a pipeline depends on level)
```

---

## Case B — `solver/case-b/solution.go`

```
area: API design — value constructor over reference-type field
finding: NewEntitlementLedger returns EntitlementLedger by value, but the struct embeds a `map[int]int`. Maps are reference-headers in Go: copying the struct shares the underlying table. Methods are pointer-receiver, so `NewEntitlementLedger().Grant(0,5)` will not compile (return value is non-addressable), and `e2 := e1; e2.Grant(...)` silently mutates e1's view. This is exactly the prescribed prompt signature, but it is an API trap.
evidence: solution.go:7-11 (constructor returns value); solution.go:13, 22, 33 (pointer-receiver methods); no test asserts aliasing semantics.
repro/test idea: `e1 := NewEntitlementLedger(); e1.Grant(0,5); e2 := e1; e2.Revoke(0,5); want e1.Eligible(0,5,1)==true; got false` — proves the shared-table footgun. Also `NewEntitlementLedger().Grant(0,5)` is a compile error.
impact: production data integrity if a consumer ever copies the value (e.g., passing across a goroutine boundary, returning from a factory). Easy to mis-use.
likelihood: medium — common Go anti-pattern to copy a struct holding a map.
fix size: prompt-level — return `*EntitlementLedger`. As-implemented, conforms to the contract.
publish decision: DOC_BEFORE_TARGET (note the limitation in case-b/notes.md and explicitly state non-aliasing-safe semantics) and reconsider for next contract revision
```

```
area: memory / lifecycle
finding: Revoking entries to zero leaves them in the map. There is no compaction, and there is no upper bound on grant span. Long-running services drift toward map size = peak distinct granted IDs.
evidence: solution.go:22-31. No deletion or pruning.
repro/test idea: Grant(0, 1_000_000); Revoke(0, 1_000_000); inspect len(e.credits) — still 1_000_000.
impact: unbounded memory growth in long-running consumers; matches the "consumer-app scale" lens.
likelihood: medium for long-lived processes; low for ephemeral test harnesses.
fix size: small — delete entries when count drops to zero. Caveat: O(span) is intrinsic to the per-id model; a true production design should use a segment tree / interval-counter to avoid linear scan and per-id storage entirely.
publish decision: POST_TARGET
```

```
area: time complexity — Grant/Revoke/Eligible are O(span)
finding: Linear-per-id model is fine for small ranges but pathological for large half-open spans. The "long operation sequence" test exercises only spans ≤200 and 5 operations — it does not actually stress the linear model. The prompt asks for "long operation sequences"; the test name implies coverage that the body does not deliver.
evidence: solution.go:17-19, 26-30, 40-44; solution_test.go:158-190.
repro/test idea: e.Grant(0, 1<<24) timing — and contrast with a property-based stress that mixes Grant/Revoke/Eligible at random over a wide id space.
impact: misnamed test masks a known scalability ceiling.
likelihood: medium — the test name signals coverage that is not actually provided.
fix size: small — rename the test or add a real long-sequence stress (e.g., 10k mixed ops over 10k ids); or call out the O(span) limit in notes.
publish decision: DOC_BEFORE_TARGET
```

```
area: concurrency posture
finding: No declared concurrency contract on the type. Notes assert "no concurrent access" but neither method documents this, and Go maps panic on concurrent writes. Calling Eligible while another goroutine is mid-Grant will race.
evidence: solution.go (no doc comments, no mutex); notes.md:39 mentions sequential operations only.
repro/test idea: go test -race with two goroutines, one Grant-looping and one Eligible-checking — race detector fires.
impact: misuse is silent; data race or panic in production.
likelihood: medium for any multi-threaded consumer.
fix size: tiny doc, or add a sync.RWMutex (small).
publish decision: DOC_BEFORE_TARGET
```

```
area: input domain — pathological span
finding: No upper bound check on (right - left). A single Grant(0, math.MaxInt) hangs forever. Spec is silent on bounds, but the API has no defensive cap.
evidence: solution.go:13-20.
repro/test idea: Grant(0, 1<<40) on a service exposed to untrusted callers — DoS.
impact: trivial DoS if callers are not pre-trusted.
likelihood: low for an internal library, high if exposed.
fix size: small — bound or document caller-trust requirement.
publish decision: POST_TARGET
```

---

## Case A — `solver/case-a/solution.go`

```
area: arithmetic — sum overflow ceiling
finding: Sum uses int64 but readings are `[]int`, which on 64-bit platforms is also int64. The contract says "use 64-bit arithmetic for totals" but does not bound reading magnitude. A run of, e.g., int64-max-near readings will overflow int64 silently in the running sum. The notes assert "readings fit in int32" but the function signature does not.
evidence: solution.go:20 (`sum += int64(readings[j])`); notes.md:36 declares an int32 assumption that the signature does not enforce.
repro/test idea: AuditMeterReadings([]int{math.MaxInt64, 1}, 0, 1, 2) — current code overflows to a negative and may spuriously satisfy a "<= maxTotal" check.
impact: silent overflow can flip period validity. Auditing semantics suggest this is exactly the case the team cares about.
likelihood: low for realistic meter readings, but the contract carve-out is missing.
fix size: small — document the bound, or upgrade to math/big or saturating arithmetic if the audit really must accept the full int domain.
publish decision: DOC_BEFORE_TARGET
```

```
area: contract — inverted band
finding: No validation or documentation for `minTotal > maxTotal`. Current behavior: return 0 (no sum can satisfy both inequalities). This is probably the intended degenerate behavior, but it is not specified and not tested.
evidence: solution.go:23; tests cover legal bands only.
repro/test idea: AuditMeterReadings([]int{1,2}, 5, 1, 2) — assert 0 and pin the behavior.
impact: subtle — a caller passing inverted bounds (e.g., bug upstream) gets 0 with no signal.
likelihood: low.
fix size: trivial — one test + a sentence in notes.
publish decision: DOC_BEFORE_TARGET
```

```
area: time complexity — scalability ceiling
finding: O(n × maxReadings) per call. For n = 10^6 and maxReadings = 10^4 that is 10^10 — far beyond a single-threaded request budget. Sliding window + prefix-sum + ordered-structure ("count of prefix sums in a range") yields O(n log n). Notes acknowledge O(n²) without naming the scaling ceiling.
evidence: solution.go:11-27.
repro/test idea: benchmark on n = 100k with maxReadings = 50k.
impact: scalability ceiling not signposted.
likelihood: medium at production volumes.
fix size: prefix-sum + balanced-multiset or Fenwick — moderate.
publish decision: POST_TARGET (call it out in notes now)
```

```
area: tests — overlap & boundary spread
finding: Existing tests are correct (I re-derived every expected value) and cover the four contract-required cases. They do not include a property/fuzz pass and they do not include a case where `maxReadings == 1` interacts with bounds straddling zero, nor a case where the only valid period is at the right boundary of the slice. Low-priority gap.
evidence: solution_test.go.
repro/test idea: add a property test that enumerates all O(n²) subarrays of a small random array via a reference brute force, and compares against the implementation.
impact: low — code is simple enough that I do not expect bugs.
likelihood: low.
fix size: small.
publish decision: NO_ACTION (POST_TARGET if fuzz/property tests are within the calibration's testing scope)
```

---

## Cross-cutting

```
area: package documentation
finding: None of the three solution.go files carry package-level doc comments and none of the exported identifiers (AuditMeterReadings, EntitlementLedger, NewEntitlementLedger, Grant, Revoke, Eligible, NetInventory) have GoDoc. golint/staticcheck-style review at a "senior at large tech co" level expects these.
evidence: solver/case-{a,b,c}/solution.go top of file and per-function.
repro/test idea: `go doc ./solver/case-a` returns no docs.
impact: missing API surface documentation; calibration smell at the "production-grade Go" bar.
likelihood: not a behavior bug.
fix size: small.
publish decision: DOC_BEFORE_TARGET
```

```
area: test style consistency
finding: case-a and case-c use t.Run table-driven subtests; case-b uses many separate Test* functions, no subtests, and assertions via t.Error rather than t.Fatalf, so a single mis-set credit yields a cascade of confusable errors.
evidence: solver/case-b/solution_test.go.
repro/test idea: invert one assertion mid-file and observe the noise downstream tests produce.
impact: developer-experience inconsistency in a calibration corpus.
likelihood: low.
fix size: small.
publish decision: NO_ACTION (or POST_TARGET if calibration evaluates test style across cases)
```

```
area: notes coverage of failure modes
finding: None of the three notes.md files spell out the known limitations I have flagged above (overflow ceiling for A; aliasing-on-copy + memory growth + concurrency for B; zero-count promotion + whitespace tolerance + leading-sign tolerance for C). For a release-gate review, the right place to land these is the notes' "Assumptions" section.
evidence: solver/case-{a,b,c}/notes.md "Assumptions" headings.
repro/test idea: read each notes.md and check whether each "Assumption" actually constrains the API or just restates the prompt.
impact: false sense of completeness in the calibration corpus.
likelihood: medium — calibration tooling that reads notes will not see the gaps.
fix size: small per case.
publish decision: DOC_BEFORE_TARGET
```

---

## Summary of publish decisions

- **FIX_BEFORE_TARGET**: 1 — Case C parseNumber zero/missing conflation.
- **DOC_BEFORE_TARGET**: 7 — Case A overflow + inverted band, Case B value-constructor aliasing + concurrency + test-name overstatement, Case C leading-sign + whitespace tolerance, cross-cutting package docs, notes coverage of failure modes.
- **POST_TARGET**: 5 — Case A scalability ceiling + property tests, Case B map compaction + DoS bound, Case C performance and parser robustness on malformed parens.
- **NO_ACTION**: 2 — notes heading level (unless tooling cares), test-style consistency (calibration-tooling dependent).
- **BLOCK**: 0 — none of the implementations contradict their active prompt contract in a way that hard-breaks the stated behavior. The single fix-before item is a parser primitive that silently misinterprets a literal input the prompt does not explicitly forbid; the rest are documentation/scalability/API-design surfaces.

The single most likely-to-bite issue in production-shaped use is the Case B value-receiver constructor over a reference-type field: it conforms to the prompt but produces a footgun a senior reviewer at a large org would not let ship without explicit doc and a usage example.
