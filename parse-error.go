package main

import (
	"errors"
	"fmt"
)

func newParseError(message string, c char) error {
	return errors.New(fmt.Sprint(string(message), "at line", string(c.line)+", character", c.pos))
}
