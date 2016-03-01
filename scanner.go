package main

// char represents a single character and marks where the character is in the
// code.
type char struct {
	c    rune
	pos  int
	line int
}

// scanner takes in code as a string and emits it as individual tokens.
type scanner struct {
	chars         []char
	currChar      int
	ignoringChars bool
	currPos       int
	currLine      int
}

// newScanner creates a new scanner.
func newScanner() scanner {
	return scanner{
		chars: make([]char, 0, 256)}
}

// add appends characters from the given string to the scanner.
func (s *scanner) add(chars string) {
	for _, val := range chars {
		if val == '#' {
			// The beginning of a comment. Ignore everything until a newline
			s.ignoringChars = true
		} else if val == '\n' {
			// Go to the next line and reset cursor position on newline
			s.currLine++
			s.currPos = 0
			s.ignoringChars = false // Exit comment ignoring mode if we're there
		}

		if !s.ignoringChars {
			// Append the character to the input
			s.chars = append(s.chars, char{
				c:    val,
				pos:  s.currPos,
				line: s.currLine})
		}
		s.currPos++
	}
}

// next returns the next character if one exists, or an empty character and
// false.
func (s *scanner) next() (char, bool) {
	if s.currChar < len(s.chars) {
		s.currChar++
		return s.chars[s.currChar-1], true
	} else {
		return char{}, false
	}
}
