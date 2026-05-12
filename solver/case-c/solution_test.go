package casec

import "testing"

func TestNetInventory(t *testing.T) {
	tests := []struct {
		expr     string
		expected string
	}{
		// Examples from prompt
		{"H2O+NaCl-Na", "ClH2O"},
		{"Mg(OH)2-H2", "MgO2"},
		{"H2-O3", "INVALID"},

		// Nested groups
		{"Ca(OH)2", "CaH2O2"},

		// Subtraction
		{"NaCl-Cl", "Na"},

		// Zero net counts (omitted)
		{"H2+H2-H2-H2", ""},
		{"O2-O2", ""},

		// Negative final counts
		{"H-H2", "INVALID"},
		{"O-O2", "INVALID"},

		// Lexicographic ordering
		{"Cl+Na+H", "ClHNa"},

		// Multi-digit counts
		{"H10+O10", "H10O10"},

		// Single atom
		{"H", "H"},
		{"Na", "Na"},

		// Addition (no sign at start)
		{"H+O", "HO"},

		// Complex nested parentheses
		{"C(C(C)2)2", "C7"},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			result := NetInventory(tt.expr)
			if result != tt.expected {
				t.Errorf("NetInventory(%q) = %q, want %q", tt.expr, result, tt.expected)
			}
		})
	}
}
