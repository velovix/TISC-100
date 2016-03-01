package main

// port is a two-way communication mechanism that is shared between nodes.
type port struct {
	c chan number
}

// newPort creates a new port.
func newPort() *port {
	return &port{
		c: make(chan number)}
}

// readNum reads a number from the other side of the port. This operation will
// block until the other node is sending on this port or indefinitely if there
// is no node on the other side.
func (p *port) readNum() number {
	return <-p.c
}

// writeNum sends a number to the other side of the port. This operation will
// block until the other node is ready to recieve on this port or indefinitely
// if there is no node on the other side.
func (p *port) writeNum(n number) {
	p.c <- n
}

// anyPort is a pseudo-port that reads and writes to the first available port
// of the four given. Like other ports, it blocks until the operation can be
// completed. It also keeps track of what the last used port was for the LAST
// pseudo-port.
type anyPort struct {
	up, down, left, right *port
	lastUsedPort          *port
}

// newAnyPort creates a new ANY port that queries from the given ports.
func newAnyPort(up, down, left, right *port) *anyPort {
	return &anyPort{
		up:    up,
		down:  down,
		left:  left,
		right: right}
}

// readNum reads the first available number from the ports.
func (ap *anyPort) readNum() number {
	var n number
	select {
	case n = <-ap.up.c:
		ap.lastUsedPort = ap.up
	case n = <-ap.down.c:
		ap.lastUsedPort = ap.down
	case n = <-ap.left.c:
		ap.lastUsedPort = ap.left
	case n = <-ap.right.c:
		ap.lastUsedPort = ap.right
	}

	return n
}

// writeNum writes the given number to the first available port.
func (ap *anyPort) writeNum(n number) {
	select {
	case ap.up.c <- n:
		ap.lastUsedPort = ap.up
	case ap.down.c <- n:
		ap.lastUsedPort = ap.down
	case ap.left.c <- n:
		ap.lastUsedPort = ap.left
	case ap.right.c <- n:
		ap.lastUsedPort = ap.right
	}
}
