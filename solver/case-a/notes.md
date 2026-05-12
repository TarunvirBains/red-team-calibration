# Notes: Case A

## Design

The implementation checks each possible review period in the meter readings. For each starting position `i`, it extends the ending position `j` incrementally, accumulating the sum of readings. It counts a period whenever:
1. The number of readings in the period is at most `maxReadings`
2. The sum falls within the inclusive band `[minTotal, maxTotal]`

Early termination occurs when a period exceeds `maxReadings`, since extending further will only violate the constraint.

## Correctness

The implementation correctly handles:
- **Invalid constraints**: Returns 0 immediately if `maxReadings <= 0`
- **Edge cases**: Empty readings array, single readings, all-zero readings
- **Negative values**: Uses 64-bit signed arithmetic for both readings and totals
- **Inclusive bounds**: Checks `sum >= minTotal && sum <= maxTotal`
- **Length constraints**: Breaks the inner loop when `j - i + 1 > maxReadings`, preventing longer periods from being considered

The test suite covers:
- Negative readings and negative total bounds
- Zero-length and invalid limits (`maxReadings <= 0`)
- Cases where ignoring `maxReadings` would yield incorrect results (e.g., a long sequence that sums to the target but exceeds the reading limit)
- Duplicate totals across different periods

## Complexity

**Time**: O(n²) where n is the number of readings. Each starting position generates up to `min(n, maxReadings)` inner iterations.

**Space**: O(1) auxiliary space (only counters and accumulators).

## Assumptions

- The input `readings` array may be empty (returns 0)
- Meter readings are represented as signed integers that fit in int32
- Totals are computed using 64-bit signed arithmetic to avoid overflow
- All inclusive bounds are valid (no validation of `minTotal <= maxTotal`)
- A "reviewable period" is any uninterrupted run of readings meeting the constraints; there is no requirement for periods to be non-overlapping
