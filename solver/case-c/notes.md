# Design

The solution parses a chemical inventory expression as a sequence of formulas separated by `+` and `-` operators. The top-level expression parser scans for formulas and operators, maintaining a sign (addition or subtraction). For each formula, a recursive descent parser builds an atom count map, supporting atoms, multi-digit counts, and nested parenthetical groups. Counts are accumulated in a global map with their respective signs applied. After processing all formulas, negative counts are detected and return `"INVALID"`. The final result is constructed by sorting atoms lexicographically and appending counts.

# Correctness

The implementation correctly handles:
- Simple atoms: `H`, `Na`, `Cl`, `Mg`
- Atom counts: `H2`, `O10`
- Parenthetical groups: `(OH)2`
- Nested groups: `Ca((OH)2)2`
- Expression operators: `+` (default) and `-`
- Atoms with zero net count are omitted from output
- Lexicographic sort ensures canonical output
- Negative final counts trigger `"INVALID"` return

# Complexity

- **Time**: O(n + m log m) where n is expression length and m is the number of unique atoms
- **Space**: O(m) for the atom count map

# Assumptions

- All input formulas are syntactically valid
- Top-level operators are only `+` and `-`
- Atom names start with uppercase followed by zero or more lowercase letters
- Multi-digit counts are allowed
- Missing leading sign indicates addition
