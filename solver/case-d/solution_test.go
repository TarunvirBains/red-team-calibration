package cased

import "testing"

func TestEngagementLoadScore_Empty(t *testing.T) {
	result := EngagementLoadScore([]int{})
	if result != 0 {
		t.Errorf("empty input: expected 0, got %d", result)
	}
}

func TestEngagementLoadScore_SingleElement(t *testing.T) {
	result := EngagementLoadScore([]int{5})
	expected := int64(25) // min=5, sum=5, score=5*5=25
	if result != expected {
		t.Errorf("single element [5]: expected %d, got %d", expected, result)
	}
}

func TestEngagementLoadScore_AllZeros(t *testing.T) {
	result := EngagementLoadScore([]int{0, 0, 0})
	expected := int64(0) // min=0, sum=0, score=0*0=0
	if result != expected {
		t.Errorf("all zeros: expected %d, got %d", expected, result)
	}
}

func TestEngagementLoadScore_SingleZero(t *testing.T) {
	result := EngagementLoadScore([]int{0})
	expected := int64(0) // min=0, sum=0, score=0*0=0
	if result != expected {
		t.Errorf("single zero: expected %d, got %d", expected, result)
	}
}

func TestEngagementLoadScore_ZerosAndNonZeros(t *testing.T) {
	result := EngagementLoadScore([]int{0, 5, 0})
	// Run 1: [0] → min=0, sum=0, score=0
	// Run 2: [5] → min=5, sum=5, score=25
	// Run 3: [0] → min=0, sum=0, score=0
	// Total: 25
	expected := int64(25)
	if result != expected {
		t.Errorf("[0,5,0]: expected %d, got %d", expected, result)
	}
}

func TestEngagementLoadScore_RepeatedEqualValues(t *testing.T) {
	result := EngagementLoadScore([]int{3, 3, 3})
	// Run 1: [3, 3, 3] → min=3, sum=9, score=3*9=27
	expected := int64(27)
	if result != expected {
		t.Errorf("[3,3,3]: expected %d, got %d", expected, result)
	}
}

func TestEngagementLoadScore_MixedHighsAndLows(t *testing.T) {
	result := EngagementLoadScore([]int{5, 5, 3, 3, 3, 2})
	// Run 1: [5, 5] → min=5, sum=10, score=5*10=50
	// Run 2: [3, 3, 3] → min=3, sum=9, score=3*9=27
	// Run 3: [2] → min=2, sum=2, score=2*2=4
	// Total: 50+27+4=81
	expected := int64(81)
	if result != expected {
		t.Errorf("[5,5,3,3,3,2]: expected %d, got %d", expected, result)
	}
}

func TestEngagementLoadScore_AlternatingValues(t *testing.T) {
	result := EngagementLoadScore([]int{1, 2, 1, 2, 1})
	// Run 1: [1] → score=1*1=1
	// Run 2: [2] → score=2*2=4
	// Run 3: [1] → score=1*1=1
	// Run 4: [2] → score=2*2=4
	// Run 5: [1] → score=1*1=1
	// Total: 1+4+1+4+1=11
	expected := int64(11)
	if result != expected {
		t.Errorf("[1,2,1,2,1]: expected %d, got %d", expected, result)
	}
}

func TestEngagementLoadScore_LargeValuesSingleRun(t *testing.T) {
	result := EngagementLoadScore([]int{999999999, 999999999})
	// dayValue = 999999999
	// runSum = (999999999 + 999999999) % MOD = 1999999998 % 1000000007 = 999999991
	// minValue = 999999999 % MOD = 999999999
	// score = (999999999 * 999999991) % MOD
	minVal := int64(999999999)
	sumVal := int64(1999999998 % MOD)
	expected := (minVal * sumVal) % int64(MOD)
	if result != expected {
		t.Errorf("[999999999,999999999]: expected %d, got %d", expected, result)
	}
}

func TestEngagementLoadScore_LargeCountWithModulo(t *testing.T) {
	// Create 100 consecutive 10M loads
	daily := make([]int, 100)
	for i := range daily {
		daily[i] = 10000000
	}
	result := EngagementLoadScore(daily)
	// dayValue = 10000000
	// runSum = (100 * 10000000) % MOD = 1000000000 % 1000000007 = 999999993
	// minValue = 10000000 % MOD = 10000000
	// score = (10000000 * 999999993) % MOD
	minVal := int64(10000000)
	sumVal := int64(1000000000 % MOD)
	expected := (minVal * sumVal) % int64(MOD)
	if result != expected {
		t.Errorf("100 x 10M: expected %d, got %d", expected, result)
	}
}

func TestEngagementLoadScore_MultipleRunsWithLargeNumbers(t *testing.T) {
	result := EngagementLoadScore([]int{999999999, 999999999, 1, 1})
	// Run 1: [999999999, 999999999]
	//   runSum = (999999999 + 999999999) % MOD = 999999991
	//   score = (999999999 * 999999991) % MOD = 998999936
	// Run 2: [1, 1]
	//   runSum = (1 + 1) % MOD = 2
	//   score = (1 * 2) % MOD = 2
	// Total = (998999936 + 2) % MOD = 998999938
	minVal1 := int64(999999999)
	sumVal1 := int64((999999999 + 999999999) % MOD)
	score1 := (minVal1 * sumVal1) % int64(MOD)
	score2 := int64(2)
	expected := (score1 + score2) % int64(MOD)
	if result != expected {
		t.Errorf("[999999999,999999999,1,1]: expected %d, got %d", expected, result)
	}
}

func TestEngagementLoadScore_OneOfEach(t *testing.T) {
	result := EngagementLoadScore([]int{1})
	expected := int64(1) // 1*1=1
	if result != expected {
		t.Errorf("[1]: expected %d, got %d", expected, result)
	}
}
