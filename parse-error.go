package main

import (
	"errors"
	"fmt"
)

func newParseError(message string, c char) error {
	return errors.New(fmt.Sprint(message, " at line ", c.line, ", character ", c.pos))
}
