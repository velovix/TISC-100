package main

import "testing"

// TestScanning tests the scanner for correct character order, character count,
// and proper positioning and line count.
func TestScanning(t *testing.T) {
	scan := newScanner()

	scan.add("test\n!")

	var c char
	var hasNext bool
	expected := char{c: 't', pos: 0, line: 0}
	if c, hasNext = scan.next(); c != expected {
		t.Error("expected scanner to emit", expected, "but gave", c)
	} else if !hasNext {
		t.Error("scanner ran out of characters prematurely")
	}

	expected = char{c: 'e', pos: 1, line: 0}
	if c, hasNext = scan.next(); c != expected {
		t.Error("expected scanner to emit", expected, "but gave", c)
	} else if !hasNext {
		t.Error("scanner ran out of characters prematurely")
	}

	expected = char{c: 's', pos: 2, line: 0}
	if c, hasNext = scan.next(); c != expected {
		t.Error("expected scanner to emit", expected, "but gave", c)
	} else if !hasNext {
		t.Error("scanner ran out of characters prematurely")
	}

	expected = char{c: 't', pos: 3, line: 0}
	if c, hasNext = scan.next(); c != expected {
		t.Error("expected scanner to emit", expected, "but gave", c)
	} else if !hasNext {
		t.Error("scanner ran out of characters prematurely")
	}

	expected = char{c: '\n', pos: 4, line: 0}
	if c, hasNext = scan.next(); c != expected {
		t.Error("expected scanner to emit", expected, "but gave", c)
	} else if !hasNext {
		t.Error("scanner ran out of characters prematurely")
	}

	expected = char{c: '!', pos: 0, line: 1}
	if c, hasNext = scan.next(); c != expected {
		t.Error("expected scanner to emit", expected, "but gave", c)
	} else if !hasNext {
		t.Error("scanner ran out of characters prematurely")
	}

	if _, hasNext = scan.next(); hasNext {
		t.Error("scanner still claims to have characters after it should have run out")
	}
}

// TestScannerIgoresComments checks that comment characters are not emitted by
// the scanner.
func TestScannerIgnoresComments(t *testing.T) {
	scan := newScanner()

	scan.add("hello\n#world\n! # Test")

	expected := []rune("hello\n! ")
	expectedPos := 0
	for c, hasNext := scan.next(); hasNext; c, hasNext = scan.next() {
		if c.c != expected[expectedPos] {
			t.Error("expected scanner to emit a '" + string(expected[expectedPos]) + "' but got a '" + string(c.c) + "'")
		}
		expectedPos++
	}
}
