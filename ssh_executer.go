package main

import (
	"code.google.com/p/go.crypto/ssh"
	"code.google.com/p/go.crypto/ssh/agent"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
)

type SSHExecuter struct {
	mutex   sync.Mutex
	auths   []ssh.AuthMethod
	clients map[string]*ssh.Client
}

func (s *SSHExecuter) loadAuths() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.auths == nil {
		sock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
		if err == nil {
			agent := agent.NewClient(sock)
			signers, err := agent.Signers()
			if err == nil {
				s.auths = []ssh.AuthMethod{ssh.PublicKeys(signers...)}
			}
		}
	}
}

func (s *SSHExecuter) Execute(host Host, command string, out io.Writer) error {
	s.loadAuths()

	s.mutex.Lock()
	if s.clients == nil {
		s.clients = map[string]*ssh.Client{}
	}
	s.mutex.Unlock()

	if _, ok := s.clients[host.Identifier()]; !ok {
		clientConfig := &ssh.ClientConfig{
			User: host.User,
			Auth: s.auths,
		}
		clientConfig.SetDefaults()

		port := 22
		if host.Port > 0 {
			port = host.Port
		}
		addressAndPort := fmt.Sprintf("%v:%d", host.Address, port)

		client, err := ssh.Dial("tcp", addressAndPort, clientConfig)
		if err != nil {
			return err
		}
		s.mutex.Lock()
		s.clients[host.Identifier()] = client
		s.mutex.Unlock()
	}
	client := s.clients[host.Identifier()]

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = out
	session.Stderr = out

	err = session.Run(command)
	return err

}
