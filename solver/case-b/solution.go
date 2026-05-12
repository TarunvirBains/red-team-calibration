package caseb

type EntitlementLedger struct {
	credits map[int]int
}

func NewEntitlementLedger() EntitlementLedger {
	return EntitlementLedger{
		credits: make(map[int]int),
	}
}

func (e *EntitlementLedger) Grant(left int, right int) {
	if left >= right {
		return
	}
	for id := left; id < right; id++ {
		e.credits[id]++
	}
}

func (e *EntitlementLedger) Revoke(left int, right int) {
	if left >= right {
		return
	}
	for id := left; id < right; id++ {
		if e.credits[id] > 0 {
			e.credits[id]--
		}
	}
}

func (e *EntitlementLedger) Eligible(left int, right int, minCredits int) bool {
	if left >= right {
		return true
	}
	if minCredits <= 0 {
		return true
	}
	for id := left; id < right; id++ {
		if e.credits[id] < minCredits {
			return false
		}
	}
	return true
}
