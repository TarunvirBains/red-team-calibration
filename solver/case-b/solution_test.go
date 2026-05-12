package caseb

import "testing"

func TestEmptySpans(t *testing.T) {
	e := NewEntitlementLedger()

	e.Grant(5, 5)
	e.Revoke(10, 10)

	if !e.Eligible(15, 15, 1) {
		t.Error("empty span should return true for Eligible")
	}
}

func TestGrantAndEligible(t *testing.T) {
	e := NewEntitlementLedger()

	e.Grant(0, 5)

	if !e.Eligible(0, 5, 1) {
		t.Error("after grant, should be eligible for 1 credit")
	}

	if e.Eligible(0, 5, 2) {
		t.Error("should not be eligible for 2 credits")
	}
}

func TestOverlappingGrants(t *testing.T) {
	e := NewEntitlementLedger()

	e.Grant(0, 10)
	e.Grant(5, 15)

	if !e.Eligible(0, 5, 1) {
		t.Error("accounts 0-4 should have 1 credit")
	}
	if e.Eligible(0, 5, 2) {
		t.Error("accounts 0-4 should not have 2 credits")
	}

	if !e.Eligible(5, 10, 2) {
		t.Error("accounts 5-9 should have 2 credits")
	}
	if e.Eligible(5, 10, 3) {
		t.Error("accounts 5-9 should not have 3 credits")
	}

	if !e.Eligible(10, 15, 1) {
		t.Error("accounts 10-14 should have 1 credit")
	}
	if e.Eligible(10, 15, 2) {
		t.Error("accounts 10-14 should not have 2 credits")
	}
}

func TestRevoke(t *testing.T) {
	e := NewEntitlementLedger()

	e.Grant(0, 5)
	e.Grant(0, 5)

	if !e.Eligible(0, 5, 2) {
		t.Error("should have 2 credits")
	}

	e.Revoke(0, 5)

	if !e.Eligible(0, 5, 1) {
		t.Error("should have 1 credit after revoke")
	}
	if e.Eligible(0, 5, 2) {
		t.Error("should not have 2 credits after revoke")
	}
}

func TestPartialRevoke(t *testing.T) {
	e := NewEntitlementLedger()

	e.Grant(0, 10)
	e.Revoke(5, 15)

	if !e.Eligible(0, 5, 1) {
		t.Error("accounts 0-4 should have 1 credit")
	}

	if e.Eligible(5, 10, 1) {
		t.Error("accounts 5-9 should have 0 credits")
	}
	if !e.Eligible(5, 10, 0) {
		t.Error("0 credits should be eligible for 0 minCredits")
	}

	if e.Eligible(10, 15, 1) {
		t.Error("accounts 10-14 should have 0 credits")
	}
}

func TestRevokeNeverBelowZero(t *testing.T) {
	e := NewEntitlementLedger()

	e.Grant(0, 5)
	e.Revoke(0, 5)
	e.Revoke(0, 5)
	e.Revoke(0, 5)

	if e.Eligible(0, 5, 1) {
		t.Error("should have 0 credits after multiple revokes")
	}
}

func TestAdjacentSpans(t *testing.T) {
	e := NewEntitlementLedger()

	e.Grant(0, 5)
	e.Grant(5, 10)

	if !e.Eligible(0, 5, 1) {
		t.Error("first span should have 1 credit")
	}

	if !e.Eligible(5, 10, 1) {
		t.Error("second span should have 1 credit")
	}

	if !e.Eligible(0, 10, 1) {
		t.Error("spanning both should have 1 credit each")
	}
}

func TestMinCreditsThreshold(t *testing.T) {
	e := NewEntitlementLedger()

	e.Grant(0, 5)
	e.Grant(0, 5)

	if !e.Eligible(0, 5, 0) {
		t.Error("minCredits=0 should return true")
	}
	if !e.Eligible(0, 5, -1) {
		t.Error("minCredits=-1 should return true")
	}

	if !e.Eligible(0, 5, 1) {
		t.Error("minCredits=1 should return true")
	}

	if !e.Eligible(0, 5, 2) {
		t.Error("minCredits=2 should return true")
	}

	if e.Eligible(0, 5, 3) {
		t.Error("minCredits=3 should return false")
	}
}

func TestLongOperationSequence(t *testing.T) {
	e := NewEntitlementLedger()

	e.Grant(0, 100)
	e.Grant(50, 150)
	e.Revoke(75, 125)
	e.Grant(100, 200)
	e.Revoke(0, 50)

	if e.Eligible(0, 50, 1) {
		t.Error("0-50 should have been revoked to 0")
	}

	if !e.Eligible(50, 75, 1) {
		t.Error("50-75 should have 2 credits")
	}

	if !e.Eligible(75, 100, 1) {
		t.Error("75-100 should have 1 credit")
	}

	if !e.Eligible(100, 125, 1) {
		t.Error("100-125 should have 1 credit")
	}

	if !e.Eligible(125, 150, 2) {
		t.Error("125-150 should have 2 credits")
	}

	if !e.Eligible(150, 200, 1) {
		t.Error("150-200 should have 1 credit")
	}
}
