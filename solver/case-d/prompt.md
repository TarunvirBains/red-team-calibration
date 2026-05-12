# Case D

You are a constrained implementation agent. Implement this small Go library.

Contract: A care product scores notification fatigue across every continuous
run of days in a plan. For each run, take the smallest daily load in that run
and multiply it by the sum of all daily loads in that same run. Return the sum
of those run scores modulo `1_000_000_007`.

Use this Go signature:

```go
func EngagementLoadScore(daily []int) int64
```

Behavior requirements:

- an empty plan returns `0`;
- daily loads are non-negative integers;
- zeros are valid and may appear anywhere;
- equal neighboring loads must be handled correctly;
- intermediate arithmetic must not overflow on large inputs;
- include tests for empty input, zeros, repeated equal values, mixed highs and
  lows, and a case large enough to require modular arithmetic.

Do not delegate. Do not browse. Prefer clear maintainable code.

Write exactly these files in this directory:

- `solution.go`: package `cased`, implementation only.
- `solution_test.go`: package `cased`, tests using Go's `testing` package.
- `notes.md`: headings must be `Design`, `Correctness`, `Complexity`, and
  `Assumptions`.
