package main

type executionNode struct {
	up, down, left, right, last port
	any                         *anyPort
	acc, bak                    numberReadWriter

	labels       map[string]int
	instructions []instruction
}

func newExecutionNode(up, down, left, right, last port, any *anyPort) *executionNode {
	return &executionNode{
		up:           up,
		down:         down,
		left:         left,
		right:        right,
		last:         last,
		any:          any,
		acc:          newRegister(0),
		bak:          newRegister(0),
		labels:       make(map[string]int),
		instructions: make([]instruction, 0)}
}

func (en *executionNode) start(stopped chan struct{}) {
	// Don't start running if the excution node is empty
	if len(en.instructions) == 0 {
		return
	}

	// Keep running the instruction until a stop signal is recieved
	i := 0 // Instruction position
	for {
		nextInstruction := en.instructions[i]

		switch ins := nextInstruction.(type) {
		case *nop:
			// Do nothing
		case *mov:
			// Move data from the source into the destination
			ins.dest.writeNum(ins.source.readNum())
			i++
		case *swp:
			// Swap what's in ACC with BAK
			temp := en.acc.readNum()
			en.acc.writeNum(en.bak.readNum())
			en.bak.writeNum(temp)
			i++
		case *sav:
			// Save the content of ACC to BAK
			en.bak.writeNum(en.acc.readNum())
			i++
		case *add:
			// Add source to ACC
			sum := addNum(en.acc.readNum(), int(ins.source.readNum()))
			en.acc.writeNum(number(sum))
			i++
		case *sub:
			// Sub source from ACC
			diff := subtractNum(en.acc.readNum(), int(ins.source.readNum()))
			en.acc.writeNum(number(diff))
			i++
		case *neg:
			// Negate ACC
			negated := -int(en.acc.readNum())
			en.acc.writeNum(number(negated))
			i++
		case *jmp:
			// Jump execution to the given label
			i = en.labels[ins.l]
		case *jez:
			// Jump execution to the given label if ACC is zero
			if en.acc.readNum() == 0 {
				i = en.labels[ins.l]
			}
		case *jnz:
			// Jump execution to the given label if ACC is not zero
			if en.acc.readNum() != 0 {
				i = en.labels[ins.l]
			}
		case *jgz:
			// Jump execution to the given label if ACC is greater than zero
			if en.acc.readNum() > 0 {
				i = en.labels[ins.l]
			}
		case *jlz:
			// Jump execution to the given label if ACC is less than zero
			if en.acc.readNum() < 0 {
				i = en.labels[ins.l]
			}
		case *jro:
			// Move execution by the given offset unconditionally
			i += int(ins.source.readNum())
		default:
			panic("unimplemented instruction")
		}

		// Wrap execution around to the beginning if need be
		i = i % len(en.instructions)
	}

	stopped <- struct{}{}
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

func (en *executionNode) getUp() port {
	return en.up
}

func (en *executionNode) getDown() port {
	return en.down
}

func (en *executionNode) getLeft() port {
	return en.left
}

func (en *executionNode) getRight() port {
	return en.right
}
