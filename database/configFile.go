package database

import (
	"fmt"
	"one-backup/cmd"
	"one-backup/ssh"
	"strconv"
	"strings"
)

// file info
type File struct {
	// 压缩包文件名
	TarFilename string
	// 保存目录
	SaveDir string
	// 备份目录
	BackupDir string
	// name
	Name string
	// 主机
	Host string
	// 端口
	Port string
	// 账号
	Username string
	// 密码
	Password string
	KeyFile  string
	// 路径
	Path string
}

// backup
func (ctx File) Backup() error {
	if ctx.Host == "local" {
		destName := strings.Split(ctx.Path, `/`)
		cmd_str := fmt.Sprintf("/bin/cp -rf %v %v/%v", ctx.Path, ctx.BackupDir, destName[len(destName)-1])
		return cmd.Run(cmd_str, Debug)
	} else {
		cliConf := new(ssh.ClientConfig)
		sshPort, _ := strconv.ParseInt(ctx.Port, 10, 64)
		cliConf.CreateClient(ctx.Host, sshPort, ctx.Username, ctx.Password, ctx.KeyFile)
		pathStrList := strings.Split(ctx.Path, `/`)
		fileName := pathStrList[len(pathStrList)-1]
		if fileName == "" {
			return cliConf.DownloadDirectory(ctx.Path, ctx.BackupDir+"/"+fileName)
		} else {

			return cliConf.Download(ctx.Path, ctx.BackupDir+"/"+fileName)
		}

	}

}
