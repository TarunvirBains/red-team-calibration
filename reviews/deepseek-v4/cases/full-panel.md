Let me run the tests to verify baseline correctness before compiling findings.
All 38 tests pass. Let me verify the critical bugs I identified with targeted edge cases.
Confirmed: `H)2` crashes with index out of range panic. Now compiling findings for all three cases.

---

## Case A: `AuditMeterReadings`

```
area: correctness / overflow
finding: Sum accumulation in the inner loop uses int64 arithmetic but has no overflow detection. If individual readings are near math.MaxInt64 and several are summed, the total wraps silently, producing incorrect counts. The contract requires "use 64-bit arithmetic for totals" which is satisfied, but no defense against wrapping is present.
evidence: solution.go:20 — sum += int64(readings[j]) with no bounds check.
repro/test idea: readings := []int{math.MaxInt64/2, math.MaxInt64/2, math.MaxInt64/2} with maxReadings=3 — the sum wraps to a negative value, and the count may incorrectly include or exclude it.
impact: incorrect audit counts for extreme-valued inputs; could produce misleading compliance data.
likelihood: low (typical meter readings are far below int64 limits).
fix size: small (add overflow check or document input range constraints).
publish decision: POST_TARGET
```

```
area: correctness / maxReadings boundary
finding: The `maxReadings` parameter is typed `int`. The expression `j-i+1 > maxReadings` on line 15 could theoretically overflow if `j-i+1` exceeds int max on a 32-bit platform, allowing periods longer than maxReadings to be counted. On 64-bit (go 1.26 target), this is infeasible but the type mismatch between `maxReadings int` and the returned `int64` count is inconsistent.
evidence: solution.go:15 — j-i+1 > maxReadings; solution.go:4 — func AuditMeterReadings(... maxReadings int) int64
repro/test idea: N/A (infeasible on 64-bit).
impact: none on 64-bit; conceptual inconsistency only.
likelihood: extremely low.
fix size: trivial (change maxReadings to int64 for consistency).
publish decision: NO_ACTION
```

```
area: test coverage
finding: No tests exercise values near int64 boundaries, empty readings with negative maxReadings, or the interaction of maxReadings=1 with large input arrays. The contract only requires tests for "negative values, duplicate totals, zero-length/invalid limits, and a case where ignoring maxReadings gives the wrong answer" — all of which are covered. Coverage is contract-complete but not hardening-complete.
evidence: solution_test.go — all test values are small integers.
repro/test idea: test case with readings containing math.MaxInt32 values.
impact: latent overflow bugs remain undetected.
likelihood: low.
fix size: small (add high-value boundary tests).
publish decision: POST_TARGET
```

---

## Case B: `EntitlementLedger`

```
area: resource / memory leak
finding: The credits map grows without bound. Revoke decrements credits but never deletes map entries for accounts that reach zero credits. Over a long sequence of grant/revoke cycles across varying account ID ranges, the map retains zero-valued entries indefinitely, consuming memory proportional to the union of all account IDs ever touched.
evidence: solution.go:27-28 — Revoke checks > 0 and decrements but never calls delete(); solution.go:8-9 — Grant inserts without bound.
repro/test idea: loop of Grant(0,1000) then Revoke(0,1000) repeated 1000 times — map size stays at 1000 entries even though all credits are zero. Use unsafe.Sizeof or runtime to measure map growth.
impact: unbounded memory growth in long-running service; eventual OOM.
likelihood: medium (depends on deployment lifecycle and account ID diversity).
fix size: small (delete entries when credits reach zero, or add periodic compaction).
publish decision: FIX_BEFORE_TARGET
```

```
area: correctness / overflow
finding: Credit counts are stored as `int`. Repeated grants to the same range (e.g., Grant(0,1) called MaxInt times) cause the credit count to overflow silently, making an account with MaxInt+1 grants appear to have a negative or zero count. Revoke's saturation check would no longer fire correctly because the int value has wrapped past zero.
evidence: solution.go:18 — e.credits[id]++ with no overflow guard.
repro/test idea: create ledger, call Grant(0,1) math.MaxInt32+1 times, then check Eligible(0,1,1) — returns false incorrectly.
impact: incorrect eligibility results under extreme credit accumulation.
likelihood: very low (requires MaxInt operations on same account).
fix size: small (clamp at MaxInt or use a saturated-add).
publish decision: POST_TARGET
```

```
area: correctness / boundary
finding: Grant and Revoke use `left >= right` to detect empty spans, which correctly returns early. However, if `left` or `right` are negative, the for-loop `for id := left; id < right; id++` would iterate through negative IDs. This is handled consistently (the map accepts negative int keys), but the notes.md assumption "Account IDs are non-negative integers" is not enforced by the code. A caller passing negative IDs would get silently accepted.
evidence: solution.go:14-19 — no validation that left/right are non-negative.
repro/test idea: Grant(-5, 0); Eligible(-5, 0, 1) — returns true for negative account IDs.
impact: negative account IDs silently accepted, violating the documented assumption.
likelihood: low (controlled by caller).
fix size: small (assert/panic on negative IDs, or document acceptance in notes).
publish decision: DOC_BEFORE_TARGET
```

---

## Case C: `NetInventory`

```
area: correctness / panic (crash)
finding: The parseFormula stack implementation panics with an index-out-of-range error when encountering a `)` without a matching `(`. The function pops the top map from the stack and then accesses `stack[len(stack)-1]`, but if the stack was only one element (the root), the pop leaves it empty, making the subsequent access panic.
evidence: solution.go:91-93 — stack = stack[:len(stack)-1]; parent := stack[len(stack)-1]; no guard against stack underflow.
repro/test idea: NetInventory("H)2") → panic: runtime error: index out of range [-1]. Confirmed above.
impact: service crash from any malformed input containing unbalanced `)`; denial-of-service vector.
likelihood: low under "all formulas valid" assumption; medium in deployed systems receiving external input.
fix size: small (check stack length before popping and return error or handle gracefully).
publish decision: BLOCK
```

```
area: correctness / data loss
finding: Unclosed `(` causes silent data loss. The parseFormula loop exits without draining the stack, and only `stack[0]` (the root map) is returned. Any atoms parsed inside an unclosed parenthetical group are discarded. This can produce a net-inventory result that is missing atoms the input explicitly declared.
evidence: solution.go:112 — returns stack[0] unconditionally; no check that len(stack) == 1.
repro/test idea: NetInventory("H(Cl") → "H" (should include Cl). Confirmed above.
impact: wrong inventory counts produced from malformed input; potential downstream errors in lab inventory systems that rely on this output.
likelihood: low under "all formulas valid" assumption; medium in deployed systems.
fix size: small (check stack length before returning; return "INVALID" or error for malformed input).
publish decision: BLOCK
```

```
area: correctness / input validation
finding: Trailing `+` or `-` operators are silently accepted. The top-level parser creates an empty string term when a trailing operator is encountered (e.g., "H+"), passes it to parseFormula, which returns an empty map, and this is silently accumulated. The result is the same as the expression without the trailing operator, but the invalid syntax is not rejected.
evidence: solution.go:14-17 — inner loop skips operators; empty term is passed to parseFormula; solution.go:20 — formulaCounts from parseFormula("") is an empty map, silently accepted.
repro/test idea: NetInventory("H+") → "H" (silently accepted). Confirmed above.
impact: silently accepts invalid input, could mask transcription errors.
likelihood: low under "valid formulas" assumption; medium for external input.
fix size: small (reject empty tokens or return error).
publish decision: DOC_BEFORE_TARGET
```

```
area: correctness / ambiguity
finding: parseNumber returns 0 for both "no digits present" and "the digit 0". Callers treat 0 as "no number" and substitute 1. This means a chemical formula with a count of zero (e.g., "H0") is silently interpreted as count 1, producing wrong net inventory. While count-0 is unusual in standard chemical notation, the contract does not restrict counts to positive values.
evidence: solution.go:133-136 — returns 0 for non-digit; solution.go:101-102 — if count == 0 { count = 1 }; solution.go:86-87 — same for multiplier.
repro/test idea: NetInventory("H0") → "H" (should arguably be "" since net count of H is 0). Confirmed above.
impact: subtle incorrect behavior if zero-count atoms appear; limited practical impact.
likelihood: very low (zero-count atoms are unconventional).
fix size: small (use a sentinel or bool return to distinguish "no number" from "number is zero").
publish decision: DOC_BEFORE_TARGET
```

```
area: performance / overflow
finding: parseNumber uses int for accumulation and can overflow on extremely long digit strings (20+ digits). The expression `num = num*10 + int(s[*i]-'0')` wraps silently for values exceeding int max, producing an incorrect multiplier or count. This can produce both wrong atom counts and wrong inventory results.
evidence: solution.go:137-141 — num*10 with no overflow check.
repro/test idea: NetInventory("H99999999999999999999") — count wraps to an incorrect value.
impact: incorrect inventory results from inputs with extreme digit counts.
likelihood: very low (requires unrealistic input sizes).
fix size: small (cap at MaxInt, or validate digit string length).
publish decision: POST_TARGET
```

---

## Summary

| Decision | Count | Cases |
|---|---|---|
| **BLOCK** | 2 | Case C: unbalanced `)` panic, unclosed `(` data loss |
| **FIX_BEFORE_TARGET** | 1 | Case B: unbounded map growth |
| **DOC_BEFORE_TARGET** | 3 | Case B: negative IDs not validated; Case C: trailing operator acceptance; Case C: parseNumber 0 ambiguity |
| **POST_TARGET** | 4 | Case A: overflow, test gaps; Case B: credit overflow; Case C: parseNumber overflow |
| **NO_ACTION** | 1 | Case A: int/int64 type inconsistency |
