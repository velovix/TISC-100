package main

import (
	"strings"
	"unicode"
)

// tokenType represents the various types a token can be.
type tokenType int

//go:generate stringer -type=tokenType
const (
	_           tokenType = iota
	tokenName             // Token is the name of something
	tokenLabel            // Token is a label
	tokenNumber           // Token is a number literal
)

// token is a single lexical token.
type token struct {
	tType        tokenType
	startingChar char
	data         string
}

// lexerState represents the type of token the lexer is currently constructing
type lexerState int

//go:generate stringer -type=lexerState
const (
	lexerStateNone lexerState = iota
	lexerStateNameOrLabel
	lexerStateNumber
)

type lexer struct {
	scan      scanner
	state     lexerState
	tokens    []token
	currToken int
}

// newLexer creates a new lexer that reads characters from the given scanner
// object.
func newLexer(scan scanner) lexer {
	return lexer{
		scan:   scan,
		tokens: make([]token, 0, 20)}
}

// lex starts lexing the input from the scanner object.
func (l *lexer) lex() error {
	var data string
	var startingChar char

	// Loop through ever character
	for character, hasNextChar := l.scan.next(); hasNextChar; character, hasNextChar = l.scan.next() {
		switch l.state {
		case lexerStateNone:
			// The lexer is waiting for a new state

			if unicode.IsLetter(character.c) {
				// A letter can either mean a name of some kind or a label
				l.state = lexerStateNameOrLabel
				startingChar = character
				data += string(character.c)
			} else if unicode.IsDigit(character.c) || character.c == '-' {
				// A digit or a negative sign means a numeric literal
				l.state = lexerStateNumber
				startingChar = character
				data += string(character.c)
			} else if !unicode.IsSpace(character.c) {
				// If the character isn't a letter, number, or space, it isn't valid
				return newParseError("unexpected character '"+string(character.c)+"'", character)
			}
		case lexerStateNameOrLabel:
			// The lexer expects more characters or a sign that the token is finished

			if unicode.IsLetter(character.c) {
				// Another letter, so the name or label is still being constructed
				data += string(character.c)
			} else if character.c == ':' {
				// A colon denotes the end of a label
				l.tokens = append(l.tokens, token{
					tType:        tokenLabel,
					startingChar: startingChar,
					data:         strings.ToUpper(data)})
				data = ""
				l.state = lexerStateNone
			} else if unicode.IsSpace(character.c) {
				// A space denotes the end of a name
				l.tokens = append(l.tokens, token{
					tType:        tokenName,
					startingChar: startingChar,
					data:         strings.ToUpper(data)})
				data = ""
				l.state = lexerStateNone // Reset the lexer state
			} else {
				// An invalid character
				return newParseError("unexpected character '"+string(character.c)+"'", character)
			}
		case lexerStateNumber:
			// The lexer expects more numbers or a sign that the token is finished

			if unicode.IsDigit(character.c) {
				// Another digit, so the number is still being constructed
				data += string(character.c)
			} else if unicode.IsSpace(character.c) {
				// A space denotes the end of the number
				l.tokens = append(l.tokens, token{
					tType:        tokenNumber,
					startingChar: startingChar,
					data:         data})
				data = ""
				l.state = lexerStateNone // Reset the lexer state
			} else {
				// An invalid character
				return newParseError("unexpected character '"+string(character.c)+"'", character)
			}
		}
	}

	return nil
}

// next returns the next parsed token, or an empty token and false if no more
// exist.
func (l *lexer) next() (token, bool) {
	if l.currToken < len(l.tokens) {
		l.currToken++
		return l.tokens[l.currToken-1], true
	} else {
		return token{}, false
	}
}
