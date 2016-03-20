package main

import (
	"fmt"
)

type executionNode struct {
	up, down, left, right, last port
	any                         *anyPort
	acc, bak                    numberReadWriter

	labels       map[string]int
	instructions []instruction

	name string
}

func newExecutionNode(name string, up, down, left, right, last port, any *anyPort) *executionNode {
	return &executionNode{
		name:         name,
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

func (en *executionNode) String() string {
	return en.name
}

func (en *executionNode) start(stopped chan struct{}) {
	// Don't start running if the excution node is empty
	if len(en.instructions) == 0 {
		return
	}

	var ok bool

	// Keep running the instruction until a stop signal is recieved
	i := 0 // Instruction position
execution:
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
			if i, ok = en.labels[ins.l]; !ok {
				fmt.Println("unknown label '" + ins.l + "'")
				break execution // The label doesn't exist, so we halt execution
			}
		case *jez:
			// Jump execution to the given label if ACC is zero
			if en.acc.readNum() == 0 {
				if i, ok = en.labels[ins.l]; !ok {
					fmt.Println("unknown label '" + ins.l + "'")
					break execution // The label doesn't exist, so we halt execution
				}
			} else {
				i++
			}
		case *jnz:
			// Jump execution to the given label if ACC is not zero
			if en.acc.readNum() != 0 {
				if i, ok = en.labels[ins.l]; !ok {
					fmt.Println("unknown label '" + ins.l + "'")
					break execution // The label doesn't exist, so we halt execution
				}
			} else {
				i++
			}
		case *jgz:
			// Jump execution to the given label if ACC is greater than zero
			if en.acc.readNum() > 0 {
				if i, ok = en.labels[ins.l]; !ok {
					fmt.Println("unknown label '" + ins.l + "'")
					break execution // The label doesn't exist, so we halt execution
				}
			} else {
				i++
			}
		case *jlz:
			// Jump execution to the given label if ACC is less than zero
			if en.acc.readNum() < 0 {
				if i, ok = en.labels[ins.l]; !ok {
					fmt.Println("unknown label '" + ins.l + "'")
					break execution // The label doesn't exist, so we halt execution
				}
			} else {
				i++
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
