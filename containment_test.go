package main

import (
	"testing"
)

type FakeExecuter struct {
	host    Host
	command string
}

func (f *FakeExecuter) Execute(host Host, command string) ([]byte, error) {
	f.host = host
	f.command = command
	return []byte("Fakely Executed"), nil
}

func (f *FakeExecuter) Validate(t *testing.T, expectedAddress string, expectedPort int, expectedUser string, expectedCommand string) {
	if f.host.Address != expectedAddress {
		t.Errorf("Address should have been %q but was %q", expectedAddress, f.host.Address)
	}
	if f.host.Port != expectedPort {
		t.Errorf("Port should have been %q but was %q", expectedPort, f.host.Port)
	}
	if f.host.User != expectedUser {
		t.Errorf("User should have been %q but was %q", expectedUser, f.host.User)
	}
	if f.command != expectedCommand {
		t.Errorf("Command should have been \n%q\nbut was\n%q", expectedCommand, f.command)
	}
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

	fakeExecuter.Validate(t, "1.1.1.2", 45, "derp", "sudo docker pull ubuntu")
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

	fakeExecuter.Validate(t, "1.1.1.2", 45, "derp", "sudo docker run -d --name something-something -p 80:80 -p 123:123 something/something")
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

	fakeExecuter.Validate(t, "1.1.1.2", 45, "derp", "sudo docker stop something-something && sudo docker rm something-something")
}
