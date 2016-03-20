package main

import (
	"sync"
)

type stackNode struct {
	up, down, left, right port
	any                   *anyPort
	head                  *stackNodeLink
	nextServedData        chan number
	sync.Mutex
}

type stackNodeLink struct {
	n    number
	prev *stackNodeLink
}

func newStackNode(up, down, left, right port) *stackNode {
	return &stackNode{
		up:             up,
		down:           down,
		left:           left,
		right:          right,
		nextServedData: make(chan number)}
}

// read returns number popped off the top of the stack. Since it doesn't matter
// what direction a read request comes from to the stack node, all reads use
// this method.
func (sn *stackNode) read() number {
	sn.Lock() // Request ownership of the stack node

	if sn.head == nil {
		// The stack is empty, so we relinquish ownership of the stack node and
		// listen on the channel for the next served piece of data
		sn.Unlock()
		return <-sn.nextServedData
	}

	defer sn.Unlock()
	// Pop the data off the top of the stack
	n := sn.head.n
	sn.head = sn.head.prev
	return n
}

// write puts a number on the top of the stack. Since it doesn't matter to the
// stack node what direction a write request comes from, all writes use this
// method.
func (sn *stackNode) write(n number) {
	sn.Lock() // Request ownership of the stack node
	defer sn.Unlock()

	// Check if a read request is waiting on an empty stack
	select {
	case sn.nextServedData <- n:
		return
	}

	// Add new data to the stop of the stack
	link := &stackNodeLink{
		n:    n,
		prev: sn.head}
	sn.head = link
}

func (sn *stackNode) start(stopped chan struct{}) {
	// Wait for incoming data
	go func(stopped chan struct{}) {
		for {
			var n number
			select {
			case n = <-sn.up.getChan():
			case n = <-sn.down.getChan():
			case n = <-sn.left.getChan():
			case n = <-sn.right.getChan():
			}

			sn.write(n)
		}
	}(stopped)

	// Wait for data requests
	go func(stopped chan struct{}) {
		for {
			select {
			case sn.up.getChan() <- sn.read():
			case sn.down.getChan() <- sn.read():
			case sn.left.getChan() <- sn.read():
			case sn.right.getChan() <- sn.read():
			}
		}
	}(stopped)

	<-stopped
}

func (sn *stackNode) getLeft() port {
	return sn.left
}

func (sn *stackNode) getRight() port {
	return sn.up
}

func (sn *stackNode) getUp() port {
	return sn.up
}

func (sn *stackNode) getDown() port {
	return sn.down
}
