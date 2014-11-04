package main

import (
	"code.google.com/p/go.crypto/ssh"
	"code.google.com/p/go.crypto/ssh/agent"
	"fmt"
	"net"
	"os"
)

type SSHExecuter struct{}

func (s *SSHExecuter) Execute(host Host, command string) ([]byte, error) {
	var auths []ssh.AuthMethod

	sock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err == nil {
		agent := agent.NewClient(sock)
		signers, err := agent.Signers()
		if err == nil {
			auths = []ssh.AuthMethod{ssh.PublicKeys(signers...)}
		}
	}

	clientConfig := &ssh.ClientConfig{
		User: host.User,
		Auth: auths,
	}
	clientConfig.SetDefaults()

	port := 22
	if host.Port > 0 {
		port = host.Port
	}
	addressAndPort := fmt.Sprintf("%v:%d", host.Address, port)

	client, err := ssh.Dial("tcp", addressAndPort, clientConfig)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	if err != nil {
		return output, err
	}

	return output, nil
}
