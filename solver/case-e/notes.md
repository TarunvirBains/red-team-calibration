# MaxCoveredSessions Implementation

## Design

The implementation uses exhaustive enumeration of grant assignments combined with greedy matching:

**Grant Assignment Enumeration:**
For each possible way to distribute grants to staff members (2^m possibilities where m = number of staff):
1. Create an augmented capacity array where selected staff members receive the grant
2. Perform a greedy matching of sessions to staff

**Greedy Matching:**
For each session (in sorted order from smallest to largest), match it to the unused staff member with the smallest capacity that can still cover the session. This "best fit" strategy minimizes waste and preserves capacity for remaining sessions.

**Optimal Result:**
Return the maximum number of sessions covered across all grant assignment scenarios.

The implementation avoids mutating input slices by creating independent sorted copies before processing.

## Correctness

**Enumeration guarantees optimality:**
By trying every possible grant assignment (0 to min(grantCount, len(staff)) grants distributed across all staff), the implementation considers the optimal grant distribution. For each assignment, the greedy matching is optimal: matching each session to the smallest-capacity unused staff that can cover it minimizes waste and maximizes future flexibility.

**Grant constraints enforced:**
- Each staff member receives at most one grant (mask ensures binary choice per staff)
- Each staff member covers at most one session (usedLocal tracking)
- Grant count is respected (grantUsed ≤ grantCount check)

**Edge cases:**
- Empty session or staff slices return 0 immediately
- Non-negative efforts and capacities are handled naturally by comparison
- Zero grants or zero grant size are handled correctly (skips invalid masks)
- Grants larger than needed are acceptable

**No mutation guarantee:**
Input slices are copied before sorting, leaving caller-owned data untouched.

## Complexity

**Time Complexity:** O(2^m · (m + n log n)) where n = sessions, m = staff
- Sorting: O(n log n) for sessions + O(m log m) for staff
- Enumerate grants: 2^m iterations
  - Each iteration:
    - Augment capacities: O(m)
    - Greedy matching: O(n · m) to find best staff for each session
- Total: 2^m · (m + n·m)

For practical constraints (m ≤ 20), this is acceptable. With m=20, 2^m ≈ 1M iterations, each doing O(n·m) = O(n·20) work.

**Space Complexity:** O(n + m)
- Sorted copies of sessions and staff: O(n + m)
- Boolean arrays and capacity array: O(m)

## Assumptions

- Session efforts and staff capacities are non-negative integers
- Grant count and grant size are non-negative integers
- Staff count m is small enough (m ≤ 20-25) for 2^m enumeration to be practical
- Each staff member can receive at most one grant (binary choice per staff)
- Each staff member can cover at most one session (exclusive matching)
- A session is either fully covered by one staff member or not covered at all (no partial coverage)
- Capacity comparison is straightforward: staff capacity ≥ session effort means can cover
- The goal is to maximize the absolute count of covered sessions
- All grants have the same size (cannot customize per staff member)
