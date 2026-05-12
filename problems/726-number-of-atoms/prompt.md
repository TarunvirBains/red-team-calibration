# 726 Number of Atoms

Task summary: Parse a valid chemical formula with atoms, optional counts, and
nested parenthesized groups. Return the canonical count string with atom names
sorted lexicographically and omitted count `1`.

Review pressure points:

- uppercase/lowercase atom tokenization;
- multi-digit counts;
- nested group multiplication;
- stack merge order;
- canonical output ordering;
- malformed input behavior if reused outside the LeetCode validity contract.

