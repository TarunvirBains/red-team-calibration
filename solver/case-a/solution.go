package casea

func AuditMeterReadings(readings []int, minTotal int64, maxTotal int64, maxReadings int) int64 {
	if maxReadings <= 0 {
		return 0
	}

	count := int64(0)
	n := len(readings)

	for i := 0; i < n; i++ {
		sum := int64(0)
		for j := i; j < n; j++ {
			// Check if period exceeds maxReadings
			if j-i+1 > maxReadings {
				break
			}

			// Add current reading to sum using 64-bit arithmetic
			sum += int64(readings[j])

			// Check if sum is within the inclusive bounds
			if sum >= minTotal && sum <= maxTotal {
				count++
			}
		}
	}

	return count
}
