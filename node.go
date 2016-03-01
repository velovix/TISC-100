package main

// node represents a node with four ends that can read and write from those
// ends
type node interface {
	getRight() numberReadWriter
	readRight() number
	writeRight(number)

	getLeft() numberReadWriter
	readLeft() number
	writeLeft(number)

	getUp() numberReadWriter
	readUp() number
	writeUp(number)

	getDown() numberReadWriter
	readDown() number
	writeDown(number)
}
