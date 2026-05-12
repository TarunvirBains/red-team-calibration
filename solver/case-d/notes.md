# EngagementLoadScore Implementation Notes

## Design

A "continuous run" is a maximal sequence of consecutive days with identical daily loads. The implementation:

1. Scan left-to-right through the daily loads
2. For each new run of equal consecutive values:
   - The minimum load in the run is the value itself (all are equal)
   - Accumulate the sum of all loads in the run
   - Compute run score as `(min × sum) % MOD`
3. Add each run score to a total, taking modulo after each addition
4. Return the final total modulo `1,000,000,007`

Example: `[5, 5, 3, 3, 3, 2]`
- Run 1: `[5, 5]` → min=5, sum=10, score=(5×10)%MOD=50
- Run 2: `[3, 3, 3]` → min=3, sum=9, score=(3×9)%MOD=27
- Run 3: `[2]` → min=2, sum=2, score=(2×2)%MOD=4
- **Total**: (50+27+4)%MOD = 81

## Correctness

- **Empty input**: Returns 0 as specified
- **Single elements**: Correctly handled as runs of length 1; score is value²
- **Zeros**: Contribute 0 to the final score (0×0=0), correctly placed anywhere in the sequence
- **Equal neighbors**: Consecutive identical loads form a single run; transitions to different values create separate runs
- **Overflow prevention**: All arithmetic uses `int64` and applies modulo `1,000,000,007` during accumulation to prevent overflow before multiplication
- **Modular arithmetic**: Intermediate sums and products are taken modulo before further operations, preserving mathematical correctness while preventing overflow

Test cases cover:
- Empty input, single elements, all zeros
- Runs of equal values (length 1, length 2, length 3+)
- Alternating values (each creates its own run)
- Zeros intermixed with non-zero values
- Large numbers (up to 10⁹) and large counts (100+ consecutive values) requiring modular reduction

## Complexity

- **Time**: O(n) — single left-to-right pass through the input array; each element visited once
- **Space**: O(1) — only a handful of loop variables and accumulators; no auxiliary data structures

## Assumptions

- All input values are non-negative integers
- Integer overflow is prevented by applying modulo `1,000,000,007` during accumulation and multiplication, not just at the final result
- The modulo base `1,000,000,007` is used consistently throughout
- Input arrays can be arbitrarily large without functional constraints beyond system memory
- "Notification fatigue scoring" semantics: runs of identical loads create fatigue proportional to both the load magnitude and the total burden over the run period
