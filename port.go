package main

// port is a numberReadWriter that uses channels.
type port interface {
	numberReadWriter
	getChan() chan number
}

// nodePort is a two-way communication mechanism that is shared between nodes.
type nodePort struct {
	c chan number
}

// newNodePort creates a new node port.
func newNodePort() *nodePort {
	return &nodePort{
		c: make(chan number)}
}

// readNum reads a number from the other side of the port. This operation will
// block until the other node is sending on this port or indefinitely if there
// is no node on the other side.
func (np *nodePort) readNum() number {
	return <-np.c
}

// writeNum sends a number to the other side of the port. This operation will
// block until the other node is ready to recieve on this port or indefinitely
// if there is no node on the other side.
func (np *nodePort) writeNum(n number) {
	np.c <- n
}

// getChan returns the node port's underlying channel.
func (np *nodePort) getChan() chan number {
	return np.c
}

// anyPort is a pseudo-port that reads and writes to the first available port
// of the four given. Like other ports, it blocks until the operation can be
// completed. It also keeps track of what the last used port was for the LAST
// pseudo-port.
type anyPort struct {
	up, down, left, right port
	lastUsedPort          port
}

// newAnyPort creates a new ANY port that queries from the given ports. The
// slice of ports can be no larger than four.
func newAnyPort(up, down, left, right port) *anyPort {
	ap := &anyPort{
		up:    up,
		down:  down,
		left:  left,
		right: right}

	return ap
}

// readNum reads the first available number from the ports.
func (ap *anyPort) readNum() number {
	var n number
	select {
	case n = <-ap.up.getChan():
		ap.lastUsedPort = ap.up
	case n = <-ap.down.getChan():
		ap.lastUsedPort = ap.down
	case n = <-ap.left.getChan():
		ap.lastUsedPort = ap.left
	case n = <-ap.right.getChan():
		ap.lastUsedPort = ap.right
	}

	return n
}

// writeNum writes the given number to the first available port.
func (ap *anyPort) writeNum(n number) {
	select {
	case ap.up.getChan() <- n:
		ap.lastUsedPort = ap.up
	case ap.down.getChan() <- n:
		ap.lastUsedPort = ap.down
	case ap.left.getChan() <- n:
		ap.lastUsedPort = ap.left
	case ap.right.getChan() <- n:
		ap.lastUsedPort = ap.right
	}
}
