# Case B

You are a constrained implementation agent. Implement this small Go library.

Contract: A customer platform manages eligibility credits for account id spans.
`Grant(left, right)` adds one credit to every account id in `[left, right)`.
`Revoke(left, right)` removes one credit but never below `0`.
`Eligible(left, right, minCredits)` returns true iff every account id in
`[left, right)` has at least `minCredits` credits.

Use this Go API:

```go
type EntitlementLedger struct { /* ... */ }
func NewEntitlementLedger() EntitlementLedger
func (e *EntitlementLedger) Grant(left int, right int)
func (e *EntitlementLedger) Revoke(left int, right int)
func (e *EntitlementLedger) Eligible(left int, right int, minCredits int) bool
```

Behavior requirements:

- spans are half-open;
- empty spans are no-ops, and checking an empty span returns true;
- `minCredits <= 0` returns true for any checked span;
- revokes saturate at zero;
- credit counts matter, not just present/absent state;
- include tests for overlapping grants, partial revokes, adjacent spans,
  thresholds above `1`, empty spans, and long operation sequences.

Do not delegate. Do not browse. Prefer clear maintainable code.

Write exactly these files in this directory:

- `solution.go`: package `caseb`, implementation only.
- `solution_test.go`: package `caseb`, tests using Go's `testing` package.
- `notes.md`: headings must be `Design`, `Correctness`, `Complexity`, and
  `Assumptions`.
