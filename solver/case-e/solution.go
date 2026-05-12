package casee

import "sort"

func MaxCoveredSessions(sessionEffort []int, staffCapacity []int, grantCount int, grantSize int) int {
	if len(sessionEffort) == 0 || len(staffCapacity) == 0 {
		return 0
	}

	// Make copies to avoid mutating input
	sessions := make([]int, len(sessionEffort))
	copy(sessions, sessionEffort)
	sort.Ints(sessions)

	staff := make([]int, len(staffCapacity))
	copy(staff, staffCapacity)
	sort.Ints(staff)

	// Try all possible grant assignments and compute maximum coverage
	// For smaller input sizes, iterate through which staff get grants.
	maxCovered := 0

	// For each possible subset of staff to grant (up to grantCount)
	// we try to maximize matching
	for mask := 0; mask < (1 << uint(len(staff))); mask++ {
		grantUsed := 0
		for i := 0; i < len(staff); i++ {
			if (mask & (1 << uint(i))) != 0 {
				grantUsed++
			}
		}
		if grantUsed > grantCount {
			continue
		}

		// Compute effective capacities with this grant assignment
		capacities := make([]int, len(staff))
		for i := 0; i < len(staff); i++ {
			capacities[i] = staff[i]
			if (mask & (1 << uint(i))) != 0 {
				capacities[i] += grantSize
			}
		}

		// Greedily match sessions to staff
		usedLocal := make([]bool, len(staff))
		covered := 0

		for _, effort := range sessions {
			// Find the best unused staff member for this session
			// (smallest capacity that can still cover)
			bestIdx := -1
			for j := 0; j < len(staff); j++ {
				if !usedLocal[j] && capacities[j] >= effort {
					if bestIdx == -1 || capacities[j] < capacities[bestIdx] {
						bestIdx = j
					}
				}
			}
			if bestIdx != -1 {
				usedLocal[bestIdx] = true
				covered++
			}
		}

		if covered > maxCovered {
			maxCovered = covered
		}
	}

	return maxCovered
}
