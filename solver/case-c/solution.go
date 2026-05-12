package casec

import "sort"

func NetInventory(expr string) string {
	atomCounts := make(map[string]int)

	i := 0
	sign := 1 // First term is addition

	for i < len(expr) {
		// Parse a formula term
		term := ""
		for i < len(expr) && expr[i] != '+' && expr[i] != '-' {
			term += string(expr[i])
			i++
		}

		// Parse the formula and get atom counts
		formulaCounts := parseFormula(term)
		for atom, count := range formulaCounts {
			atomCounts[atom] += sign * count
		}

		// Read next operator
		if i < len(expr) {
			if expr[i] == '+' {
				sign = 1
			} else {
				sign = -1
			}
			i++
		}
	}

	// Check for negative final counts
	for _, count := range atomCounts {
		if count < 0 {
			return "INVALID"
		}
	}

	// Sort atoms and build result
	var atoms []string
	for atom, count := range atomCounts {
		if count > 0 {
			atoms = append(atoms, atom)
		}
	}
	sort.Strings(atoms)

	result := ""
	for _, atom := range atoms {
		result += atom
		count := atomCounts[atom]
		if count > 1 {
			// Convert count to string
			s := ""
			temp := count
			for temp > 0 {
				s = string(rune('0'+(temp%10))) + s
				temp /= 10
			}
			result += s
		}
	}

	return result
}

func parseFormula(formula string) map[string]int {
	result := make(map[string]int)
	var stack []map[string]int
	stack = append(stack, result)

	i := 0
	for i < len(formula) {
		if formula[i] == '(' {
			newMap := make(map[string]int)
			stack = append(stack, newMap)
			i++
		} else if formula[i] == ')' {
			i++
			// Parse multiplier after )
			multiplier := parseNumber(formula, &i)
			if multiplier == 0 {
				multiplier = 1
			}

			// Pop and merge with multiplier
			current := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			parent := stack[len(stack)-1]
			for atom, count := range current {
				parent[atom] += count * multiplier
			}
		} else if isUpperCase(formula[i]) {
			// Parse atom name and count
			atom := parseAtom(formula, &i)
			count := parseNumber(formula, &i)
			if count == 0 {
				count = 1
			}

			current := stack[len(stack)-1]
			current[atom] += count
		} else {
			i++
		}
	}

	return stack[0]
}

func isUpperCase(b byte) bool {
	return b >= 'A' && b <= 'Z'
}

func isLowerCase(b byte) bool {
	return b >= 'a' && b <= 'z'
}

func parseAtom(s string, i *int) string {
	atom := string(s[*i])
	*i++
	for *i < len(s) && isLowerCase(s[*i]) {
		atom += string(s[*i])
		*i++
	}
	return atom
}

func parseNumber(s string, i *int) int {
	if *i >= len(s) || s[*i] < '0' || s[*i] > '9' {
		return 0
	}
	num := 0
	for *i < len(s) && s[*i] >= '0' && s[*i] <= '9' {
		num = num*10 + int(s[*i]-'0')
		*i++
	}
	return num
}
