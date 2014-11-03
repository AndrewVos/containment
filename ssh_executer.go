package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
)

type SSHExecuter struct{}

func (s *SSHExecuter) Execute(host Host, command string) ([]byte, error) {
	arguments := []string{
		"-T",
		"-o", "StrictHostKeyChecking=no",
	}
	if host.Port > 0 {
		arguments = append(arguments, "-p")
		arguments = append(arguments, strconv.Itoa(host.Port))
	}
	arguments = append(arguments, fmt.Sprintf("%v@%v", host.User, host.Address))
	arguments = append(arguments, command)

	cmd := exec.Command("ssh", arguments...)

	buffer := &bytes.Buffer{}
	cmd.Stdout = buffer
	cmd.Stderr = buffer

	err := cmd.Start()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("I couldn't launch ssh:\n%v", err))
	}

	err = cmd.Wait()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("I had some problems running %q on %q:\nError:\n%v\nOutput:\n%v", command, host.Address, err, string(buffer.Bytes())))
	}

	return buffer.Bytes(), nil
}
