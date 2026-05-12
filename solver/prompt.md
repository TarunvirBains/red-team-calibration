# Solver Prompt

You are a constrained implementation agent. Solve the three algorithm tasks
below in one response.

Use Rust for all three solutions. For each problem include:

- algorithm overview;
- correctness argument;
- complexity;
- Rust implementation;
- at least five targeted tests, including edge cases.

Do not delegate. Do not browse. Do not optimize for cleverness over clarity.
Prefer LeetCode-compatible Rust signatures where practical, and call out any
helper structs or assumptions needed for local testing.

## Problems

1. LeetCode 327, Count of Range Sum.
   Given an integer array and inclusive bounds `lower` and `upper`, count the
   number of subarray sums that fall within the inclusive range.

2. LeetCode 715, Range Module.
   Implement a mutable interval tracker supporting add, query, and remove
   operations over half-open ranges.

3. LeetCode 726, Number of Atoms.
   Parse a valid chemical formula and return the canonical atom count string
   sorted lexicographically by atom name.

## Output Shape

Use these headings exactly:

```text
# 327 Count of Range Sum
# 715 Range Module
# 726 Number of Atoms
```
