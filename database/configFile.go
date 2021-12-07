package database

import (
	"fmt"
	"one-backup/cmd"
	"strings"

	"github.com/wonderivan/logger"
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
	cmd_str := ""
	if ctx.Host == "local" {
		destName := strings.Split(ctx.Path, `/`)
		cmd_str = fmt.Sprintf("/bin/cp -rf %v %v/%v", ctx.Path, ctx.BackupDir, destName[len(destName)-1])

	} else {
		logger.Info("no support ssh")
		return nil
	}

	return cmd.Run(cmd_str)
}
