package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

var currentContainerPort = 45000

type dockerContainer struct {
	id   string
	user string
	ip   string
	port int
}

func init() {
	buildTestDockerContainer()
}

func NewDockerContainer() (*dockerContainer, error) {
	cmd := exec.Command("docker", "run", "-d", "containment/test")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		return nil, err
	}
	id := strings.Replace(string(out), "\n", "", -1)

	cmd = exec.Command("docker", "inspect", "--format", "{{ .NetworkSettings.IPAddress }}", id)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	ip := strings.Replace(string(out), "\n", "", -1)
	return &dockerContainer{id: id, user: "root", ip: ip, port: 22}, nil
}

func (c *dockerContainer) Kill() {
	cmd := exec.Command("docker", "kill", c.id)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}
}

func buildTestDockerContainer() {
	dockerFile := `
FROM ubuntu

RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list
RUN apt-get update

RUN apt-get install -y openssh-server
RUN mkdir /var/run/sshd

RUN apt-get install -y docker

RUN passwd -d root
RUN echo "PermitEmptyPasswords yes" > /etc/ssh/sshd_config

CMD /usr/sbin/sshd -D`

	file, err := os.Create("Dockerfile")
	defer os.Remove("Dockerfile")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write([]byte(dockerFile))
	cmd := exec.Command("docker", "build", "-t", "containment/test", ".")
	buffer := &bytes.Buffer{}
	cmd.Stdout = buffer
	cmd.Stderr = buffer
	err = cmd.Run()
	if err != nil {
		log.Fatalln(err, buffer.Bytes())
	}
}

func captureStdout(f func() error) (string, error) {
	tempFile, _ := ioutil.TempFile("", "stdout")
	oldStdout := os.Stdout
	os.Stdout = tempFile
	err := f()
	os.Stdout = oldStdout
	tempFile.Close()
	b, _ := ioutil.ReadFile(tempFile.Name())
	return string(b), err
}
