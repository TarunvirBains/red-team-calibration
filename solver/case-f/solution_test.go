package casef

import "testing"

func TestEmpty(t *testing.T) {
	result := EarliestChecklistArrival([][]int{})
	if result != -1 {
		t.Errorf("Empty input: expected -1, got %d", result)
	}
}

func TestEmptyColumn(t *testing.T) {
	result := EarliestChecklistArrival([][]int{{}})
	if result != -1 {
		t.Errorf("Empty column: expected -1, got %d", result)
	}
}

func TestJagged(t *testing.T) {
	result := EarliestChecklistArrival([][]int{{0, 1}, {2}})
	if result != -1 {
		t.Errorf("Jagged input: expected -1, got %d", result)
	}
}

func TestOneCheckpoint(t *testing.T) {
	result := EarliestChecklistArrival([][]int{{0}})
	if result != 0 {
		t.Errorf("One checkpoint: expected 0, got %d", result)
	}
}

func TestOneCheckpointNonzeroReady(t *testing.T) {
	result := EarliestChecklistArrival([][]int{{5}})
	if result != 0 {
		t.Errorf("One checkpoint with readyAt=5: expected 0 (already there at time 0), got %d", result)
	}
}

func TestSimple2x2AllReady(t *testing.T) {
	readyAt := [][]int{{0, 0}, {0, 0}}
	result := EarliestChecklistArrival(readyAt)
	if result != 2 {
		t.Errorf("Simple 2x2: expected 2, got %d", result)
	}
}

func TestBlockedFirstMoves(t *testing.T) {
	// Both immediate neighbors are blocked
	readyAt := [][]int{{0, 10}, {10, 0}}
	result := EarliestChecklistArrival(readyAt)
	if result != -1 {
		t.Errorf("Blocked first moves: expected -1, got %d", result)
	}
}

func TestDestinationBlocked(t *testing.T) {
	// Destination ready at time 5; can reach by looping
	readyAt := [][]int{{0, 0}, {0, 5}}
	result := EarliestChecklistArrival(readyAt)
	// Path with loop: (0,0)->(1,0)->(0,0)->(0,1)->(1,0)->(0,0)->(1,0)->(1,1)
	// Better: (0,0)->(1,0)->(0,0)->(0,1)->(1,0)->(1,1)
	// Times:   0  ->  1  ->  2  ->  3  ->  4  ->  5? No, 4<5
	// Even better: add more loops to reach at time >= 5
	// (0,0)->(1,0)->(0,0)->(1,0)->(0,0)->(0,1)->(1,1): 0->1->2->3->4->5->6
	if result != 6 {
		t.Errorf("Destination blocked: expected 6, got %d", result)
	}
}

func TestDetourAroundBlockade(t *testing.T) {
	// Must go around a blocked cell
	readyAt := [][]int{{0, 10, 0}, {0, 0, 0}}
	result := EarliestChecklistArrival(readyAt)
	// Path: (0,0) -> (1,0) -> (1,1) -> (1,2)
	// Times: 0 -> 1 -> 2 -> 3 (destination is (1,2))
	if result != 3 {
		t.Errorf("Detour around blockade: expected 3, got %d", result)
	}
}

func TestDelayedDestinationReachable(t *testing.T) {
	// Destination becomes ready at time of arrival
	readyAt := [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 4}}
	result := EarliestChecklistArrival(readyAt)
	// Shortest path to (2,2): 4 moves
	// (0,0) -> (1,0) -> (2,0) -> (2,1) -> (2,2)
	// Arrives at time 4, ready at 4: accessible
	if result != 4 {
		t.Errorf("Delayed destination reachable at time 4: expected 4, got %d", result)
	}
}

func TestDelayedDestinationUnreachable(t *testing.T) {
	// Destination ready at 5; shortest path is 4, but can loop to reach at time >= 5
	readyAt := [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 5}}
	result := EarliestChecklistArrival(readyAt)
	// By revisiting cells, we can reach (2,2) at time 6:
	// (0,0)->(1,0)->(0,0)->(1,0)->(2,0)->(2,1)->(2,2): 0->1->2->3->4->5->6
	if result != 6 {
		t.Errorf("Delayed destination (ready 5): expected 6, got %d", result)
	}
}

func TestLargeReadyAtUnreachable(t *testing.T) {
	// Test handling of large readyAt values
	readyAt := [][]int{{0, 0}, {0, 1000000000}}
	result := EarliestChecklistArrival(readyAt)
	if result != -1 {
		t.Errorf("Large readyAt (unreachable): expected -1, got %d", result)
	}
}

func TestLargeReadyAtReachable(t *testing.T) {
	// Destination ready at a very large time, but reachable at that time
	readyAt := [][]int{{0, 0}, {0, 10}}
	result := EarliestChecklistArrival(readyAt)
	// (0,0) -> (1,0) at time 1
	// (1,0) -> (1,1) at time 2, ready at 10, so blocked
	if result != -1 {
		t.Errorf("Large readyAt (2 < 10): expected -1, got %d", result)
	}
}

func TestComplexPath(t *testing.T) {
	// Multi-step path with various delays requiring revisits
	// Grid: {{0, 2, 3}, {0, 3, 0}} - destination is (1,2)
	readyAt := [][]int{{0, 2, 3}, {0, 3, 0}}
	result := EarliestChecklistArrival(readyAt)
	// Path: (0,0) at 0 -> (1,0) at 1 -> (0,0) at 2 -> (0,1) at 3 -> (1,1) at 4 -> (1,2) at 5
	// Revisiting (0,0) allows us to reach (0,1) when it's ready at time 2
	if result != 5 {
		t.Errorf("Complex path: expected 5, got %d", result)
	}
}

func TestArrivingLaterUnblocks(t *testing.T) {
	// Case where forced delay allows reaching previously blocked cell
	readyAt := [][]int{{0, 2}, {0, 4}}
	result := EarliestChecklistArrival(readyAt)
	// Path: (0,0) -> (1,0) -> (0,0) -> (0,1) -> (1,1)
	// Times: 0 -> 1 -> 2 -> 3 -> 4
	// At time 3, (0,1) is ready (readyAt=2)
	// At time 4, (1,1) is ready (readyAt=4)
	if result != 4 {
		t.Errorf("Arriving later unblocks destination: expected 4, got %d", result)
	}
}

func TestWideGrid(t *testing.T) {
	// Wider grid to test horizontal traversal
	readyAt := [][]int{{0, 0, 0, 0, 0}}
	result := EarliestChecklistArrival(readyAt)
	// (0,0) -> (0,1) -> (0,2) -> (0,3) -> (0,4)
	// Times: 0 -> 1 -> 2 -> 3 -> 4
	if result != 4 {
		t.Errorf("Wide grid: expected 4, got %d", result)
	}
}

func TestTallGrid(t *testing.T) {
	// Tall grid to test vertical traversal
	readyAt := [][]int{{0}, {0}, {0}, {0}, {0}}
	result := EarliestChecklistArrival(readyAt)
	// (0,0) -> (1,0) -> (2,0) -> (3,0) -> (4,0)
	// Times: 0 -> 1 -> 2 -> 3 -> 4
	if result != 4 {
		t.Errorf("Tall grid: expected 4, got %d", result)
	}
}
