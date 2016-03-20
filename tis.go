package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func main() {
	// Load the machine config file
	machConfig, err := newMachineConfig("./machine.json")
	if err != nil {
		fmt.Println("Error parsing machine.json:", err)
		os.Exit(1)
	}

	// Create a machine from the config information
	mach, err := newMachine(machConfig)
	if err != nil {
		fmt.Println("Error assembling TIS-100:", err)
	}

	// Load a source file for each executable node
	for y, row := range mach.nodes {
		for x, elem := range row {
			switch t := elem.(type) {
			case *executionNode:
				// The node is an execution node, so it may have a source file
				// associated with it

				file := strconv.Itoa(x) + "-" + strconv.Itoa(y) + ".tis"

				// Try to load a source file for the node
				data, err := ioutil.ReadFile(file)
				if err != nil && !os.IsNotExist(err) {
					fmt.Println("Error opening code for node "+file+":\n", err)
				}

				// Create a scanner from the code
				scan := newScanner()
				scan.add(string(data))
				scan.add("\n") // Add a newline to the end of the file in case one isn't there

				// Lex tokens out of the code
				lex := newLexer(scan)
				err = lex.lex()
				if err != nil {
					fmt.Println("Error opening code for node "+file+":\n", err)
				}

				// Parse the tokens
				parse := newParser(lex)
				err = parse.parse(t)
				if err != nil {
					fmt.Println("Error opening code for node "+file+":\n", err)
				}
			}
		}
	}

	// Start the machine
	mach.start()

	<-mach.stopSignal
}
