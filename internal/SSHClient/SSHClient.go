package sshclient

import (
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	Host     string
	Port     string
	Username string
	Conn     *ssh.Client
	session  *ssh.Session
	stdin    io.Reader
	stdout   io.Writer
	stderr   io.Writer
}

func (c *SSHClient) NewTerminalSession() error {
	session, err := c.Conn.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return fmt.Errorf("failed to request PTY: %w", err)
	}

	session.Stdin = c.stdin
	session.Stdout = c.stdout
	session.Stderr = c.stderr

	if err := session.Shell(); err != nil {
		session.Close()
		return fmt.Errorf("failed to start shell: %w", err)
	}

	c.session = session
	return nil
}

func (c *SSHClient) WaitSession() error {
	if c.session == nil {
		return fmt.Errorf("no active session")
	}
	return c.session.Wait()
}

func (c *SSHClient) Close() error {
	var err error
	if c.session != nil {
		err = c.session.Close()
	}
	if c.Conn != nil {
		connErr := c.Conn.Close()
		if connErr != nil && err == nil {
			err = connErr
		}
	}
	return err
}

// Getters for private fields
func (c *SSHClient) SetIO(stdin io.Reader, stdout, stderr io.Writer) {
	c.stdin = stdin
	c.stdout = stdout
	c.stderr = stderr
}
