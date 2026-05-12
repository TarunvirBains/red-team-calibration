# Case C

You are a constrained implementation agent. Implement this small Go library.

Contract: A lab inventory import records delivered and withdrawn compounds in a
single expression. Parse an expression made of valid chemical formulas joined
by `+` and `-`. Deliveries add formula counts and withdrawals subtract formula
counts. Parentheses and multipliers work inside each formula. Return the
canonical positive net atom count string sorted lexicographically, omitting
atoms whose net count is zero. If any atom would have a negative final count,
return `"INVALID"`.

Examples:

- `H2O+NaCl-Na` -> `ClH2O`
- `Mg(OH)2-H2` -> `MgO2`
- `H2-O3` -> `INVALID`

Use this Go signature:

```go
func NetInventory(expr string) string
```

Behavior requirements:

- formulas inside each inventory term are valid;
- separators are only top-level `+` and `-`;
- a missing leading sign means addition;
- multi-digit counts are allowed;
- atom names use an uppercase letter followed by zero or more lowercase letters;
- omit atoms with net count zero;
- include tests for nested groups, subtraction, zero net counts, invalid
  negative final counts, and lexicographic output ordering.

Do not delegate. Do not browse. Prefer clear maintainable code.

Write exactly these files in this directory:

- `solution.go`: package `casec`, implementation only.
- `solution_test.go`: package `casec`, tests using Go's `testing` package.
- `notes.md`: headings must be `Design`, `Correctness`, `Complexity`, and
  `Assumptions`.
