package casee

import "testing"

func TestEmptyInputs(t *testing.T) {
	tests := []struct {
		name          string
		sessionEffort []int
		staffCapacity []int
		grantCount    int
		grantSize     int
		expected      int
	}{
		{"no sessions", []int{}, []int{1, 2}, 0, 0, 0},
		{"no staff", []int{1, 2}, []int{}, 0, 0, 0},
		{"both empty", []int{}, []int{}, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxCoveredSessions(tt.sessionEffort, tt.staffCapacity, tt.grantCount, tt.grantSize)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestExactFits(t *testing.T) {
	tests := []struct {
		name          string
		sessionEffort []int
		staffCapacity []int
		grantCount    int
		grantSize     int
		expected      int
	}{
		{"single exact match", []int{5}, []int{5}, 0, 0, 1},
		{"all exact matches", []int{1, 2, 3}, []int{3, 2, 1}, 0, 0, 3},
		{"no matches", []int{5, 6, 7}, []int{1, 2, 3}, 0, 0, 0},
		{"partial matches", []int{1, 2, 3}, []int{2, 4}, 0, 0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxCoveredSessions(tt.sessionEffort, tt.staffCapacity, tt.grantCount, tt.grantSize)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestDuplicateValues(t *testing.T) {
	tests := []struct {
		name          string
		sessionEffort []int
		staffCapacity []int
		grantCount    int
		grantSize     int
		expected      int
	}{
		{"duplicate sessions same capacity", []int{5, 5, 5}, []int{5, 5}, 0, 0, 2},
		{"duplicate staff same effort", []int{2, 2, 2}, []int{2, 2}, 0, 0, 2},
		{"all duplicates", []int{3, 3}, []int{3, 3}, 0, 0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxCoveredSessions(tt.sessionEffort, tt.staffCapacity, tt.grantCount, tt.grantSize)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGrantsNecessary(t *testing.T) {
	tests := []struct {
		name          string
		sessionEffort []int
		staffCapacity []int
		grantCount    int
		grantSize     int
		expected      int
	}{
		{"grant enables one match", []int{5}, []int{3}, 1, 2, 1},
		{"grants enable multiple matches", []int{5, 5}, []int{3, 3}, 2, 2, 2},
		{"grant barely makes it", []int{10}, []int{5}, 1, 5, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxCoveredSessions(tt.sessionEffort, tt.staffCapacity, tt.grantCount, tt.grantSize)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGrantsInsufficient(t *testing.T) {
	tests := []struct {
		name          string
		sessionEffort []int
		staffCapacity []int
		grantCount    int
		grantSize     int
		expected      int
	}{
		{"grant too small", []int{10}, []int{3}, 1, 2, 0},
		{"grants too small for gap", []int{5, 5}, []int{1, 1}, 2, 2, 0},
		{"partial coverage with insufficient grants", []int{1, 10}, []int{2, 3}, 1, 2, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxCoveredSessions(tt.sessionEffort, tt.staffCapacity, tt.grantCount, tt.grantSize)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestExcessGrants(t *testing.T) {
	tests := []struct {
		name          string
		sessionEffort []int
		staffCapacity []int
		grantCount    int
		grantSize     int
		expected      int
	}{
		{"more grants than staff", []int{5}, []int{1, 1, 1}, 10, 10, 1},
		{"more grants than useful", []int{1, 2, 3}, []int{4, 5}, 10, 10, 2},
		{"no grants needed but plenty available", []int{1, 2}, []int{3, 4}, 100, 100, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxCoveredSessions(tt.sessionEffort, tt.staffCapacity, tt.grantCount, tt.grantSize)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestLocalGreedyOrder(t *testing.T) {
	// Test that proper sorting prevents locally attractive pairings from reducing overall coverage
	tests := []struct {
		name          string
		sessionEffort []int
		staffCapacity []int
		grantCount    int
		grantSize     int
		expected      int
	}{
		{
			"greedy matching optimal with unsorted input",
			[]int{1, 5, 5},
			[]int{1, 5},
			1, 4,
			2,
		},
		{
			"sorting ensures good matching over naive pairing",
			[]int{10, 10, 1},
			[]int{11, 5},
			1, 6,
			2, // After sort: sessions [1, 10, 10], staff [5, 11]
			// Session 1 with staff 5, session 10 with staff 11
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxCoveredSessions(tt.sessionEffort, tt.staffCapacity, tt.grantCount, tt.grantSize)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestInputNotMutated(t *testing.T) {
	sessions := []int{3, 1, 2}
	staff := []int{4, 2, 3}
	sessionsCopy := make([]int, len(sessions))
	staffCopy := make([]int, len(staff))
	copy(sessionsCopy, sessions)
	copy(staffCopy, staff)

	MaxCoveredSessions(sessions, staff, 0, 0)

	for i := range sessions {
		if sessions[i] != sessionsCopy[i] {
			t.Errorf("session slice was mutated at index %d", i)
		}
	}
	for i := range staff {
		if staff[i] != staffCopy[i] {
			t.Errorf("staff slice was mutated at index %d", i)
		}
	}
}

func TestComplexScenarios(t *testing.T) {
	tests := []struct {
		name          string
		sessionEffort []int
		staffCapacity []int
		grantCount    int
		grantSize     int
		expected      int
	}{
		{
			"grants cover gaps between two groups",
			[]int{2, 5, 8},
			[]int{1, 6},
			1, 3,
			2, // Staff 1 with grant (1+3=4) covers 2, staff 6 covers 5
		},
		{
			"strategic grant on higher capacity staff",
			[]int{5, 5},
			[]int{3, 4},
			1, 2,
			1, // Staff 3 or 4 with grant=5+1=5 covers one session; other can't match
		},
		{
			"zero grant size has no effect",
			[]int{5},
			[]int{3},
			10, 0,
			0,
		},
		{
			"multiple sessions multiple staff with grants",
			[]int{3, 6},
			[]int{2, 5},
			2, 2,
			2, // Staff 2 with grant (2+2=4) covers 3, staff 5 with grant (5+2=7) covers 6
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxCoveredSessions(tt.sessionEffort, tt.staffCapacity, tt.grantCount, tt.grantSize)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
