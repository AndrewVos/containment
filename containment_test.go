package main

import (
	"testing"
)

type FakeExecuter struct {
	address string
	port    int
	command string
	user    string
}

func (s *FakeExecuter) Execute(address string, port int, user string, command string) ([]byte, error) {
	s.address = address
	s.port = port
	s.user = user
	s.command = command
	return nil, nil
}

func TestPullsContainers(t *testing.T) {
	oldExecuter := executer
	defer func() { executer = oldExecuter }()
	fakeExecuter := &FakeExecuter{}
	executer = fakeExecuter

	configuration := Configuration{
		Clusters: []Cluster{
			Cluster{
				Name: "some-cluster",
				Hosts: []Host{
					Host{Address: "1.1.1.2", Port: 45, User: "derp"},
				},
			},
		},
		Containers: []Container{
			Container{Image: "ubuntu", Clusters: []string{"some-cluster"}},
		},
	}

	update(configuration, "ubuntu")

	expectedAddress := "1.1.1.2"
	expectedPort := 45
	expectedUser := "derp"
	expectedCommand := "docker pull ubuntu"

	if fakeExecuter.address != expectedAddress {
		t.Errorf("Address should have been %v but was %v", expectedAddress, fakeExecuter.address)
	}
	if fakeExecuter.port != expectedPort {
		t.Errorf("Port should have been %v but was %v", expectedPort, fakeExecuter.port)
	}
	if fakeExecuter.user != expectedUser {
		t.Errorf("User should have been %v but was %v", expectedUser, fakeExecuter.user)
	}
	if fakeExecuter.command != expectedCommand {
		t.Errorf("Command should have been %v but was %v", expectedCommand, fakeExecuter.command)
	}
}
