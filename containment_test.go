package main

import (
	"io"
	"testing"
)

type FakeExecuter struct {
	executedTask bool
	host         Host
	command      string
}

func (f *FakeExecuter) Execute(host Host, command string, writer io.Writer) error {
	f.executedTask = true
	f.host = host
	f.command = command
	writer.Write([]byte("Fakely Executed"))
	return nil
}

func (f *FakeExecuter) Validate(t *testing.T, expectedAddress string, expectedPort int, expectedUser string, expectedCommand string) {
	if f.executedTask == false {
		t.Fatal("Expected a task to be executed")
	}
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

var oldExecuter Executer
var fakeExecuter *FakeExecuter

func enableFakeExecuter() {
	oldExecuter = executer
	fakeExecuter = &FakeExecuter{}
	executer = fakeExecuter
}

func disablefakeExecuter() {
	executer = oldExecuter
}

func simpleConfiguration() Configuration {
	return Configuration{
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
}

func TestPullsContainers(t *testing.T) {
	enableFakeExecuter()
	defer disablefakeExecuter()

	configuration := simpleConfiguration()

	_, err := captureStdout(func() error {
		return update(configuration, "something/something")
	})
	if err != nil {
		t.Error(err)
	}

	fakeExecuter.Validate(t, "1.1.1.2", 45, "derp", "sudo docker pull something/something")
}

func TestStartsContainers(t *testing.T) {
	enableFakeExecuter()
	defer disablefakeExecuter()

	configuration := simpleConfiguration()

	_, err := captureStdout(func() error {
		return start(configuration, "something/something")
	})
	if err != nil {
		t.Error(err)
	}

	fakeExecuter.Validate(t, "1.1.1.2", 45, "derp", "sudo docker run -d --name something-something -p 80:80 -p 123:123 something/something")
}

func TestStopsContainers(t *testing.T) {
	enableFakeExecuter()
	defer disablefakeExecuter()

	configuration := simpleConfiguration()

	_, err := captureStdout(func() error {
		return stop(configuration, "something/something")
	})
	if err != nil {
		t.Error(err)
	}

	fakeExecuter.Validate(t, "1.1.1.2", 45, "derp", "sudo docker stop something-something && sudo docker rm something-something")
}

func TestRestartsContainers(t *testing.T) {
	enableFakeExecuter()
	defer disablefakeExecuter()

	configuration := simpleConfiguration()

	_, err := captureStdout(func() error {
		return restart(configuration, "something/something")
	})
	if err != nil {
		t.Error(err)
	}

	fakeExecuter.Validate(
		t,
		"1.1.1.2",
		45,
		"derp",
		"sudo docker stop something-something && sudo docker rm something-something && sudo docker run -d --name something-something -p 80:80 -p 123:123 something/something",
	)
}

func TestListsContainerStatus(t *testing.T) {
	enableFakeExecuter()
	defer disablefakeExecuter()

	configuration := simpleConfiguration()

	output, err := captureStdout(func() error {
		return status(configuration, "something/something")
	})
	if err != nil {
		t.Error(err)
	}

	fakeExecuter.Validate(t, "1.1.1.2", 45, "derp", "sudo docker inspect -f '{{.State.Running}}' something-something")

	if expected := "[derp@1.1.1.2] something/something stopped\n"; output != expected {
		t.Errorf("Expected output to be\n%q\nbut was\n%q", expected, output)
	}
}
