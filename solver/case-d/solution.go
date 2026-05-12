package cased

const MOD = 1000000007

// EngagementLoadScore calculates notification fatigue by scoring each
// continuous run of equal daily loads: (min_load * sum_of_loads) % MOD.
func EngagementLoadScore(daily []int) int64 {
	if len(daily) == 0 {
		return 0
	}

	var totalScore int64 = 0
	i := 0

	for i < len(daily) {
		// Start of a new run of consecutive equal values
		dayValue := int64(daily[i])
		var runSum int64 = 0

		// Accumulate sum of all consecutive equal values
		for i < len(daily) && int64(daily[i]) == dayValue {
			runSum = (runSum + dayValue) % MOD
			i++
		}

		// Score for this run: (min * sum) % MOD
		// Min of the run is dayValue itself (all are equal)
		minValue := dayValue % MOD
		runScore := (minValue * runSum) % MOD
		totalScore = (totalScore + runScore) % MOD
	}

	return totalScore
}
