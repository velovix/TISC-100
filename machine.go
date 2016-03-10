package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// machineConfig contains the configuration information of a single machine.
// This can be used to construct a machine object.
type machineConfig struct {
	Name      string     `json:"name"`
	Nodes     [][]string `json:"nodes"`
	ConsoleIn struct {
		Side string `json:"side"`
		Pos  int    `json:"Pos"`
	} `json:"consoleIn"`
	ConsoleOut struct {
		Side string `json:"side"`
		Pos  int    `json:"pos"`
	} `json:"consoleOut"`
}

// newMachineConfig creates a new machine configuration object based on the
// provided config file location.
func newMachineConfig(config string) (machineConfig, error) {
	var mc machineConfig

	// Read the data from the config file
	data, err := ioutil.ReadFile(config)
	if err != nil {
		return machineConfig{}, err
	}

	// Interpret the data as JSON and populate the machine configuration object
	err = json.Unmarshal(data, &mc)
	if err != nil {
		return machineConfig{}, err
	}

	// Make sure the given nodes create a rectangle
	nodeWidth := len(mc.Nodes[0])
	nodeHeight := len(mc.Nodes)
	for _, val := range mc.Nodes {
		if len(val) != nodeWidth {
			return machineConfig{}, errors.New("node array must form a rectangle")
		}
	}

	// Check console in position for validity
	switch mc.ConsoleIn.Side {
	case "top":
		fallthrough
	case "bottom":
		if mc.ConsoleIn.Pos < 0 || mc.ConsoleIn.Pos >= nodeWidth {
			return machineConfig{}, errors.New("consoleIn pos must be within the width of the node array")
		}
	case "right":
		fallthrough
	case "left":
		if mc.ConsoleIn.Pos < 0 || mc.ConsoleIn.Pos >= nodeHeight {
			return machineConfig{}, errors.New("consoleIn pos must be within the height of the node array")
		}
	default:
		return machineConfig{}, errors.New("consoleIn has an invalid side value")
	}

	// Check console out position for validity
	switch mc.ConsoleOut.Side {
	case "top":
		fallthrough
	case "bottom":
		if mc.ConsoleOut.Pos < 0 || mc.ConsoleOut.Pos >= nodeWidth {
			return machineConfig{}, errors.New("consoleOut pos must be within the width of the node array")
		}
	case "right":
		fallthrough
	case "left":
		if mc.ConsoleOut.Pos < 0 || mc.ConsoleOut.Pos >= nodeHeight {
			return machineConfig{}, errors.New("consoleOut pos must be within the height of the node array")
		}
	default:
		return machineConfig{}, errors.New("consoleOut has an invalid side value")
	}

	return mc, nil
}

// machine represents the TIS-100 instance. It is a collection of nodes.
type machine struct {
	nodes      [][]node
	stopSignal chan struct{}

	consoleIn  port
	consoleOut port
}

// newMachine creates a new machine from the given machine config . It
// creates empty nodes based on the configuration and wires them up to each
// other.
func newMachine(config machineConfig) (machine, error) {
	var m machine

	m.stopSignal = make(chan struct{})

	// Construct the console input and output
	m.consoleIn = newConsoleIn()
	m.consoleOut = newConsoleOut()

	// Construct an empty array of nodes based on the size of the nodes in the config
	m.nodes = make([][]node, len(config.Nodes))
	for i := range m.nodes {
		m.nodes[i] = make([]node, len(config.Nodes[0]))
	}

	nodeWidth := len(m.nodes[0])
	nodeHeight := len(m.nodes)

	for x, valX := range config.Nodes {
		for y, valY := range valX {
			// The node's ports default to ones that go nowhere
			var up, down, left, right port = newNodePort(), newNodePort(), newNodePort(), newNodePort()

			// See if this is the node that console in should be wired to
			switch config.ConsoleIn.Side {
			case "top":
				if y == 0 && x == config.ConsoleIn.Pos {
					up = m.consoleIn
				}
			case "bottom":
				if y == nodeHeight-1 && x == config.ConsoleIn.Pos {
					down = m.consoleIn
				}
			case "left":
				if x == 0 && y == config.ConsoleIn.Pos {
					left = m.consoleIn
				}
			case "right":
				if x == nodeWidth-1 && y == config.ConsoleIn.Pos {
					right = m.consoleIn
				}
			}

			// See if this is the node that console out should be wired to
			switch config.ConsoleOut.Side {
			case "top":
				if y == 0 && x == config.ConsoleOut.Pos {
					up = m.consoleOut
				}
			case "bottom":
				if y == nodeHeight-1 && x == config.ConsoleOut.Pos {
					down = m.consoleOut
				}
			case "left":
				if x == 0 && y == config.ConsoleOut.Pos {
					left = m.consoleOut
				}
			case "right":
				if x == nodeWidth-1 && y == config.ConsoleOut.Pos {
					right = m.consoleOut
				}
			}

			switch valY {
			case "e":
				// The node is an execution node

				any := newAnyPort(up, down, left, right)

				if y-1 >= 0 {
					// If there is a node above this node, connect it
					up = m.nodes[x][y-1].getDown()
				}
				if x-1 >= 0 {
					// If there is a node to the left of this node, connect it
					left = m.nodes[x-1][y].getRight()
				}

				m.nodes[x][y] = newExecutionNode(up, down, left, right, any.lastUsedPort, any)
			default:
				// The node is invalid

				return machine{}, errors.New("invalid node type '" + valY + "'")
			}
		}
	}

	return m, nil
}

// start starts all execution nodes.
func (m *machine) start() {
	for _, row := range m.nodes {
		for _, elem := range row {
			if execNode, ok := elem.(*executionNode); ok {
				go execNode.start(m.stopSignal)
			}
		}
	}
}
