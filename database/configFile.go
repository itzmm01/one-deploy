package database

import (
	"fmt"
	"one-backup/cmd"
	"one-backup/ssh"
	"strconv"
	"strings"
)

type File struct {
	TarFilename string
	SaveDir     string
	BackupDir   string
	Name        string
	Host        string
	Port        string
	Username    string
	Password    string
	Path        string
}

func (ctx File) Backup() error {
	if ctx.Host == "local" {
		destName := strings.Split(ctx.Path, `/`)
		cmd_str := fmt.Sprintf("/bin/cp -rf %v %v/%v", ctx.Path, ctx.BackupDir, destName[len(destName)-1])
		return cmd.Run(cmd_str)
	} else {
		cliConf := new(ssh.ClientConfig)
		sshPort, _ := strconv.ParseInt(ctx.Port, 10, 64)
		cliConf.CreateClient(ctx.Host, sshPort, ctx.Username, ctx.Password)
		pathStrList := strings.Split(ctx.Path, `/`)
		fileName := pathStrList[len(pathStrList)-1]

		return cliConf.Download(ctx.Path, ctx.BackupDir+"/"+fileName)
	}

}
