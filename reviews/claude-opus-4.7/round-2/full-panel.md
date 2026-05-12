I have the full picture. Conducting the three-pass review now and reporting only novel findings, material disagreements, sharper repros, and important areas where I checked and found no issue.

---

# Release-Gate Red-Team Review — solver/{case-d,case-e,case-f}

Scope: `solver/case-d`, `solver/case-e`, `solver/case-f`. I read each `prompt.md`, `solution.go`, `solution_test.go`, `notes.md` end to end. I will not repeat known issues except where I have sharper evidence or severity disagreement.

## Case D — `EngagementLoadScore`

### Novel finding D-1: Tests encode the wrong contract interpretation; fixing the implementation will cascade-break the test suite

```text
area: tests, contract conformance
finding: Every non-trivial test in solution_test.go computes its `expected`
  value under the same wrong interpretation as the implementation
  ("same-value consecutive runs"). When the implementation is fixed to score
  every continuous span, all of these tests will fail simultaneously.
  TestEngagementLoadScore_AlternatingValues expects 11 for [1,2,1,2,1];
  the contract-correct answer is 60 (sum of min*sum over all 15 subarrays).
  TestEngagementLoadScore_MixedHighsAndLows expects 81 for [5,5,3,3,3,2];
  contract-correct value is much larger. Repeated equal-value tests
  (TestEngagementLoadScore_RepeatedEqualValues expects 27 for [3,3,3])
  similarly encode the wrong rule.
evidence: solution_test.go lines 58-66, 70-80, 117-132. The notes.md at
  lines 5 and 15-19 explicitly redefine a "continuous run" as "a maximal
  sequence of consecutive days with identical daily loads" — a private
  re-interpretation of the prompt, which says "every continuous run of days
  in a plan".
repro/test idea: Replace one test with the contract-canonical example
  [1,2] → 8 from the known-issue list. Add [1,2,1,2,1] → 60. Re-run.
impact: Test suite is anti-helpful: a green run today implies bug-correct
  behavior, not contract-correct behavior. The BLOCK fix is therefore larger
  than "change the algorithm"; it includes a complete test rewrite. Without
  flagging this, downstream agents may swap in a correct algorithm, see
  every test fail, and incorrectly assume their algorithm is broken.
likelihood: certain
fix size: medium (rewrite test expectations; ideally add a brute-force
  reference to derive them, since hand-derivation for non-trivial cases is
  error-prone)
publish decision: BLOCK
```

### Novel finding D-2: `notes.md` documents the wrong contract under a publicly correct heading

```text
area: docs, maintainability
finding: notes.md Section 1 ("Design") declares a definition of "continuous
  run" that does not appear in the prompt. The prompt's contract text is
  "every continuous run of days in a plan"; the notes silently narrow it to
  "maximal sequence of consecutive days with identical daily loads". The
  notes also justify this with a fake worked example.
evidence: notes.md lines 5, 15-19, 26.
repro/test idea: Diff prompt.md against notes.md; the words "identical
  daily loads" appear only in the notes.
impact: When this BLOCK is fixed, the notes must be rewritten or future
  reviewers will re-anchor on this incorrect definition and re-introduce
  the bug. Treat the doc as part of the contract-violation surface.
likelihood: certain
fix size: small (re-derive Design / Correctness against the real contract)
publish decision: BLOCK (gated together with the algorithm/test fixes)
```

### Material commentary on known issue D (algorithm)

```text
area: algorithm
finding: I concur with the BLOCK on the algorithm. Adding evidence: the
  classical contract interpretation is the standard "sum of (min × sum)
  over all subarrays" problem, which has a known O(n) monotonic-stack
  solution. The current implementation collapses entire subarray families
  into a single run-length contribution, so the gap grows as O(n^2) in the
  number of distinct neighboring values, not a constant factor.
evidence: For sorted input [1,2,3,...,n] of length n, the wrong code
  produces sum of i*i = n(n+1)(2n+1)/6; the correct value grows ~n^4/24.
  At n=10 this is 385 vs 7,260 — a ratio of ~19x at moderate sizes.
repro/test idea: Property test: brute-force over all subarrays vs current
  output, fail loudly on the first non-monotone input.
impact: same as known.
likelihood: certain
fix size: small (monotonic-stack rewrite is ~30 lines and well-known)
publish decision: BLOCK (no change in severity, sharper evidence only)
```

---

## Case E — `MaxCoveredSessions`

### Novel finding E-1: `1 << uint(len(staff))` silently returns 0 for ≥64 staff (silent correctness bug distinct from known runtime bug)

```text
area: algorithm, numerical safety
finding: The outer loop bound is `mask < (1 << uint(len(staff)))`. On
  64-bit systems (Go `int` is 64-bit), once `len(staff) >= 64` the shift
  evaluates to 0 (the bit shifts out of range), the loop body never
  executes, and the function returns the initial `maxCovered = 0`. The
  function silently returns "no sessions can be covered" for inputs that
  are obviously coverable. This is qualitatively different from the
  known "exponential hang" bug: it is a fast, silent, wrong answer.
evidence: solution.go line 25 uses `(1 << uint(len(staff)))` with no
  guard. Go spec: shift operations on signed types may overflow with
  defined modular wrap; `1 << 64` on int64 is 0.
repro/test idea:
  staff := make([]int, 64); for i := range staff { staff[i] = 100 }
  sessions := []int{10}
  MaxCoveredSessions(sessions, staff, 0, 0)  // returns 0, contract says 1
impact: A user with a large staff list gets "zero coverage" results that
  silently disagree with reality. Worse, this happens at the exact size
  where the runtime issue would also be observed; the failure mode shifts
  from "hangs" to "wrong answer in microseconds", which is harder to
  detect.
likelihood: certain (purely deterministic)
fix size: small as a side effect of the proper algorithm rewrite (a real
  fix replaces the 2^m loop entirely with a polynomial matching strategy)
publish decision: BLOCK (additive to known runtime issue, not a duplicate)
```

### Novel finding E-2: Tests that claim to validate the matching strategy do not actually exercise it

```text
area: tests
finding: TestLocalGreedyOrder and the "strategic grant on higher capacity
  staff" case in TestComplexScenarios announce themselves as guards
  against locally-attractive pairings, but every expected value in those
  tests is achievable by *any* valid mask + matching choice (including
  unsorted naive pairing). They do not falsify a wrong-greedy
  implementation; they only verify the trivial coverage count.
evidence:
  - TestLocalGreedyOrder/"greedy matching optimal with unsorted input":
    sessions=[1,5,5], staff=[1,5], grantCount=1, grantSize=4. Any sensible
    matching achieves 2; a backwards-greedy that matches the largest
    session first against the largest staff would also score 2.
  - TestLocalGreedyOrder/"sorting ensures good matching over naive
    pairing": sessions=[10,10,1], staff=[11,5], grantCount=1, grantSize=6.
    With grant on staff[0]=11 → 17, can cover one 10 + the 1, score 2.
    Without grant, staff=11 covers 10, staff=5 covers 1, score 2.
    Without sorting, output is also 2.
  - TestComplexScenarios/"strategic grant on higher capacity staff":
    expected = 1; achievable trivially.
repro/test idea: Construct a case where naive (input-order) greedy
  matches one session, sorted-smallest-fit matches two. Example:
  sessions=[1, 4], staff=[2, 4], no grants. Naive: session 1 → staff 2
  (capacity 2, ok), session 4 → staff 4 (ok) → 2. Sorted-fit: session 1
  matched to staff 2 (smallest fit), session 4 → staff 4, → 2. Both 2.
  A genuine differentiating case is harder than the test name suggests;
  the test author should construct one or remove the misleading label.
impact: When the implementation is rewritten away from 2^m + greedy
  (e.g., to a Hungarian / sorted two-pointer approach), the new tests
  will pass but the *kind* of bug these tests pretend to guard against
  remains unverified.
likelihood: high (the rewrite path is now likely)
fix size: small (add 2–3 distinguishing tests, ideally generated against
  a brute-force reference)
publish decision: FIX_BEFORE_TARGET
```

### Novel finding E-3: No brute-force property test or randomized cross-check

```text
area: tests, validation
finding: For a "try all subsets" implementation that depends on the
  subsequent greedy being optimal for the fixed augmented capacities,
  the natural confidence-builder is a randomized property test that
  cross-checks against an O(2^m · m!) brute force on small inputs (say
  m,n ≤ 6). None is present.
evidence: solution_test.go contains only hand-rolled cases.
repro/test idea: For 200 random (sessions, staff, grantCount, grantSize)
  with bounded sizes, compare MaxCoveredSessions to a brute force that
  enumerates all permutations and all grant assignments.
impact: After any algorithm rewrite (which this case will need), there
  is no scaffolding to detect regressions. This is a release-gate quality
  concern even after the BLOCK is resolved.
likelihood: high
fix size: small
publish decision: FIX_BEFORE_TARGET
```

### Material commentary on known issue E (exponential enumeration)

```text
area: algorithm
finding: I concur with the BLOCK. Sharper evidence: the contract sets no
  bound on staff count; the notes admit "m ≤ 20-25" — a self-acknowledged
  contract violation in writing. The correct approach is polynomial: sort
  both arrays, decide for each candidate grant placement using a
  two-pointer matching over the augmented capacities, and pick the best
  among at most min(grantCount, m) sorted placements. A clean version is
  ~40 lines and runs in O((n+m) log (n+m) + m · grantCount).
evidence: notes.md lines 39-49, 55-61. solution.go line 25.
repro/test idea: m = 30, sessions of length 30, grantCount = 5; the
  current code attempts 2^30 ≈ 1e9 iterations and either hangs or OOMs.
impact: as known.
likelihood: certain.
fix size: medium (polynomial rewrite + new tests).
publish decision: BLOCK (no severity change; mention coupling with E-1).
```

---

## Case F — `EarliestChecklistArrival`

### Novel finding F-1: TWO existing tests assert results that are inconsistent with the implementation — the tests are wrong, and the suite cannot have been run to completion

```text
area: tests, correctness
finding: TestLargeReadyAtReachable (readyAt = [[0,0],[0,10]], expects -1)
  and TestLargeReadyAtUnreachable (readyAt = [[0,0],[0,1000000000]],
  expects -1) both assert that the destination is unreachable. The
  implementation will, given enough time and memory, reach the destination
  in both cases via legal oscillation between (0,0), (0,1), and (1,0).
  For the first test the implementation will return 10, not -1. For the
  second it will return 1_000_000_000 (or OOM before doing so) — never -1.
evidence:
  - readyAt[[0,0],[0,10]]: from (0,0) the worker can oscillate (0,0)↔(0,1)
    or (0,0)↔(1,0) while time increases by 1 per move. At t=9 the worker
    can stand at (0,1) (readyAt 0); at t=10 it can enter (1,1) since
    readyAt[1][1] = 10 ≤ 10. The implementation's BFS visits states keyed
    by (row, col, time), and will enqueue these revisit states because
    they are not in `visited` yet. Return value = 10.
  - readyAt[[0,0],[0,1e9]]: same logic with K=1e9; return = 1e9 if memory
    allowed, or OOM before returning. Cannot equal -1 unless killed
    externally.
  - The implementation's own TestArrivingLaterUnblocks and
    TestDestinationBlocked rely on exactly the same oscillation logic to
    pass with positive return values. The contradiction is internal to
    the test file.
repro/test idea:
  - `go test -run TestLargeReadyAtReachable -timeout 30s` against the
    current code: this will fail with `result = 10` (or hang then time
    out, but the value-comparison is the unambiguous failure).
  - For the 1e9 case, run with `GOMEMLIMIT=512MiB` and observe OOM-kill
    or runtime panic.
impact: This is a release-gate red flag. Either (a) `go test ./...` has
  not been executed against this case since the tests were written, or
  (b) the agent observed a crash and recorded it as "expected -1". Both
  interpretations imply the case-f deliverable was never validated to its
  own claimed standard. The notes.md "Edge cases handled: Large readyAt
  values don't cause overflow" claim (line 36) is unsupported.
likelihood: certain
fix size: medium — these tests need redesign. The simplest correct test
  for "very large readyAt unreachable" requires constructing a layout
  where the destination is truly unreachable (e.g., a wall of cells with
  readyAt so large that the destination's readyAt is *less* than them and
  the only physical paths are gated through them). Alternative: split
  into two test classes — "reachable but expensive" (expected K) and
  "structurally unreachable" (expected -1) — and use small K (e.g., 20)
  to stay inside test runtime/memory budgets.
publish decision: BLOCK
```

### Novel finding F-2: `notes.md` describes a *different and incorrect* algorithm than what is implemented

```text
area: docs, design
finding: notes.md describes the visited tracking as "(row, col) →
  earliest_arrival_time" with the explicit (and wrong) justification:
  "Once we reach a cell at its earliest possible time, any later arrival
  at that same cell can only lead to equal or worse overall arrival times
  at the destination." This is false for this contract: cells become more
  accessible over time (readyAt unblocks them), so a *later* arrival at
  an intermediate cell may be the only way to reach a far cell. The
  implementation, in contrast, uses `map[[3]int]bool` keyed by
  (row, col, time), which is correct in principle but unbounded in
  practice for large readyAt — the known BLOCK.
evidence:
  - notes.md lines 7-13.
  - solution.go lines 53-54, 84-87: visited keyed by [3]int = row/col/time.
  - TestArrivingLaterUnblocks (solution_test.go lines 137-148) literally
    depends on the BFS revisiting (0,0) at a later time, which the notes'
    described scheme would not allow.
repro/test idea: With the notes' described visited-by-(row,col) scheme,
  TestArrivingLaterUnblocks would return -1. Implement that scheme and
  observe.
impact: This is a critical signal of design confusion. A maintainer
  reading the notes will believe the algorithm is O(rows × cols) BFS and
  may "fix" the memory issue by deleting `time` from the visited key —
  which is *what the notes describe* — and silently break passing tests.
  The known-issue BLOCK and this doc-drift compound each other: any
  attempt to fix the memory bound by following the notes will break
  correctness.
likelihood: high — the notes are an explicit instruction to a future
  maintainer to make the wrong fix.
fix size: small (rewrite notes against the actual algorithm) but the
  underlying algorithm needs replacement anyway, so this is best bundled
  with the algorithmic fix.
publish decision: BLOCK
```

### Novel finding F-3: `notes.md` complexity claim is wrong by a factor of `maxReady`

```text
area: docs
finding: notes.md asserts Time = O(rows × cols) and Space = O(rows × cols).
  With (row, col, time) visited tracking the actual bounds are
  O(rows × cols × maxTime) where maxTime is bounded by rows + cols +
  maxReady. For maxReady = 1e9 the realized complexity is ~1e9 × grid
  size — the source of the known BLOCK.
evidence: notes.md lines 41-49 vs solution.go lines 42, 53, 84.
repro/test idea: Add a microbenchmark on a 2×2 grid with progressively
  larger readyAt; observe runtime and memory growing linearly in readyAt.
impact: Misleads adopters about scaling and obscures the known BLOCK's
  root cause.
likelihood: certain
fix size: small (rewrite complexity section)
publish decision: FIX_BEFORE_TARGET
```

### Novel finding F-4: `maxTime = rows + cols + maxReady` overflows for near-MaxInt readyAt — silent wrong answer

```text
area: numerical safety
finding: When maxReady is close to MaxInt (or rows+cols+maxReady exceeds
  MaxInt), the addition wraps to a negative number. The cutoff check
  `state.time >= maxTime` is then true for any non-negative state.time,
  the BFS skips all expansions, and the function returns -1 for inputs
  whose destination is reachable. Contract: "large readiness values must
  be handled without overflow."
evidence: solution.go line 42 (`maxTime := rows + cols + maxReady`) and
  line 63 (`if state.time >= maxTime`). No saturating-add or int64
  promotion.
repro/test idea: readyAt = [[0,0],[0, math.MaxInt-3]] on a 2×2 grid;
  expected: implementation should return either the very large arrival
  time, or -1 only if structurally unreachable. With the overflow, it
  returns -1 immediately (and instantly, because the BFS terminates after
  the first dequeue).
impact: A silent contract violation for large inputs — exactly the case
  the contract called out.
likelihood: medium (requires near-MaxInt input, but contract explicitly
  asks for "large readiness values" handling).
fix size: small (saturate or use int64 explicitly; cap at MaxInt)
publish decision: FIX_BEFORE_TARGET
```

### Novel finding F-5: Contract ambiguity on re-entering the start cell — implementation chose one interpretation undocumented

```text
area: contract, docs
finding: The contract says "the starting checkpoint is already occupied
  at minute 0". The implementation treats the start specially only for
  the initial enqueue: any later re-entry to (0,0) is gated by
  readyAt[0][0]. If readyAt[0][0] > 0, the worker leaves and cannot
  return to start until time ≥ readyAt[0][0]. This interpretation is
  defensible but is neither stated in the contract nor in notes.md
  Section "Assumptions".
evidence: solution.go line 78-81 unconditionally checks
  `newTime < readyAt[newRow][newCol]`; notes.md "Assumptions" section
  doesn't address re-entry of the origin.
repro/test idea: readyAt = [[5, 0], [0, 0]] on a 2×2 grid. Worker is at
  (0,0) at time 0 (free), moves to (1,0) at t=1, then *cannot* return to
  (0,0) at t=2 (readyAt[0][0]=5). Add a test asserting whatever the
  intended contract behavior is.
impact: Low-severity adoption risk; an early adopter might assume start
  is special forever and write a test that fails.
likelihood: medium
fix size: small (document in notes.md or add a test)
publish decision: DOC_BEFORE_TARGET
```

### Novel finding F-6: Queue uses `queue = queue[1:]` re-slicing — known Go anti-pattern, compounds the BLOCK's memory issue

```text
area: performance, maintainability
finding: solution.go line 60 dequeues with `queue = queue[1:]`. The
  underlying array is never released, so the queue's backing storage
  grows monotonically with total enqueues. Under the BLOCK's
  unbounded-time scenarios, this compounds the visited-map blow-up:
  observed memory is roughly (queue allocations) + (visited entries).
evidence: solution.go line 60.
repro/test idea: pprof a small `readyAt = [[0,0],[0,100]]` run and
  observe the queue's backing array dominating runtime allocation.
impact: Makes the BLOCK fail-mode louder than it needs to be (OOM
  rather than just slow). Once the BLOCK is fixed, this becomes
  POST_TARGET cleanup.
likelihood: high under BLOCK conditions
fix size: small (use an index-based ring queue, or container/list, or
  reset the slice when the head index passes a threshold)
publish decision: POST_TARGET (gated under BLOCK fix)
```

### Material commentary on known issue F (unbounded readyAt)

```text
area: algorithm
finding: I concur with the BLOCK. Sharper framing: the correct algorithm
  is Dijkstra/0-1 BFS over (row, col) with arrival-time as the priority,
  exploiting that the minimum arrival time at any cell c is the smallest
  t such that t ≥ readyAt[c] and t has the same parity as the Manhattan
  distance from the previous frontier cell. The naïve fix of "memoize
  visited by (row, col) earliest_time" — which the notes describe — is
  wrong (see F-2). The correct memo is the earliest *parity-aware*
  arrival time, since each move flips parity. Mention this in the fix
  brief; otherwise the maintainer will rediscover it the slow way.
evidence: notes.md describes an incorrect simpler visited scheme.
repro/test idea: Construct a layout where the minimum-time path requires
  one extra detour to flip parity; verify a parity-naive memo gives a
  wrong answer.
impact: same as known.
likelihood: certain.
fix size: medium.
publish decision: BLOCK (severity unchanged; embedding fix guidance).
```

---

## Important areas checked, no new issue found

- **Case D modular arithmetic**: `int64` is sufficient and modulo is applied after every multiplication; no intermediate overflow path. Both `runSum + dayValue` (≤ 2·MOD) and `minValue · runSum` (≤ MOD²·≈10^18, inside int64) are safe.
- **Case D input mutation**: implementation does not write to `daily`. Contract doesn't forbid mutation explicitly but caller-safe behavior is preserved.
- **Case D empty / single-element handling**: matches contract and tests.
- **Case E input mutation**: explicit copy + sort of both slices; `TestInputNotMutated` exercises this.
- **Case E grantCount > len(staff) edge**: the popcount filter (`grantUsed > grantCount`) correctly handles this; tested.
- **Case E grantSize = 0**: contract-relevant (grant must "help cover" a session); the enumeration naturally returns the same max regardless of mask, so the result is identical to no-grant; tested explicitly.
- **Case F input validation**: empty outer, empty inner, jagged-with-shorter-second-row are all handled by the `cols`/per-row-length check. The asymmetric jagged direction (first row shorter than second) is handled correctly even though not directly tested.
- **Case F 1×1 grid**: short-circuit at line 25 returns 0, including for non-zero `readyAt[0][0]`; explicit test `TestOneCheckpointNonzeroReady`.
- **Case F starting position**: not gated by `readyAt[0][0]` for the initial enqueue (contract: "already occupied"); correct.
- **Cases D/E/F concurrency**: none of these functions touch globals, goroutines, or shared state; no concurrency hazards.

---

## Roll-up

| Case | Known publish decision | After this review |
|------|------------------------|-------------------|
| D    | BLOCK                  | BLOCK (algorithm + tests + notes all wrong-anchored) |
| E    | BLOCK                  | BLOCK (adds silent zero-return at m≥64, weak greedy tests, missing brute-force property test) |
| F    | BLOCK                  | BLOCK (two existing tests assert impossible results; notes describe a different algorithm than implemented; near-MaxInt overflow path) |

For each case the BLOCK fix is larger than the known issue alone:
- **D**: algorithm + test expectations + notes redefinition all need rework as a single change.
- **E**: polynomial rewrite + brute-force-cross-check tests; the `1 << uint(...)` silent-zero bug must be retired with the loop, not patched in place.
- **F**: algorithm rewrite (Dijkstra/0-1 BFS with parity-aware memo, *not* the (row, col)-only memo described in notes) + correction of two miscalibrated tests + overflow-safe `maxTime` computation.
