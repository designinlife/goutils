package goutils

import (
	"testing"
)

func TestSSHClient_Upload(t *testing.T) {
	client := NewSSHClient("127.0.0.1", 22, "root", "~/.ssh/id_rsa", false, SSHOptionWithProxy("http://127.0.0.1:3128"))
	err := client.Upload("d:/tmp/php-8.0.7.tar.gz", "/root/php-8.0.7.tar.gz")

	if err != nil {
		t.Errorf("%v", err)
	}
}
