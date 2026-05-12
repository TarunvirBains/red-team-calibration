# 327 Count of Range Sum

Task summary: Given an integer array and inclusive integer bounds `lower` and
`upper`, count the number of contiguous subarrays whose sums lie in that
inclusive range.

Review pressure points:

- inclusive lower/upper bounds;
- empty prefix and zero-length mistakes;
- negative numbers and duplicates;
- 32-bit overflow risks in prefix sums;
- O(n log n) expectation for large inputs;
- merge-sort, Fenwick, or balanced-tree off-by-one errors.

