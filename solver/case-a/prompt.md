# Case A

You are a constrained implementation agent. Implement this small Go library.

Contract: A utility company audits reviewable periods in meter data. A
reviewable period starts and ends on recorded readings, contains no more than
`maxReadings` readings, and has a signed usage total inside the inclusive
manual-review band `[minTotal, maxTotal]`. Return the number of reviewable
periods in `readings`.

Use this Go signature:

```go
func AuditMeterReadings(readings []int, minTotal int64, maxTotal int64, maxReadings int) int64
```

Behavior requirements:

- bounds are inclusive;
- periods longer than `maxReadings` do not count;
- if `maxReadings <= 0`, return `0`;
- use 64-bit arithmetic for totals;
- include tests for negative values, duplicate totals, zero-length/invalid
  limits, and a case where ignoring `maxReadings` gives the wrong answer.

Do not delegate. Do not browse. Prefer clear maintainable code.

Write exactly these files in this directory:

- `solution.go`: package `casea`, implementation only.
- `solution_test.go`: package `casea`, tests using Go's `testing` package.
- `notes.md`: headings must be `Design`, `Correctness`, `Complexity`, and
  `Assumptions`.
