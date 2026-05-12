package casea

import "testing"

func TestAuditMeterReadings(t *testing.T) {
	tests := []struct {
		name       string
		readings   []int
		minTotal   int64
		maxTotal   int64
		maxReadings int
		expected   int64
	}{
		{
			name:        "simple case with multiple valid periods",
			readings:    []int{1, 2, 3},
			minTotal:    3,
			maxTotal:    5,
			maxReadings: 3,
			expected:    3, // [1,2]=3, [2,3]=5, [3]=3
		},
		{
			name:        "negative readings",
			readings:    []int{-5, 10, -3},
			minTotal:    -3,
			maxTotal:    10,
			maxReadings: 3,
			expected:    5, // [-5,10]=5, [-5,10,-3]=2, [10]=10, [10,-3]=7, [-3]=-3
		},
		{
			name:        "maxReadings zero returns zero",
			readings:    []int{1, 2, 3},
			minTotal:    1,
			maxTotal:    10,
			maxReadings: 0,
			expected:    0,
		},
		{
			name:        "maxReadings negative returns zero",
			readings:    []int{1, 2, 3},
			minTotal:    1,
			maxTotal:    10,
			maxReadings: -1,
			expected:    0,
		},
		{
			name:        "empty readings",
			readings:    []int{},
			minTotal:    0,
			maxTotal:    10,
			maxReadings: 5,
			expected:    0,
		},
		{
			name:        "maxReadings prevents full sequence from counting",
			readings:    []int{1, 2, 3, 4, 5},
			minTotal:    15,
			maxTotal:    15,
			maxReadings: 4,
			expected:    0, // [1,2,3,4,5]=15 but has 5 readings > maxReadings=4
		},
		{
			name:        "duplicate totals in different windows",
			readings:    []int{1, 2, 3, 2, 1},
			minTotal:    3,
			maxTotal:    3,
			maxReadings: 5,
			expected:    3, // [1,2]=3, [3]=3, [2,1]=3
		},
		{
			name:        "single element",
			readings:    []int{5},
			minTotal:    5,
			maxTotal:    5,
			maxReadings: 1,
			expected:    1, // [5]=5
		},
		{
			name:        "no valid periods",
			readings:    []int{1, 1, 1},
			minTotal:    10,
			maxTotal:    20,
			maxReadings: 3,
			expected:    0,
		},
		{
			name:        "negative minTotal and maxTotal",
			readings:    []int{-5, -3, 2},
			minTotal:    -8,
			maxTotal:    -2,
			maxReadings: 3,
			expected:    4, // [-5]=-5, [-5,-3]=-8, [-5,-3,2]=-6, [-3]=-3
		},
		{
			name:        "all readings negative",
			readings:    []int{-2, -3, -5},
			minTotal:    -10,
			maxTotal:    -5,
			maxReadings: 3,
			expected:    4, // [-2,-3]=-5, [-2,-3,-5]=-10, [-3,-5]=-8, [-5]=-5
		},
		{
			name:        "zero readings",
			readings:    []int{0, 0, 5},
			minTotal:    0,
			maxTotal:    5,
			maxReadings: 3,
			expected:    6, // [0]=0, [0,0]=0, [0,0,5]=5, [0]=0, [0,5]=5, [5]=5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AuditMeterReadings(tt.readings, tt.minTotal, tt.maxTotal, tt.maxReadings)
			if got != tt.expected {
				t.Errorf("got %d, want %d", got, tt.expected)
			}
		})
	}
}
