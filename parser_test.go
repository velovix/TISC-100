package main

import "testing"

func TestParser(t *testing.T) {
	empty := newNodePort()
	emptyAny := newAnyPort(empty, empty, empty, empty)
	ex := newExecutionNode(empty, empty, empty, empty,
		emptyAny.lastUsedPort, emptyAny)

	// Create a scanner with the test code
	scan := newScanner()
	scan.add("nop\nMyLabel: add 14\njmp MyLabel\n")

	// Lex the test code
	lex := newLexer(scan)
	lex.lex()

	// Parse the tokens
	newParser(lex)

	if len(ex.instructions) < 3 {
		t.Error("parser created fewer instructions than expected")
	} else if len(ex.instructions) > 3 {
		t.Error("parser created more instructions than expected")
	}

	if line, ok := ex.labels["MYLABEL"]; ok {
		t.Error("parser failed to find a label")
	} else if line != 1 {
		t.Error("parser found the label, but didn't point it at the correct line: expected 1, found", line)
	}

	if _, ok := ex.instructions[0].(*nop); !ok {
		t.Error("parser failed to create an instruction of type nop")
	}

	if ins, ok := ex.instructions[1].(*add); !ok {
		t.Error("parser failed to create an instruction of type add")
	} else if ins.source == nil {
		t.Error("parser failed to parse the first argument of the add instruction")
	} else if n := ins.source.readNum(); n != number(15) {
		t.Error("the value of the first argument in the add instruction is incorrect: expected 15, found", n)
	}

	if ins, ok := ex.instructions[2].(*jmp); !ok {
		t.Error("parser failed to create an instruction of type jmp")
	} else if ins.l != "MYLABEL" {
		t.Error("parser read the jmp label incorrectly: expected 'MYLABEL', found '" + ins.l + "'")
	}
}
