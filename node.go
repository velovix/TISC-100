package main

// node represents a node with four ends that can read and write from those
// ends
type node interface {
	getRight() port
	readRight() number
	writeRight(number)

	getLeft() port
	readLeft() number
	writeLeft(number)

	getUp() port
	readUp() number
	writeUp(number)

	getDown() port
	readDown() number
	writeDown(number)
}
