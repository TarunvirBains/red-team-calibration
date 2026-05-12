# EarliestChecklistArrival Solution Notes

## Design

The solution uses a queue-based walk to find the shortest path from the top-left to the bottom-right checkpoint. The key insight is that since each move takes exactly 1 minute and time only increases monotonically, queue order naturally explores states in order of increasing time, claiming that the first arrival at the destination is the earliest possible.

The implementation maintains:
- A queue of states (row, col, time) representing the worker's position and arrival time
- A visited map tracking (row, col) → earliest_arrival_time to avoid reprocessing cells
- Four cardinal directions for neighbor exploration

**Why visited by (row, col) only:** Once we reach a cell at its earliest possible time, any later arrival at that same cell can only lead to equal or worse overall arrival times at the destination. This is because time strictly increases, and reaching an intermediate cell earlier allows more time to navigate to the destination.

Input validation detects:
- Empty arrays
- Jagged arrays (inconsistent row widths)
- Single checkpoint (base case: already at destination)

## Correctness

The implementation is correct because:

1. **Completeness:** BFS explores all reachable cells. If the destination is reachable, we will find it.

2. **Optimality:** BFS processes states in order of increasing time. The first time we reach the destination is guaranteed to be via a shortest path (minimum time).

3. **Constraint satisfaction:** We only move to a neighbor if `newTime ≥ readyAt[neighbor]`, ensuring cells are entered at valid times.

4. **No-wait enforcement:** The implementation doesn't create artificial delay states. We only progress time by moving to valid neighbors. If all neighbors are blocked, we cannot proceed from that cell, which follows the contract's "no waiting" rule.

5. **Unreachability detection:** If the queue empties without reaching the destination, we correctly return -1.

Edge cases handled:
- Starting cell readyAt value doesn't matter (we're already there at time 0)
- One-checkpoint grids return 0 (destination is starting point)
- Large readyAt values don't cause overflow under the chosen integer type
- Blocked paths correctly return -1

## Complexity

**Time Complexity:** O(rows × cols)
- Each cell is visited at most once in the optimal path
- The visited map ensures we process each (row, col) coordinate only once
- Each neighbor check is O(1)

**Space Complexity:** O(rows × cols)
- Visited map stores at most rows × cols entries
- Queue can contain at most O(rows × cols) elements in the worst case (if entire grid is reachable)

This is efficient even for large layouts. The implementation avoids exploring infinite loops or redundant paths through the visited tracking.

## Assumptions

1. **Integer times:** The readyAt values and computed arrival times fit within Go's `int` range. While large readyAt values are handled, overflow is assumed not to occur under the stated contract.

2. **Rectangular layout:** The input is either empty, jagged (invalid), or represents a valid m×n rectangular layout.

3. **Movement costs:** Every move between adjacent cells (up/down/left/right) costs exactly 1 minute. Diagonal movement is not allowed.

4. **Starting position:** The worker always starts at (0, 0) at time 0, and this is always valid regardless of readyAt[0][0].

5. **Destination:** The destination is always (rows-1, cols-1). There is no alternative goal cell.

6. **No waiting semantics:** "Waiting is not allowed" means we cannot stay in a checkpoint. If all neighbors are blocked, the path is dead, and alternatives are explored by the queue.

7. **Readiness is permanent:** Once a checkpoint becomes ready (time ≥ readyAt), it remains accessible for all future times.
