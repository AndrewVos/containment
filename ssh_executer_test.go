package main

import (
	"strings"
	"testing"
)

func TestExecutesCommand(t *testing.T) {
	container, err := NewDockerContainer()
	if err != nil {
		t.Fatalf("Couldn't create docker container cause:\n%v\n", err)
	}
	defer container.Kill()

	executer := &SSHExecuter{}
	b, _ := executer.Execute(Host{Address: container.ip, Port: container.port, User: container.user}, "echo hello!")
	output := string(b)
	if strings.Contains(output, "hello!") == false {
		t.Errorf("Expected output to contain some echoed text but was this:\n%v\n", output)
	}
}
