package main

import (
	"fmt"
	"os"
)

func main() {
	machConfig, err := newMachineConfig("./machine.json")
	if err != nil {
		fmt.Println("Error parsing machine.json:", err)
		os.Exit(1)
	}

	_, err = newMachine(machConfig)
	if err != nil {
		fmt.Println("Error assembling TIS-100:", err)
	}
}
