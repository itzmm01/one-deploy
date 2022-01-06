package ftpclient

import (
	"fmt"
	"os"
	"time"

	"github.com/jlaffaye/ftp"
)

// FtpClient-INFO
type FtpClient struct {
	Host     string
	Port     string
	Username string
	Password string
}

// upload
func (ctx FtpClient) Upload(srcFile, dstPath, fileName string) error {
	ftpServer := fmt.Sprintf("%v:%v", ctx.Host, ctx.Port)
	c, err := ftp.Dial(ftpServer, ftp.DialWithTimeout(5*time.Second))
	defer c.Quit()
	if err != nil {
		return err
	}
	err = c.Login(ctx.Username, ctx.Password)
	if err != nil {
		return err
	}
	_ = c.MakeDir(dstPath)
	err = c.ChangeDir(dstPath)
	if err != nil {
		return err
	}
	file, _ := os.Open(srcFile)
	defer file.Close()
	err = c.Stor(fileName, file)
	if err != nil {
		return err
	}

	c.Logout()
	return nil
}
