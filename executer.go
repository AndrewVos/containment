package main

type Executer interface {
	Execute(host Host, command string) ([]byte, error)
}
