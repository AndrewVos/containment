package main

import (
	"io"
)

type Executer interface {
	Execute(host Host, command string, writer io.Writer) error
}
