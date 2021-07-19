package goutils

import (
	"fmt"
	"testing"
)

func TestSSHClient_Upload(t *testing.T) {
	client := NewSSHClient("127.0.0.1", 22, "root", "~/.ssh/id_rsa", false, SSHOptionWithProxy("http://127.0.0.1:3128"))
	err := client.Upload("d:/tmp/php-8.0.7.tar.gz", "/root/php-8.0.7.tar.gz")

	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestSSHClient_Tunnel(t *testing.T) {
	client := NewSSHClient("127.0.0.1", 0, "root", "~/.ssh/lilei", false, SSHOptionWithTunnel(&SSHTunnel{
		Local:  &SSHTunnelEndpoint{Host: "127.0.0.1", Port: 0},
		Remote: &SSHTunnelEndpoint{Host: "10.8.0.2", Port: 22},
		Server: &SSHTunnelEndpoint{Host: "192.168.1.1", Port: 22},
	}))

	exitcode, err := client.Run("whoami && pwd && hostname")

	if err != nil {
		t.Errorf("%v", err)
	}

	fmt.Println(exitcode)
}
