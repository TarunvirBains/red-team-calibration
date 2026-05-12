package casef

func EarliestChecklistArrival(readyAt [][]int) int {
	// Check for empty input
	if len(readyAt) == 0 {
		return -1
	}

	// Check dimensions and validate
	rows := len(readyAt)
	cols := len(readyAt[0])

	if cols == 0 {
		return -1
	}

	// Detect jagged array
	for _, row := range readyAt {
		if len(row) != cols {
			return -1
		}
	}

	// Single checkpoint: already at destination
	if rows == 1 && cols == 1 {
		return 0
	}

	// Find maximum readyAt value to bound exploration time
	maxReady := 0
	for _, row := range readyAt {
		for _, val := range row {
			if val > maxReady {
				maxReady = val
			}
		}
	}

	// Maximum useful exploration time: can't reach destination faster than
	// the grid allows plus the maximum readyAt constraint
	// In practice, if we haven't found a path by this time, it's unreachable
	maxTime := rows + cols + maxReady

	// BFS to find earliest arrival
	type State struct {
		row  int
		col  int
		time int
	}

	queue := []State{{0, 0, 0}}
	// Track visited (row, col, time) states
	visited := make(map[[3]int]bool)
	visited[[3]int{0, 0, 0}] = true

	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for len(queue) > 0 {
		state := queue[0]
		queue = queue[1:]

		// Prevent exploration beyond reasonable time limit
		if state.time >= maxTime {
			continue
		}

		// Try all four cardinal directions
		for _, dir := range directions {
			newRow := state.row + dir[0]
			newCol := state.col + dir[1]
			newTime := state.time + 1

			// Check bounds
			if newRow < 0 || newRow >= rows || newCol < 0 || newCol >= cols {
				continue
			}

			// Check if we can enter this cell at this time
			if newTime < readyAt[newRow][newCol] {
				continue
			}

			// Check if we've visited this exact state
			key := [3]int{newRow, newCol, newTime}
			if visited[key] {
				continue
			}

			// Check if we've reached the destination
			if newRow == rows-1 && newCol == cols-1 {
				return newTime
			}

			visited[key] = true
			queue = append(queue, State{newRow, newCol, newTime})
		}
	}

	// Destination unreachable
	return -1
}
