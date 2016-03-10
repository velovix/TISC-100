package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type consoleIn struct {
	c chan number
}

func newConsoleIn() *consoleIn {
	var cin consoleIn

	cin.c = make(chan number)

	// Start the goroutine that will feed console input to the console in
	// channel.
	go func() {
		stdin := bufio.NewReader(os.Stdin)

		// Keep reading user input until the end of time
		for {
			// Read a line of console input
			input, err := stdin.ReadString('\n')
			if err != nil {
				fmt.Println("Failure to read input:", err)
				continue
			}
			input = input[:len(input)-1] // Trim off newline

			// Convert the input to an integer
			inputInt, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid integer value")
				continue
			}

			// Convert the integer to a number
			inputNum := newNumber(inputInt)
			if inputInt != int(inputNum) {
				fmt.Println("Given integer is outside TIS-100 number bounds")
				continue
			}

			// Feed the results into the console in channel
			cin.c <- inputNum
		}
	}()

	return &cin
}

func (cin *consoleIn) readNum() number {
	return <-cin.c
}

func (cin *consoleIn) writeNum(n number) {
	panic("attempt to write to console input")
}

func (cin *consoleIn) getChan() chan number {
	return cin.c
}

type consoleOut struct{}

func newConsoleOut() *consoleOut {
	return &consoleOut{}
}

func (c *consoleOut) writeNum(n number) {
	fmt.Println(n)
}

func (c *consoleOut) readNum() number {
	panic("attempt to read from console output")
}

func (cout *consoleOut) getChan() chan number {
	return nil
}
