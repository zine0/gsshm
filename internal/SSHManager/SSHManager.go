package sshmanager

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"

	sshclient "github.com/zine0/gsshm/internal/SSHClient"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

type SSHManager struct {
	connection *sshclient.SSHClient
}

func NewSSHManager() *SSHManager {
	return &SSHManager{}
}

func (m *SSHManager) Connect(host, port, username, password, keypath string) error {
	hostKeyCallback, err := knownhosts.New(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	if err != nil {
		return fmt.Errorf("failed to create host key callback: %w", err)
	}

	authMethods := []ssh.AuthMethod{}

	if keypath != "" {
		key, err := ioutil.ReadFile(keypath)
		if err != nil {
			return fmt.Errorf("failed to read private key file: %w", err)
		}

		var signer ssh.Signer
		if password != "" {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(password))
		} else {
			signer, err = ssh.ParsePrivateKey(key)
		}
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	if len(authMethods) == 0 {
		return fmt.Errorf("no authentication methods provided")
	}

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            authMethods,
		HostKeyCallback: hostKeyCallback,
	}

	address := net.JoinHostPort(host, port)
	conn, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	client := &sshclient.SSHClient{
		Host:     host,
		Port:     port,
		Username: username,
		Conn:     conn,
	}
	client.SetIO(os.Stdin, os.Stdout, os.Stderr)

	m.connection = client
	return nil
}

func (m *SSHManager) Close() error {
	if m.connection == nil {
		return nil
	}
	return m.connection.Close()
}

func (m *SSHManager) StartTerminal() error {
	if m.connection == nil {
		return fmt.Errorf("not connected")
	}

	if err := m.connection.NewTerminalSession(); err != nil {
		return err
	}

	return m.connection.WaitSession()
}
