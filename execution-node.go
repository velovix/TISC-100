package main

type executionNode struct {
	up, down, left, right, any, last numberReadWriter
	acc, bak                         numberReadWriter

	labels       map[string]int
	instructions []instruction
}

func newExecutionNode(up, down, left, right, any, last numberReadWriter) *executionNode {
	return &executionNode{
		up:           up,
		down:         down,
		left:         left,
		right:        right,
		any:          any,
		last:         last,
		acc:          newRegister(0),
		bak:          newRegister(0),
		labels:       make(map[string]int),
		instructions: make([]instruction, 0)}
}

func (en *executionNode) start(stop chan struct{}) {
	for {
		for _, val := range en.instructions {
			switch ins := val.(type) {
			case *nop:
				// Do nothing
			case *mov:
				ins.dest.writeNum(ins.source.readNum())
			case *swp:
				temp := en.acc.readNum()
				en.acc.writeNum(en.bak.readNum())
				en.bak.writeNum(temp)
			case *sav:
				en.bak.writeNum(en.acc.readNum())
			case *add:
				sum := int(en.acc.readNum()) + int(ins.source.readNum())
				en.acc.writeNum(number(sum))
			case *sub:
				diff := int(en.acc.readNum()) - int(ins.source.readNum())
				en.acc.writeNum(number(diff))
			case *neg:
				negated := -int(en.acc.readNum())
				en.acc.writeNum(number(negated))
			default:
				panic("unimplemented instruction")
			}
		}
	}
}

func (en *executionNode) readRight() number {
	return en.up.readNum()
}

func (en *executionNode) readLeft() number {
	return en.left.readNum()
}

func (en *executionNode) readUp() number {
	return en.up.readNum()
}

func (en *executionNode) readDown() number {
	return en.down.readNum()
}

func (en *executionNode) writeRight(n number) {
	en.right.writeNum(n)
}

func (en *executionNode) writeLeft(n number) {
	en.left.writeNum(n)
}

func (en *executionNode) writeUp(n number) {
	en.up.writeNum(n)
}

func (en *executionNode) writeDown(n number) {
	en.down.writeNum(n)
}

func (en *executionNode) getUp() numberReadWriter {
	return en.up
}

func (en *executionNode) getDown() numberReadWriter {
	return en.down
}

func (en *executionNode) getLeft() numberReadWriter {
	return en.left
}

func (en *executionNode) getRight() numberReadWriter {
	return en.right
}
