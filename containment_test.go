package main

import (
	"fmt"
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
	return []byte("Fakely Executed"), nil
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

	err := update(configuration, "ubuntu")
	if err != nil {
		t.Error(err)
	}

	expectedAddress := "1.1.1.2"
	expectedPort := 45
	expectedUser := "derp"
	expectedCommand := "sudo docker pull ubuntu"

	if fakeExecuter.address != expectedAddress {
		t.Errorf("Address should have been %q but was %q", expectedAddress, fakeExecuter.address)
	}
	if fakeExecuter.port != expectedPort {
		t.Errorf("Port should have been %q but was %q", expectedPort, fakeExecuter.port)
	}
	if fakeExecuter.user != expectedUser {
		t.Errorf("User should have been %q but was %q", expectedUser, fakeExecuter.user)
	}
	if fakeExecuter.command != expectedCommand {
		t.Errorf("Command should have been \n%q\nbut was\n%q", expectedCommand, fakeExecuter.command)
	}
}

func TestStartsContainers(t *testing.T) {
	oldExecuter := executer
	defer func() { executer = oldExecuter }()
	fakeExecuter := &FakeExecuter{}
	executer = fakeExecuter

	configuration := Configuration{
		Clusters: []Cluster{
			Cluster{
				Name: "some-cluster", Hosts: []Host{Host{Address: "1.1.1.2", Port: 45, User: "derp"}},
			},
		},
		Containers: []Container{
			Container{
				Image:    "something/something",
				Clusters: []string{"some-cluster"},
				Ports:    []string{"80:80", "123:123"},
			},
		},
	}

	err := start(configuration, "something/something")
	if err != nil {
		t.Error(err)
	}

	expectedAddress := "1.1.1.2"
	expectedPort := 45
	expectedUser := "derp"
	expectedCommand := "sudo docker run -d --name something-something -p 80:80 -p 123:123 something/something"

	fmt.Println(fakeExecuter)
	if fakeExecuter.address != expectedAddress {
		t.Errorf("Address should have been %q but was %q", expectedAddress, fakeExecuter.address)
	}
	if fakeExecuter.port != expectedPort {
		t.Errorf("Port should have been %d but was %d", expectedPort, fakeExecuter.port)
	}
	if fakeExecuter.user != expectedUser {
		t.Errorf("User should have been %q but was %q", expectedUser, fakeExecuter.user)
	}
	if fakeExecuter.command != expectedCommand {
		t.Errorf("Command should have been \n%q\nbut was\n%q", expectedCommand, fakeExecuter.command)
	}
}

func TestStopsContainers(t *testing.T) {
	oldExecuter := executer
	defer func() { executer = oldExecuter }()
	fakeExecuter := &FakeExecuter{}
	executer = fakeExecuter

	configuration := Configuration{
		Clusters: []Cluster{
			Cluster{
				Name: "some-cluster", Hosts: []Host{Host{Address: "1.1.1.2", Port: 45, User: "derp"}},
			},
		},
		Containers: []Container{
			Container{
				Image:    "something/something",
				Clusters: []string{"some-cluster"},
				Ports:    []string{"80:80", "123:123"},
			},
		},
	}

	err := stop(configuration, "something/something")
	if err != nil {
		t.Error(err)
	}

	expectedAddress := "1.1.1.2"
	expectedPort := 45
	expectedUser := "derp"
	expectedCommand := "sudo docker stop something-something && sudo docker rm something-something"

	if fakeExecuter.address != expectedAddress {
		t.Errorf("Address should have been %q but was %q", expectedAddress, fakeExecuter.address)
	}
	if fakeExecuter.port != expectedPort {
		t.Errorf("Port should have been %d but was %d", expectedPort, fakeExecuter.port)
	}
	if fakeExecuter.user != expectedUser {
		t.Errorf("User should have been %q but was %q", expectedUser, fakeExecuter.user)
	}
	if fakeExecuter.command != expectedCommand {
		t.Errorf("Command should have been \n%q\nbut was\n%q", expectedCommand, fakeExecuter.command)
	}
}
