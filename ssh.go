package goutils

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-isatty"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	logger "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
)

type SSHTunnelEndpoint struct {
	Host string
	Port int
}

func (endpoint *SSHTunnelEndpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

type SSHTunnel struct {
	Local  *SSHTunnelEndpoint
	Server *SSHTunnelEndpoint
	Remote *SSHTunnelEndpoint

	Config *ssh.ClientConfig
}

func (tunnel *SSHTunnel) Start() error {
	listener, err := net.Listen("tcp", tunnel.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	logger.Debugf("[SSHTunnel] Listen: %s", tunnel.Local.String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go tunnel.forward(conn)
	}
}

func (tunnel *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", tunnel.Server.String(), tunnel.Config)
	if err != nil {
		logger.Errorf("[SSHTunnel] Server dial error: %s", err)
		return
	}

	remoteConn, err := serverConn.Dial("tcp", tunnel.Remote.String())
	if err != nil {
		logger.Errorf("[SSHTunnel] Remote dial error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			logger.Errorf("[SSHTunnel] io.Copy error: %s", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

// SSHClient SSH 客户端
type SSHClient struct {
	// RSA 私钥证书路径或全文内容
	PrivateKey string
	// 主机 Domain/IP 地址
	Host string
	// SSH 端口
	Port int
	// SSH 登录用户名
	User string
	// 静默方式: 不输出 Stdout 信息
	Quiet bool
	// 是否已连接？
	Connected bool
	// SSH 客户端实例
	Client *ssh.Client
	// 代理服务器地址 (支持 http,https,socks5,socks5h, 例如: http://127.0.0.1:3128, socks5://127.0.0.1:1080)
	Proxy string
	// 超时时间 (默认不超时)
	Timeout time.Duration
	// 开启 TTY 终端模式
	TTY bool
	// 缓冲区大小 (默认: 8192 bytes)
	ChunkSize uint16
	// SSH 隧道
	Tunnel *SSHTunnel
}

type SSHClientOption func(*SSHClient)

// SSHOptionWithProxy 支持 http/https/socks5/socks5h 代理协议。
func SSHOptionWithProxy(proxyUrl string) SSHClientOption {
	return func(c *SSHClient) {
		c.Proxy = proxyUrl
	}
}

func SSHOptionWithTimeout(timeout time.Duration) SSHClientOption {
	return func(c *SSHClient) {
		c.Timeout = timeout
	}
}

func SSHOptionWithChunkSize(chunkSize uint16) SSHClientOption {
	return func(c *SSHClient) {
		c.ChunkSize = chunkSize
	}
}

func SSHOptionWithTunnel(tunnel *SSHTunnel) SSHClientOption {
	return func(c *SSHClient) {
		c.Tunnel = tunnel
	}
}

func NewSSHClient(host string, port int, user string, privateKey string, quiet bool, opts ...SSHClientOption) *SSHClient {
	c := &SSHClient{
		Host:       host,
		Port:       port,
		User:       user,
		PrivateKey: privateKey,
		Quiet:      quiet,
		ChunkSize:  8192,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func newSSHClientWithProxy(proxyAddress, sshServerAddress string, sshConfig *ssh.ClientConfig) (*ssh.Client, error) {
	// dialer, err := proxy.SOCKS5("tcp", proxyAddress, nil, proxy.Direct)
	proxyUrl, err := url.Parse(proxyAddress)

	if err != nil {
		return nil, err
	}

	dialer, err := proxy.FromURL(proxyUrl, proxy.Direct)

	if err != nil {
		return nil, err
	}

	conn, err := dialer.Dial("tcp", sshServerAddress)
	if err != nil {
		return nil, err
	}

	c, chans, reqs, err := ssh.NewClientConn(conn, sshServerAddress, sshConfig)
	if err != nil {
		return nil, err
	}

	return ssh.NewClient(c, chans, reqs), nil
}

func (s *SSHClient) Connect() error {
	if !s.Connected {
		// var hostKey ssh.PublicKey
		var key []byte

		if strings.HasPrefix(s.PrivateKey, "~/") || strings.HasPrefix(s.PrivateKey, "/") {
			pkey, err := homedir.Expand(s.PrivateKey)

			if err != nil {
				return errors.Wrapf(err, "Load private key path error.")
			}

			key, err = ioutil.ReadFile(pkey)

			if err != nil {
				return errors.Wrapf(err, "Unable to read private key: %s", s.PrivateKey)
			}
		} else {
			key = []byte(s.PrivateKey)
		}

		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return errors.Wrapf(err, "Ubable to parse private key %s", s.PrivateKey)
		}

		config := &ssh.ClientConfig{
			User: s.User,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			// HostKeyCallback: ssh.FixedHostKey(hostKey),
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         s.Timeout,
		}

		var client *ssh.Client

		if s.Proxy != "" {
			client, err = newSSHClientWithProxy(s.Proxy, fmt.Sprintf("%s:%d", s.Host, s.Port), config)
		} else {
			client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port), config)
		}

		if err != nil {
			return errors.Wrapf(err, "Unable to connect %s:%d.", s.Host, s.Port)
		}

		s.Client = client
		s.Connected = true
	}

	return nil
}

func (s *SSHClient) Close() error {
	if s.Connected {
		return s.Client.Close()
	}

	return nil
}

func (s *SSHClient) Run(command string) (int, error) {
	return s.RunWithWriter(command, nil)
}

func (s *SSHClient) RunWithWriter(command string, w io.Writer) (int, error) {
	err := s.Connect()
	if err != nil {
		return -1, err
	}

	session, err := s.Client.NewSession()
	if err != nil {
		return -2, err
	}

	if s.TTY {
		// Set up terminal modes
		modes := ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}

		err = session.RequestPty("xterm", 24, 80, modes)

		if err != nil {
			return -2, errors.Wrapf(err, "Failed to set tty. (%s:%d)", s.Host, s.Port)
		}
	}

	defer session.Close()

	stderr, _ := session.StderrPipe()
	stdout, _ := session.StdoutPipe()

	if err := session.Start(command); err != nil {
		return -3, err
	}

	var scanner *bufio.Scanner

	if w != nil {
		scanner = bufio.NewScanner(io.TeeReader(io.MultiReader(stdout, stderr), w))
	} else {
		scanner = bufio.NewScanner(io.MultiReader(stdout, stderr))
	}

	for scanner.Scan() {
		m := scanner.Text()

		if !s.Quiet {
			logger.Info(m)
		}
	}

	if err := session.Wait(); err != nil {
		if exiterr, ok := err.(*ssh.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			exitstatus := exiterr.ExitStatus()

			return exitstatus, err
		} else {
			return -4, err
		}
	}

	return 0, nil
}

func (s *SSHClient) Upload(src, dst string) error {
	err := s.Connect()
	if err != nil {
		return err
	}

	sftpClient, err := sftp.NewClient(s.Client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	dstFile, err := sftpClient.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	fileInfo, err := srcFile.Stat()

	if err != nil {
		return err
	}

	totalByteCount := fileInfo.Size()
	readByteCount := 0

	buf := make([]byte, s.ChunkSize)

	isTty := isatty.IsTerminal(os.Stdout.Fd())

	for {
		n, err := srcFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				readByteCount = readByteCount + n

				_, _ = dstFile.Write(buf[:n])

				if isTty && !s.Quiet {
					fmt.Printf("\r%.2f%%", float32(readByteCount)*100/float32(totalByteCount))
				}
				// logger.Debugf("readByteCount=%d, n=%d, %v", readByteCount, n, err)
				break
			}
		}

		readByteCount = readByteCount + n

		_, _ = dstFile.Write(buf[:n])

		if isTty && !s.Quiet {
			fmt.Printf("\r%.2f%%", float32(readByteCount)*100/float32(totalByteCount))
		}
	}

	if isTty && !s.Quiet {
		logger.Infof("Uploaded. (%s -> %s)", src, dst)
	}

	// s.Run(fmt.Sprintf("ls -lh %s", dst))

	return nil
}

func (s *SSHClient) Download(src, dst string) error {
	err := s.Connect()
	if err != nil {
		return err
	}

	sftpClient, err := sftp.NewClient(s.Client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := sftpClient.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcFileInfo, err := srcFile.Stat()

	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	totalByteCount := srcFileInfo.Size()
	readByteCount := 0

	fmt.Println(totalByteCount)

	buf := make([]byte, s.ChunkSize)

	isTty := isatty.IsTerminal(os.Stdout.Fd())

	for {
		n, err := srcFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				readByteCount = readByteCount + n

				_, _ = dstFile.Write(buf[:n])

				if isTty && !s.Quiet {
					fmt.Printf("\r%.2f%%", float32(readByteCount)*100/float32(totalByteCount))
				}
				// logger.Debugf("readByteCount=%d, n=%d, %v", readByteCount, n, err)
				break
			}
		}

		readByteCount = readByteCount + n

		_, _ = dstFile.Write(buf[:n])

		if isTty && !s.Quiet {
			fmt.Printf("\r%.2f%%", float32(readByteCount)*100/float32(totalByteCount))
		}
	}

	// if _, err := srcFile.WriteTo(dstFile); err != nil {
	// 	return err
	// }

	if isTty && !s.Quiet {
		logger.Infof("Downloaded. (%s -> %s)", src, dst)
	}

	return nil
}
