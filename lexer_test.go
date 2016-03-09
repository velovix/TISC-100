package main

import "testing"

// TestLexingNames tests the lexer for the ability to generate the correct
// number of name tokens with case insensitive data and the character position
// that they exist at in the text.
func TestLexingNames(t *testing.T) {
	// Create a new scanner with the test data
	scan := newScanner()
	scan.add("sub ACC\n\njMP myLabel\n")

	// Lex the data
	lex := newLexer(scan)
	err := lex.lex()
	if err != nil {
		t.Error(err)
	}

	// How the lexer should react for each token
	testCases := []struct {
		expected token
		hasNext  bool
	}{
		{
			hasNext: true,
			expected: token{
				tType:        tokenName,
				startingChar: char{c: 's', pos: 0, line: 0},
				data:         "SUB"}},
		{
			hasNext: true,
			expected: token{
				tType:        tokenName,
				startingChar: char{c: 'A', pos: 4, line: 0},
				data:         "ACC"}},

		{
			hasNext: true,
			expected: token{
				tType:        tokenName,
				startingChar: char{c: 'j', pos: 0, line: 2},
				data:         "JMP"}},
		{
			hasNext: true,
			expected: token{
				tType:        tokenName,
				startingChar: char{c: 'm', pos: 4, line: 2},
				data:         "MYLABEL"}},
		{
			hasNext:  false,
			expected: token{}}}

	// Check against each test case
	for _, testCase := range testCases {
		// Get the next token
		tok, hasNext := lex.next()
		// Check if the lexer's token count is as expected
		if hasNext != testCase.hasNext {
			if hasNext {
				t.Error("lexer has an unexpected extra token")
			} else {
				t.Error("lexer ran out of tokens prematurely")
			}
		}
		// Check if the value of the emitted token is as expected
		if tok != testCase.expected {
			t.Error("expected the token", testCase.expected, "but got", tok)
		}
	}
}

// TestLexingNumbers tests the lexer's ability to generate positive and negative
// number tokens with the correct data and character positions.
func TestLexingNumbers(t *testing.T) {
	// Create a new scanner with the test data
	scan := newScanner()
	scan.add("add 9\nadd -14\n")

	// Lex the data
	lex := newLexer(scan)
	err := lex.lex()
	if err != nil {
		t.Error(err)
	}

	// How the lexer should react for each token
	testCases := []struct {
		expected token
		hasNext  bool
	}{
		{
			hasNext: true,
			expected: token{
				tType:        tokenName,
				startingChar: char{c: 'a', pos: 0, line: 0},
				data:         "ADD"}},
		{
			hasNext: true,
			expected: token{
				tType:        tokenNumber,
				startingChar: char{c: '9', pos: 4, line: 0},
				data:         "9"}},
		{
			hasNext: true,
			expected: token{
				tType:        tokenName,
				startingChar: char{c: 'a', pos: 0, line: 1},
				data:         "ADD"}},
		{
			hasNext: true,
			expected: token{
				tType:        tokenNumber,
				startingChar: char{c: '-', pos: 4, line: 1},
				data:         "-14"}},
		{
			hasNext:  false,
			expected: token{}}}

	// Check against each test case
	for _, testCase := range testCases {
		// Get the next token
		tok, hasNext := lex.next()
		// Check if the lexer's token count is as expected
		if hasNext != testCase.hasNext {
			if hasNext {
				t.Error("lexer has an unexpected extra token")
			} else {
				t.Error("lexer ran out of tokens prematurely")
			}
		}
		// Check if the value of the emitted token is as expected
		if tok != testCase.expected {
			t.Error("expected the token", testCase.expected, "but got", tok)
		}
	}
}

// TestLexingLabels tests the lexer's ability to emit label tokens case-
// insensitively both when they're on their own lines and when the label shares
// the line with an instruction.
func TestLexingLabels(t *testing.T) {
	// Create a new scanner with the test data
	scan := newScanner()
	scan.add("nop\ntestLabel:\nnop\nanotherLabel: nop\n")

	// Lex the data
	lex := newLexer(scan)
	err := lex.lex()
	if err != nil {
		t.Error(err)
	}

	// How the lexer should react for each token
	testCases := []struct {
		expected token
		hasNext  bool
	}{
		{
			hasNext: true,
			expected: token{
				tType:        tokenName,
				startingChar: char{c: 'n', pos: 0, line: 0},
				data:         "NOP"}},
		{
			hasNext: true,
			expected: token{
				tType:        tokenLabel,
				startingChar: char{c: 't', pos: 0, line: 1},
				data:         "TESTLABEL"}},
		{
			hasNext: true,
			expected: token{
				tType:        tokenName,
				startingChar: char{c: 'n', pos: 0, line: 2},
				data:         "NOP"}},
		{
			hasNext: true,
			expected: token{
				tType:        tokenLabel,
				startingChar: char{c: 'a', pos: 0, line: 3},
				data:         "ANOTHERLABEL"}},
		{
			hasNext: true,
			expected: token{
				tType:        tokenName,
				startingChar: char{c: 'n', pos: 14, line: 3},
				data:         "NOP"}},

		{
			hasNext:  false,
			expected: token{}}}

	// Check against each test case
	for _, testCase := range testCases {
		// Get the next token
		tok, hasNext := lex.next()
		// Check if the lexer's token count is as expected
		if hasNext != testCase.hasNext {
			if hasNext {
				t.Error("lexer has an unexpected extra token")
			} else {
				t.Error("lexer ran out of tokens prematurely")
			}
		}
		// Check if the value of the emitted token is as expected
		if tok != testCase.expected {
			t.Error("expected the token", testCase.expected, "but got", tok)
		}
	}
}

// TestLexingInvalidSyntax tests the lexer's ability to fail when various
// syntactically invalid code is fed to it.
func TestLexingInvalidSyntax(t *testing.T) {
	// Create a new scanner with the test data
	scan := newScanner()
	scan.add("nop\n&\n")

	// Lex the erroneous code
	lex := newLexer(scan)
	err := lex.lex()
	if err == nil {
		t.Error("lexer didn't fail even though an invalid beginning to a token was given")
	}

	// Create a new scanner with the test data
	scan = newScanner()
	scan.add("add -14a")

	// Lex the erroneous code
	lex = newLexer(scan)
	err = lex.lex()
	if err == nil {
		t.Error("lexer didn't fail even though an invalid number was given")
	}
}
