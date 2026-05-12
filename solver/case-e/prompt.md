# Case E

You are a constrained implementation agent. Implement this small Go library.

Contract: A coaching operations tool matches requested care sessions to
available staff capacity. Each session needs one staff member whose capacity is
at least that session's effort. Each staff member can cover at most one
session. A limited number of one-time support grants may be applied, and each
grant increases one staff member's capacity by `grantSize`. A staff member can
receive at most one grant. Return the largest number of sessions that can be
covered.

Use this Go signature:

```go
func MaxCoveredSessions(sessionEffort []int, staffCapacity []int, grantCount int, grantSize int) int
```

Behavior requirements:

- session and staff slices may be empty;
- a grant may be used only when it helps cover one selected session;
- do not mutate caller-owned slices;
- duplicate efforts and duplicate capacities are common;
- return `0` when no session can be covered;
- include tests for exact fits, duplicate values, grants that are necessary,
  grants that are insufficient, more grants than useful placements, and inputs
  where a locally attractive pairing leaves fewer sessions covered.

Do not delegate. Do not browse. Prefer clear maintainable code.

Write exactly these files in this directory:

- `solution.go`: package `casee`, implementation only.
- `solution_test.go`: package `casee`, tests using Go's `testing` package.
- `notes.md`: headings must be `Design`, `Correctness`, `Complexity`, and
  `Assumptions`.
