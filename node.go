package main

// node represents a node with four ends that can read and write from those
// ends
type node interface {
	getRight() port
	getLeft() port
	getUp() port
	getDown() port

	start(stopped chan struct{})
}
