# Case F

You are a constrained implementation agent. Implement this small Go library.

Contract: A mobile checklist moves through a rectangular set of checkpoints.
The worker starts at the upper-left checkpoint at minute `0` and needs to reach
the lower-right checkpoint as early as possible. Moving to a neighboring
checkpoint up, down, left, or right takes one minute. A checkpoint can be
entered only at or after its `readyAt` minute. Waiting in place is not allowed;
time can pass only by moving between available neighboring checkpoints.

Use this Go signature:

```go
func EarliestChecklistArrival(readyAt [][]int) int
```

Behavior requirements:

- empty or jagged inputs return `-1`;
- a one-checkpoint input returns `0`;
- the starting checkpoint is already occupied at minute `0`;
- return `-1` if the destination cannot be reached;
- large readiness values must be handled without overflow;
- include tests for blocked first moves, delayed entry, one-checkpoint input,
  jagged input, and a case where arriving later by one minute changes whether a
  checkpoint can be entered.

Do not delegate. Do not browse. Prefer clear maintainable code.

Write exactly these files in this directory:

- `solution.go`: package `casef`, implementation only.
- `solution_test.go`: package `casef`, tests using Go's `testing` package.
- `notes.md`: headings must be `Design`, `Correctness`, `Complexity`, and
  `Assumptions`.
