package main

type Executer interface {
	Execute(address string, port int, user string, command string) ([]byte, error)
}
