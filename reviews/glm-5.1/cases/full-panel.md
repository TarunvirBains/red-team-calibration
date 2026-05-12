Let me read all the target files first.
Let me run the tests and verify the code compiles:
All tests pass and `go vet` is clean. Now performing direct analysis.

---

# Red-Team Review: PR_GATE Calibration

## CASE A

### Finding A-1

```
area:           arithmetic overflow
finding:        Sum accumulation in the inner loop can overflow int64. With large maxReadings and readings near int32-max, the running sum on line 20 (sum += int64(readings[j])) exceeds 2^63-1, wrapping to a negative value and producing silently wrong counts or missing/phantom reviewable periods.
evidence:       solution.go:20. Each reading is cast to int64 and added without overflow check. If maxReadings=2e9 and each reading=2.1e9, the sum reaches ~4.2e18, within int64 range. But maxReadings=5e9 and readings=2.1e9 gives ~1.05e19, exceeding int64 max (9.2e18). The notes claim "fit in int32" (line 35 of notes.md) but the signature accepts []int, not int32, and the constraint is neither enforced nor tested.
repro/test idea: Test with readings of value 2e9, maxReadings=5000, minTotal=0, maxTotal=math.MaxInt64. Verify the returned count matches a BigInt reference implementation.
impact:         Silent data corruption — wrong audit count, possible false negatives in meter review
likelihood:     medium — feasible with large meter datasets
fix size:       small — add overflow check on sum += or bound maxReadings/reading values
publish decision: FIX_BEFORE_TARGET
```

### Finding A-2

```
area:           missing validation
finding:        When minTotal > maxTotal the function silently returns 0 instead of signaling invalid input. This is documented as an assumption (notes.md:37) but not tested, and a caller passing swapped bounds gets no diagnostic.
evidence:       solution.go:23. Condition `sum >= minTotal && sum <= maxTotal` is never true when minTotal > maxTotal, so the loop never counts anything. No test exercises this case.
repro/test idea: Test AuditMeterReadings([]int{1,2,3}, 10, 5, 3) — expect either 0 with documented guarantee or an error return.
impact:         Masked caller bugs — swapped arguments silently produce a plausible-looking zero result
likelihood:     medium — parameter ordering mistakes are common
fix size:       small — either validate and return error/panic, or add a test documenting the behavior explicitly
publish decision: DOC_BEFORE_TARGET
```

### Finding A-3

```
area:           test gap — negative-readings mixed with large positive
finding:        The test suite covers negative readings and zero readings but never exercises a period where the running sum crosses zero (e.g., alternating large positive and negative values). This is the regime where int64 overflow is most likely to produce a wrong-sign intermediate that passes the bounds check spuriously.
evidence:       solution_test.go — all test cases use small-magnitude values. No case has sum that wraps past int64 min/max.
repro/test idea: readings=[]int{math.MaxInt32, math.MaxInt32, -1}, maxReadings=3, minTotal=0, maxTotal=math.MaxInt64 — check that the sum of the first two elements (which overflows) is not counted.
impact:         Untested overflow regime is exactly where overflow manifests
likelihood:     low for typical data, high for adversarial input
fix size:       small — add one targeted test case
publish decision: FIX_BEFORE_TARGET
```

---

## CASE B

### Finding B-1

```
area:           value-semantics aliasing bug
finding:        NewEntitlementLedger returns EntitlementLedger by value, but all methods use pointer receivers. If a caller copies the struct (e.g., e2 := e1), both variables share the same underlying map because maps are reference types in Go. Mutations via e2 are visible through e1, violating value-type expectations and causing data corruption.
evidence:       solution.go:7-11. Returns value type. solution.go:13,22,33 — all pointer receivers. No copy protection (no unexported field that would prevent copying, no Copy method, no documentation warning).
repro/test idea: e1 := NewEntitlementLedger(); e1.Grant(0,5); e2 := e1; e2.Grant(5,10); assert e1.Eligible(5,10,1) == true (unexpectedly true because e1 and e2 share the map).
impact:         Silent data corruption in any code path that copies the ledger — credits granted/revoked on one copy bleed into another
likelihood:     high — value-return + pointer-receiver is a common Go pattern mistake; copying structs is idiomatic Go
fix size:       small — either (a) return *EntitlementLedger from constructor and make the type unexported, or (b) add an unexported noCopy field and document, or (c) embed sync.Locker and vet against copies
publish decision: BLOCK
```

### Finding B-2

```
area:           missing negative-id validation
finding:        The notes assume "non-negative integers" (notes.md:36) but the code accepts any int, including negative values. Grant(-5, 0) silently creates map keys -5 through -1. This is valid Go but contradicts the documented domain and could produce surprising Eligible results if callers assume only non-negative IDs exist.
evidence:       solution.go:13-19. Loop `for id := left; id < right` works with negative values. No guard.
repro/test idea: e := NewEntitlementLedger(); e.Grant(-3, 2); assert e.Eligible(-3, 0, 1) == true — passes but is outside stated domain.
impact:         Accepts out-of-domain input silently, potential for logic errors if callers rely on non-negative invariant
likelihood:     low — requires caller to pass negative IDs
fix size:       small — add guard in Grant/Revoke/Eligible or document explicitly
publish decision: DOC_BEFORE_TARGET
```

### Finding B-3

```
area:           concurrency safety
finding:        The EntitlementLedger has no synchronization. The notes acknowledge "no concurrent access" (notes.md:38) but the type is exported with public methods and no documentation of this restriction. A real "customer platform" would likely have concurrent callers, and map access without synchronization causes data races and panics in Go.
evidence:       solution.go — no mutex, no comment on concurrency. notes.md:38 is the only mention.
repro/test idea: Run Grant and Eligible concurrently from two goroutines with -race flag; expect race detector to fire.
impact:         Data races and panics under concurrent use
likelihood:     medium — depends on deployment pattern but likely for a "platform" API
fix size:       medium — add sync.Mutex or RWMutex; or document the restriction on the exported type
publish decision: FIX_BEFORE_TARGET
```

---

## CASE C

### Finding C-1

```
area:           crash on malformed input (unmatched closing paren)
finding:        An unmatched ) causes the stack to be popped below its base element, then the parent access on line 93 (stack[len(stack)-1]) indexes out of range, producing a runtime panic. This is a denial-of-service vector if input is not pre-validated.
evidence:       solution.go:82-96. When ) is encountered, the code pops stack[len(stack)-1] and then reads stack[len(stack)-1] as parent. If the stack has only the base element, the pop removes it, and the next access is stack[-1] → index out of range panic.
repro/test idea: NetInventory("H2)") — expect panic at runtime.
impact:         Unhandled panic crashes the process; DoS if input is externally supplied
likelihood:     medium — the prompt says formulas are valid, but the code exports a public API with no input guard
fix size:       small — check stack depth before popping; return error or "INVALID" on malformed input
publish decision: BLOCK
```

### Finding C-2

```
area:           custom int-to-string reimplements strconv.Itoa
finding:        Lines 58-63 implement manual digit-by-digit integer-to-string conversion instead of using the standard library strconv.Itoa. This is unidiomatic, harder to audit, and fragile — while it works for positive integers in the tested range, it has no edge-case handling and duplicates well-tested stdlib code.
evidence:       solution.go:58-63. Manual loop with modulo arithmetic. No import of "strconv".
repro/test idea: Fuzz: call the manual itoa logic with count=1 through 1000000 and compare against strconv.Itoa.
impact:         Low direct impact (appears correct for the domain), but maintenance burden and audit risk
likelihood:     low — works for current inputs
fix size:       trivial — replace with strconv.Itoa(count)
publish decision: FIX_BEFORE_TARGET
```

### Finding C-3

```
area:           undocumented H0→H1 normalization
finding:        When a count of 0 is parsed (e.g., H0, (OH)0), parseNumber returns 0, and the caller replaces it with 1 (lines 86-88, 101-103). This means H0 is treated identically to H, and (OH)0 is treated as (OH)1. The behavior is undocumented and violates chemical intuition (zero copies should mean absence).
evidence:       solution.go:86-88 (group multiplier 0→1), solution.go:101-103 (atom count 0→1). No test exercises explicit zero multipliers.
repro/test idea: NetInventory("H0") — expect "H" (current behavior) vs "" (chemical intuition). Assert and document which is correct.
impact:         Surprising semantics for a chemistry domain — zero copies meaning "one copy" is wrong
likelihood:     low — H0 is not standard notation, but (X)0 could appear in subtraction expressions
fix size:       small — either return error/INVALID for zero multiplier, or omit atoms with count 0, and add a test
publish decision: FIX_BEFORE_TARGET
```

### Finding C-4

```
area:           test gap — empty expression and single-operator expressions
finding:        NetInventory("") returns "" and NetInventory("+H") or "-H" are untested. The top-level parser assumes the first term has no sign, so "-H" would parse sign=-1, then term="" (empty string), then on the next iteration see H with sign=-1. The empty term goes through parseFormula("") which returns an empty map — harmless. But "-H" starts with sign=-1 before any term is read, so the first character is consumed as sign, leaving "H" with sign=-1, giving H:-1 → "INVALID". This behavior is undocumented and untested.
evidence:       solution.go:11-33. The sign variable starts at 1 (addition). The loop reads a term (chars until + or -), then reads the operator. For "-H", the first term is "" and the operator is consumed as the leading -, then "H" is read with sign=-1.
repro/test idea: Add test cases for "", "-H", "+H", "H+", and verify/document behavior.
impact:         Undefined behavior at API boundary for edge-case inputs
likelihood:     low — requires leading operator, which violates "missing leading sign means addition"
fix size:       small — add test cases and either reject or document
publish decision: DOC_BEFORE_TARGET
```

### Finding C-5

```
area:           integer overflow in parseNumber and group multiplication
finding:        parseNumber (line 138) accumulates digits into an int with no overflow check. Similarly, line 95 multiplies count * multiplier without overflow check. With large multipliers (e.g., (H)999999999999999999), the count overflows int, producing wrong atom counts silently.
evidence:       solution.go:138-142 (num*10 + digit overflow), solution.go:95 (count*multiplier overflow). No bounds on input.
repro/test idea: NetInventory("(H)999999999999999") — compare result against BigInt reference.
impact:         Silent overflow produces wrong atom counts or negative intermediate values that may or may not trigger the INVALID check
likelihood:     low — requires adversarial or buggy input
fix size:       medium — add overflow checks or cap multiplier size
publish decision: FIX_BEFORE_TARGET
```

---

## Summary

| ID  | Area                    | Decision            |
|-----|-------------------------|---------------------|
| A-1 | int64 overflow          | FIX_BEFORE_TARGET   |
| A-2 | missing validation      | DOC_BEFORE_TARGET   |
| A-3 | test gap (overflow)      | FIX_BEFORE_TARGET   |
| B-1 | value-semantics aliasing | **BLOCK**           |
| B-2 | negative ID accepted    | DOC_BEFORE_TARGET   |
| B-3 | no concurrency safety   | FIX_BEFORE_TARGET   |
| C-1 | panic on malformed input | **BLOCK**           |
| C-2 | custom itoa vs stdlib   | FIX_BEFORE_TARGET   |
| C-3 | H0→H1 normalization     | FIX_BEFORE_TARGET   |
| C-4 | edge-case input untested| DOC_BEFORE_TARGET   |
| C-5 | int overflow in parser  | FIX_BEFORE_TARGET   |

**2 BLOCKs** (B-1 aliasing bug, C-1 panic), 6 FIX_BEFORE_TARGET, 3 DOC_BEFORE_TARGET.
