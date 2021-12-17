package ssh

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/pkg/sftp"
	"github.com/wonderivan/logger"
	"golang.org/x/crypto/ssh"
)

//连接的配置
type ClientConfig struct {
	Host       string       //ip
	Port       int64        // 端口
	Username   string       //用户名
	Password   string       //密码
	sshClient  *ssh.Client  //ssh client
	sftpClient *sftp.Client //sftp client
	LastResult string       //最近一次运行的结果
}

// CreateClient
func (cliConf *ClientConfig) CreateClient(host string, port int64, username, password string) {
	var (
		sshClient  *ssh.Client
		sftpClient *sftp.Client
		err        error
	)
	cliConf.Host = host
	cliConf.Port = port
	cliConf.Username = username
	cliConf.Password = password
	cliConf.Port = port

	config := ssh.ClientConfig{
		User: cliConf.Username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", cliConf.Host, cliConf.Port)

	if sshClient, err = ssh.Dial("tcp", addr, &config); err != nil {
		logger.Error("error occurred:", err)
	}
	cliConf.sshClient = sshClient

	//此时获取了sshClient，下面使用sshClient构建sftpClient
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		logger.Error("error occurred:", err)
	}
	cliConf.sftpClient = sftpClient
}

// RunShell
func (cliConf *ClientConfig) RunShell(shell string) (res string, error1 error) {
	var (
		session *ssh.Session
		err     error
	)

	//获取session，这个session是用来远程执行操作的
	if session, err = cliConf.sshClient.NewSession(); err != nil {
		return "", err
	}
	//执行shell
	if output, err := session.CombinedOutput(shell); err != nil {
		return "", err
	} else {
		cliConf.LastResult = string(output)
	}
	return cliConf.LastResult, nil
}

// Upload
func (cliConf *ClientConfig) Upload(srcPath, dstPath string) error {
	srcFile, _ := os.Open(srcPath)                   //本地
	dstFile, _ := cliConf.sftpClient.Create(dstPath) //远程
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()
	buf := make([]byte, 1024)
	for {
		n, err := srcFile.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			} else {
				break
			}
		}
		_, _ = dstFile.Write(buf[:n])
	}
	return nil
}

// Download
func (cliConf *ClientConfig) Download(srcPath, dstPath string) error {
	srcFile, _ := cliConf.sftpClient.Open(srcPath) //远程
	dstFile, _ := os.Create(dstPath)               //本地
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		return err
	}
	return nil
}
