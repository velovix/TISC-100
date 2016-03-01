package main

import "strconv"

type parserState int

const (
	parserStateNone parserState = iota
	parserStateInstructionSpecific
)

type parser struct {
	lex   *lexer
	state parserState
}

func newParser(lex *lexer) parser {
	return parser{
		lex: lex}
}

func (p *parser) parse(exNode executionNode) error {
	var currPattern [][]tokenType
	var patternPos int
	var builder instruction
	var argPos int
	var instructionCnt int

	// Loop through every lexical token
	for t, hasNext := p.lex.next(); hasNext; t, hasNext = p.lex.next() {
		switch p.state {
		case parserStateNone:
			// The parser doesn't know what to expect

			switch t.tType {
			case tokenName:
				// The next token is a name, meaning that it's the start of an
				// instruction

				// Try to find the pattern for the given instruction so we know
				// how to parse it
				if val, err := patternFromName(t.data); err == nil {
					currPattern = val
					patternPos = 0
					// Get the buildable instruction from the name
					if val2, err := instructionFromName(t.data); err == nil {
						builder = val2
						argPos = 0
					} else {
						return newParseError(err.Error(), t.startingChar)
					}
					// Start parsing tokens according to the specific instruction
					p.state = parserStateInstructionSpecific
				} else {
					return newParseError(err.Error(), t.startingChar)
				}
			case tokenLabel:
				// The next token is a label

				// Check that the label does not already exist
				if _, ok := exNode.labels[t.data]; ok {
					return newParseError("duplicate label '"+t.data+"'", t.startingChar)
				}

				// Set the label to point to the "address" of the next instruction
				exNode.labels[t.data] = instructionCnt
			case tokenNumber:
				// The next token is a number. No operations start with a
				// number, so this is an error.

				return newParseError("unexpected number '"+t.data+"'", t.startingChar)
			}
		case parserStateInstructionSpecific:
			// The parser is parsing an instruction

			// Check the current token against the list of acceptable tokens
			validToken := false
			for _, val := range currPattern[patternPos] {
				if t.tType == val {
					validToken = true
				}
			}
			if !validToken {
				return newParseError("invalid token '"+t.data+"'", t.startingChar)
			}

			// Check that the token is valid and feed it into the builder
			if t.tType == tokenName {
				switch t.data {
				case "ACC":
					builder.setArg(exNode.acc, argPos)
				case "BAK":
					builder.setArg(exNode.bak, argPos)
				case "NIL":
					builder.setArg(nilReg, argPos)
				case "LEFT":
					builder.setArg(exNode.left, argPos)
				case "RIGHT":
					builder.setArg(exNode.right, argPos)
				case "UP":
					builder.setArg(exNode.up, argPos)
				case "DOWN":
					builder.setArg(exNode.down, argPos)
				case "ANY":
					builder.setArg(exNode.any, argPos)
				case "LAST":
					builder.setArg(exNode.last, argPos)
				}
			} else if t.tType == tokenNumber {
				// Check that the token's data is a valid number if it is a
				// number
				val, err := strconv.Atoi(t.data)
				if err != nil {
					return newParseError("'"+t.data+"' can't be parsed as a number", t.startingChar)
				}
				if val != int(number(val)) {
					return newParseError("'"+t.data+"' falls outside the range of an acceptable TIS-100 number", t.startingChar)
				}
			}

			// Go to the next pattern or finish the instruction
			patternPos++
			if patternPos == len(currPattern) {
				exNode.instructions = append(exNode.instructions, builder)
				p.state = parserStateNone
				instructionCnt++
			}
		}
	}

	return nil
}
