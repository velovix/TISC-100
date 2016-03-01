package main

const (
	numberMaxValue = 999
	numberMinValue = -999
)

// number is an immutible integer value that follows the TIS-100's weird number
// capping rules. It should be treated as a read-only value and not be modified
// directly.
type number int

// numberReader describes an object that can act as a source of a number.
type numberReader interface {
	readNum() number
}

// numberWriter describes an object that can take in a number.
type numberWriter interface {
	writeNum(number)
}

// numberReadWriter describes an object that can act as a source of a number
// and take in a number.
type numberReadWriter interface {
	numberReader
	numberWriter
}

func capNumber(val number) number {
	if val > numberMaxValue {
		return number(numberMaxValue)
	} else if val < numberMinValue {
		return number(numberMinValue)
	}

	return val
}

// Returns a new valid TIS-100 number based on the given integer.
func newNumber(val int) number {
	return capNumber(number(val))
}

// addNum adds the given integer to the number and returns a valid TIS-100
// number.
func addNum(n number, val int) number {
	return newNumber(int(n) + val)
}

// subtractNum subtracts the given integer from the number and returns a
// valid TIS-100 number.
func subtractNum(n number, val int) number {
	return newNumber(int(n) - val)
}
