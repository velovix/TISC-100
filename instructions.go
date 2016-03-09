package main

import (
	"errors"
	"strconv"
)

// patternFromName returns an instruction pattern, which is a two-dimensional
// slice describing the tokens that an instruction accepts as arguments.
// The first dimension is an ordered list of lists of token types, which
// describes how many tokens are accepted. The second dimension is a list of
// token types that shows what kinds of tokens are acceptable in that position.
// If no known instruction with that name exists, an error is returned.
func patternFromName(insName string) ([][]tokenType, error) {
	switch insName {
	case "NOP":
		return [][]tokenType{}, nil
	case "MOV":
		return [][]tokenType{{tokenName, tokenNumber}, {tokenName}}, nil
	case "SWP":
		return [][]tokenType{}, nil
	case "SAV":
		return [][]tokenType{}, nil
	case "ADD":
		return [][]tokenType{{tokenName, tokenNumber}}, nil
	case "SUB":
		return [][]tokenType{{tokenName, tokenNumber}}, nil
	case "NEG":
		return [][]tokenType{}, nil
	case "JMP":
		return [][]tokenType{{tokenLabel}}, nil
	case "JEZ":
		return [][]tokenType{{tokenLabel}}, nil
	case "JNZ":
		return [][]tokenType{{tokenLabel}}, nil
	case "JGZ":
		return [][]tokenType{{tokenLabel}}, nil
	case "JLZ":
		return [][]tokenType{{tokenLabel}}, nil
	case "JRO":
		return [][]tokenType{{tokenName, tokenNumber}}, nil
	default:
		return [][]tokenType{}, errors.New("invalid instruction " + insName)
	}
}

// instructionFromName returns an empty buidable instruction based on the
// given instruction name, or an error if no such instruction is known.
func instructionFromName(insName string) (instruction, error) {
	switch insName {
	case "NOP":
		return &nop{}, nil
	case "MOV":
		return &mov{}, nil
	case "SWP":
		return &swp{}, nil
	case "SAV":
		return &sav{}, nil
	case "ADD":
		return &add{}, nil
	case "SUB":
		return &sub{}, nil
	case "NEG":
		return &neg{}, nil
	case "JMP":
		return &jmp{}, nil
	case "JEZ":
		return &jez{}, nil
	case "JNZ":
		return &jnz{}, nil
	case "JGZ":
		return &jgz{}, nil
	case "JLZ":
		return &jlz{}, nil
	case "JRO":
		return &jro{}, nil
	default:
		return &nop{}, errors.New("invalid instruction " + insName)
	}
}

type instruction interface {
	setArg(data interface{}, place int)
}

type nop struct{}

func (nopIns *nop) setArg(data interface{}, place int) {
	panic("extra argument supplied for NOP command in place " + strconv.Itoa(place))
}

type mov struct {
	source numberReader
	dest   numberWriter
}

func (movIns *mov) setArg(data interface{}, place int) {
	var ok bool
	switch place {
	case 0:
		movIns.source, ok = data.(numberReader)
		if !ok {
			panic("invalid type supplied as argument " + strconv.Itoa(place) + " in MOV command")
		}
	case 1:
		movIns.dest, ok = data.(numberWriter)
		if !ok {
			panic("invalid type supplied as argument " + strconv.Itoa(place) + " in MOV command")
		}
	default:
		panic("extra argument supplied for MOV command in place " + strconv.Itoa(place))
	}
}

type swp struct{}

func (swpIns *swp) setArg(data interface{}, place int) {
	panic("extra argument supplied for SWP command in place " + strconv.Itoa(place))
}

type sav struct{}

func (savIns *sav) setArg(data interface{}, place int) {
	panic("extra argument supplied for SAV command in place " + strconv.Itoa(place))
}

type add struct {
	source numberReader
}

func (addIns *add) setArg(data interface{}, place int) {
	var ok bool
	switch place {
	case 0:
		addIns.source, ok = data.(numberReader)
		if !ok {
			panic("invalid type supplied as argument " + strconv.Itoa(place) + " in ADD command")
		}
	default:
		panic("extra argument supplied for ADD command in place " + strconv.Itoa(place))
	}
}

type sub struct {
	source numberReader
}

func (subIns *sub) setArg(data interface{}, place int) {
	var ok bool
	switch place {
	case 0:
		subIns.source, ok = data.(numberReader)
		if !ok {
			panic("invalid type supplied as argument " + strconv.Itoa(place) + " in SUB command")
		}
	default:
		panic("extra argument supplied for SUB command in place " + strconv.Itoa(place))
	}
}

type neg struct{}

func (negIns *neg) setArg(data interface{}, place int) {
	panic("extra argument supplied for NEG command in place " + strconv.Itoa(place))
}

type jmp struct {
	l string
}

func (jmpIns *jmp) setArg(data interface{}, place int) {
	var ok bool
	switch place {
	case 0:
		jmpIns.l, ok = data.(string)
		if !ok {
			panic("invalid type supplied as argument " + strconv.Itoa(place) + " in JMP command")
		}
	default:
		panic("extra argument supplied for JMP command in place " + strconv.Itoa(place))
	}
}

type jez struct {
	l string
}

func (jezIns *jez) setArg(data interface{}, place int) {
	var ok bool
	switch place {
	case 0:
		jezIns.l, ok = data.(string)
		if !ok {
			panic("invalid type supplied as argument " + strconv.Itoa(place) + " in JEZ command")
		}
	default:
		panic("extra argument supplied for JEZ command in place " + strconv.Itoa(place))
	}
}

type jnz struct {
	l string
}

func (jnzIns *jnz) setArg(data interface{}, place int) {
	var ok bool
	switch place {
	case 0:
		jnzIns.l, ok = data.(string)
		if !ok {
			panic("invalid type supplied as argument " + strconv.Itoa(place) + " in JNZ command")
		}
	default:
		panic("extra argument supplied for JNZ command in place " + strconv.Itoa(place))
	}
}

type jgz struct {
	l string
}

func (jgzIns *jgz) setArg(data interface{}, place int) {
	var ok bool
	switch place {
	case 0:
		jgzIns.l, ok = data.(string)
		if !ok {
			panic("invalid type supplied as argument " + strconv.Itoa(place) + " in JGZ command")
		}
	default:
		panic("extra argument supplied for JGZ command in place " + strconv.Itoa(place))
	}
}

type jlz struct {
	l string
}

func (jlzIns *jlz) setArg(data interface{}, place int) {
	var ok bool
	switch place {
	case 0:
		jlzIns.l, ok = data.(string)
		if !ok {
			panic("invalid type supplied as argument " + strconv.Itoa(place) + " in JLZ command")
		}
	default:
		panic("extra argument supplied for JLZ command in place " + strconv.Itoa(place))
	}
}

type jro struct {
	source numberReader
}

func (jroIns *jro) setArg(data interface{}, place int) {
	var ok bool
	switch place {
	case 0:
		jroIns.source, ok = data.(numberReader)
		if !ok {
			panic("invalid type supplied as argument " + strconv.Itoa(place) + "in JRO command")
		}
	default:
		panic("extra argument supplied for JRO command in place " + strconv.Itoa(place))
	}
}
