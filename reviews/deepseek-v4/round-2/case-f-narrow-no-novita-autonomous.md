Now let me verify one more thing — the exact counterexample proving the notes are wrong:
Now I have all the evidence needed. Here is the final review:

---

# Case F Narrow Review — Final Report

## Novel Findings

```text
area: test assertion
finding: TestLargeReadyAtReachable asserts unreachable (-1) for {{0,0},{0,10}}, but the destination IS reachable at time 10 via bouncing (0,0)↔(0,1)↔(0,0) until arrival. Correct answer: 10. This is the concrete test instance of the BLOCK known-issue about 2×2 reachability via backtracking.
evidence: go test -run TestLargeReadyAtReachable returns got 10 (not -1). Traced: (0,0)@0→(0,1)@1→(0,0)@2→(0,1)@3→(0,0)@4→(0,1)@5→(0,0)@6→(0,1)@7→(0,0)@8→(1,0)@9→(1,1)@10, readyAt=10, enters at exactly 10.
repro/test idea: Run the existing test; it fails with "got 10".
impact: Test has wrong expected value. A developer fixing the algorithm might think this test validates reachable cases and not realize the expected value is stale.
likelihood: certain (test fails every run)
fix size: small (change expected -1 to 10 in solution_test.go:120)
publish decision: BLOCK
```

```text
area: test performance / resource exhaustion
finding: TestLargeReadyAtUnreachable uses maxReady=1000000000, causing maxTime=1000000004. Since the destination IS reachable (bouncing reaches t=1000000000), the algorithm explores ~4B states before returning. The test timed out at 30s without completing. The test has two defects: (a) it expects -1 but the correct answer is 1000000000 (reachable), and (b) it demonstrates that the algorithm is catastrophically slow for even trivially small grids with large readyAt values — well below the "unbounded time and memory" threshold from the known BLOCK.
evidence: rtk go test -run TestLargeReadyAtUnreachable -timeout 30s → FAIL (timeout after 30.311s). Grid {{0,0},{0,1000000000}} has maxTime = 4 + 1e9, visited keys = 4 × 1e9 ≈ 4B map entries, each requiring multiple allocations. Real-world run time is minutes-to-hours for a 2×2 grid.
repro/test idea: Run the test with any timeout under several minutes.
impact: A test that uses 1B as a readyAt value makes the test suite unusable in CI. The choice of 1B is arbitrary and does not demonstrate correctness — it only demonstrates the performance bug.
likelihood: certain (test hangs every run)
fix size: small (change readyAt value to something reasonable, or remove until algorithm is fixed)
publish decision: BLOCK
```

```text
area: notes correctness (material disagreement with known issue #2)
finding: The notes do not just "misstate" visited-state shape. They describe a fundamentally incorrect algorithm — tracking visited by (row,col) only — and argue its correctness with a reasoning that is demonstrably false. The argument "any later arrival at a cell can only lead to equal or worse overall arrival times" is refuted by the grid {{0,3},{0,5}} (and similar counterexamples). In that grid, visiting (0,0) only at t=0 (and skipping t=2,4) prevents reaching (1,1) at t=6 — the correct answer. The code's (row,col,time) tracking is correct but expensive; the notes' (row,col)-only approach is simple but incorrect. The notes therefore provide wrong implementation guidance on top of misstating the resource cost.
evidence: Harness test on {{0,3},{0,5}} returns 6 by revisiting (0,0) at t=2 and t=4. With (row,col)-only visited, the algorithm dead-ends at t=1 since (0,0)@t=2 is skipped as "already visited." The path requiring revisitation is (0,0)@0→(1,0)@1→(0,0)@2→(0,1)@3→(0,0)@4→(1,0)@5→(1,1)@6.
repro/test idea: Implement visited as map[[2]int]bool and run on {{0,3},{0,5}}. Observe -1 instead of 6.
impact: Anyone implementing from the notes (rather than the code) would produce an incorrect solution that silently returns false negatives.
likelihood: certain (the argument in notes:52-53 is provably wrong)
fix size: medium (rewrite the Design and Correctness sections to match the code's actual algorithm)
publish decision: BLOCK
```

```text
area: overflow severity (sharper evidence for known issue #4)
finding: The overflow from `rows + cols + maxReady` is more severe than "can overflow." When maxReady ≥ math.MaxInt - rows - cols, maxTime wraps to a negative value. Since `state.time >= negative` is always true (time ≥ 0), the `if state.time >= maxTime { continue }` guard skips ALL states after (0,0,0). The queue drains immediately and the algorithm always returns -1 — even for grids where the destination is trivially reachable without bouncing. This is a silent correctness failure, not just a resource-bound miscalculation.
evidence: rows=cols=2, maxReady=math.MaxInt → maxTime = -9223372036854775805. state.time=0 >= -9223372036854775805 → true → skip. Queue processes only (0,0,0) then returns -1.
repro/test idea: Call EarliestChecklistArrival([][]int{{0,0},{0,math.MaxInt}}). Returns -1 despite the destination being reachable.
impact: Silent denial of service to correct answers for any input where maxReady is near MaxInt. The failure mode is invisible — no panic, no error, just a wrong -1.
likelihood: certain for any input with maxReady in the overflow range
fix size: small (guard maxTime against overflow, or cap it)
publish decision: BLOCK
```

```text
area: test naming
finding: TestDelayedDestinationUnreachable is misnamed — it expects result 6 (reachable, not unreachable). The destination IS reachable via bouncing. The name implies the test covers an unreachable case, which is misleading.
evidence: The test expects 6 ≠ -1. The test's own comment describes a path reaching (2,2) at time 6.
repro/test idea: Read test name vs expected value.
impact: Low — name confusion only. But contributes to the pattern of tests being wrong or misleading about reachability.
likelihood: certain
fix size: small (rename to TestDelayedDestinationReachableAfterBounce or similar)
publish decision: POST_TARGET
```

## Known-Issue Re-examination

All six known issues are confirmed with no material downgrades:

| # | Known Issue | Verdict |
|---|------------|---------|
| 1 | Unbounded time/memory from large readyAt | **Confirmed.** `TestLargeReadyAtUnreachable` with 1B hangs the suite. |
| 2 | Notes misstate visited shape and resource cost | **Confirmed and escalated** — notes describe a provably broken algorithm, not just different. |
| 3 | Tests assert unreachable for reachable 2×2 | **Confirmed.** `TestLargeReadyAtReachable` expects -1, actual=10. |
| 4 | Overflow in rows+cols+maxReady | **Confirmed and escalated** — overflow causes silent false -1 for all inputs. |
| 5 | Queue reslicing retains backing storage | **Confirmed** but not independently tested (code visible). |
| 6 | Missing single-row/col delayed-neighbor tests | **Confirmed** — still absent. |

## Areas Checked Without New Findings

- **Correctness for valid, moderate-size inputs**: All non-large tests pass correctly (15/15 tests under 10s). The algorithm is correct for bounded readyAt values.
- **Input validation (empty, jagged, 1-cell)**: All edge cases handled correctly and tested.
- **Parity constraints on arrival times**: No correctness issue — the algorithm correctly finds earliest arrival respecting parity through exhaustive (row, col, time) search.
- **Bouncing/large-readyAt unreachable vs reachable distinction**: `TestDestinationBlocked` (expects 6) is correct. `TestComplexPath` (expects 5) is correct. Only the tests involving very large readyAt values are wrong.
- **No infinite loops**: The maxTime bound prevents unbounded growth in theory, though it's insufficient as a practical stop-gap.
