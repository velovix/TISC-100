package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// machineConfig contains the configuration information of a single machine.
// This can be used to construct a machine object.
type machineConfig struct {
	Name  string     `json:"name"`
	Nodes [][]string `json:"nodes"`
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
	for _, val := range mc.Nodes {
		if len(val) != nodeWidth {
			return machineConfig{}, errors.New("node array must form a rectangle")
		}
	}

	return mc, nil
}

// machine represents the TIS-100 instance. It is a collection of nodes.
type machine struct {
	nodes      [][]node
	stopSignal chan struct{}
}

// newMachine creates a new machine from the given machine config . It
// creates empty nodes based on the configuration and wires them up to each
// other.
func newMachine(config machineConfig) (machine, error) {
	var m machine

	m.stopSignal = make(chan struct{})

	// Construct an empty array of nodes based on the size of the nodes in the config
	m.nodes = make([][]node, len(config.Nodes))
	for i := range m.nodes {
		m.nodes[i] = make([]node, len(config.Nodes[0]))
	}

	for x, valX := range config.Nodes {
		for y, valY := range valX {
			switch valY {
			case "e":
				// The node is an execution node

				// The node's ports default to ones that go nowhere
				var up, down, left, right numberReadWriter = newPort(), newPort(), newPort(), newPort()
				any := newAnyPort(up.(*port), down.(*port), left.(*port), right.(*port))

				if y-1 >= 0 {
					// If there is a node above this node, connect it
					up = m.nodes[x][y-1].getDown()
				}
				if x-1 >= 0 {
					// If there is a node to the left of this node, connect it
					left = m.nodes[x-1][y].getRight()
				}

				m.nodes[x][y] = newExecutionNode(up, down, left, right, any, any.lastUsedPort)
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
