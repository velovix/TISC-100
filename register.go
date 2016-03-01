package main

// register is a basic number storage mechanism. It does not communicate
// between nodes.
type register struct {
	value number
}

// newRegister creates a new register holding the given value.
func newRegister(value int) *register {
	return &register{
		value: newNumber(value)}
}

// readNum returns the value the register is holding.
func (r *register) readNum() number {
	return r.value
}

// writeNum sets the value of the register to the given value.
func (r *register) writeNum(val number) {
	r.value = val
}

// nilRegister is a pseudo-register that always reads as a zero and discards any
// numbers that are written to it. There exists only one instance of this
// pseudo-register.
type nilRegister struct{}

var nilReg nilRegister

func (*nilRegister) readNum() number {
	return number(0)
}

func (*nilRegister) writeNum(val number) {}
