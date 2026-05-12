# Case B: Entitlement Ledger

## Design

The implementation uses a simple map-based approach storing credits per account ID.

- `EntitlementLedger` maintains a `map[int]int` where keys are account IDs and values are credit counts
- `Grant(left, right)` iterates through the range [left, right) and increments each account's credit
- `Revoke(left, right)` iterates through the range [left, right) and decrements each account's credit, but never below zero
- `Eligible(left, right, minCredits)` iterates through the range [left, right) and checks if all accounts have at least minCredits credits

Empty ranges (left >= right) are handled as no-ops for mutations and return true for eligibility checks.

## Correctness

The implementation satisfies all contract requirements:

- Half-open interval semantics: all ranges use `left <= id < right`
- Empty span handling: `left >= right` returns early or true
- Threshold behavior: `minCredits <= 0` explicitly returns true
- Saturation: revokes check `credits[id] > 0` before decrementing
- Credit preservation: map stores exact counts, enabling overlapping grants and partial revokes
- Tests cover overlapping grants, partial revokes, adjacent spans, multi-credit thresholds, empty spans, and complex operation sequences

## Complexity

- Space: O(n) where n is the number of accounts that have been granted credits
- Grant: O(right - left)
- Revoke: O(right - left)
- Eligible: O(right - left)

Linear per-operation cost is appropriate for this use case and provides clear maintainable code.

## Assumptions

- Account IDs are non-negative integers
- The credit system is non-negative (zero or positive)
- Operations are sequential (no concurrent access)
- Empty spans are valid inputs and should be handled gracefully
