package database

import (
	"fmt"
	"one-backup/cmd"
	"one-backup/keygen"
	"one-backup/tool"
)

type Mysql struct {
	/*
		database: test1,yc
		host: 192.168.146.134
		name: mysql
		password: BhaBUTSg3lMXHLVUkHmOfw==
		port: "3306"
		type: mysql
		username: root
	*/
	TarFilename string
	SaveDir     string
	BackupDir   string
	Name        string
	Host        string
	Port        string
	Username    string
	Password    string
	Database    string
}

func (ctx *Mysql) Backup() error {
	cmdStr := fmt.Sprintf(
		"mysqldump -h %v -P %v -u%v -p'%v' ", ctx.Host, ctx.Port, ctx.Username, ctx.Password,
	)
	if ctx.Database == "alldatabase" {
		cmdStr = cmdStr + "--all-databases "
	} else {
		cmdStr = cmdStr + "--databases " + ctx.Database
	}

	cmdStr = cmdStr + fmt.Sprintf(" > %v/%v.sql", ctx.BackupDir, ctx.Database)

	err := cmd.Run(cmdStr, Debug)
	if err == nil {
		keygen.AesEncryptCBCFile(fmt.Sprintf("%v/%v.sql", ctx.BackupDir, ctx.Database), fmt.Sprintf("%v/%v-Encrypt.sql", ctx.BackupDir, ctx.Database))
		return cmd.Run(fmt.Sprintf("rm -f %v/%v.sql", ctx.BackupDir, ctx.Database), Debug)
	} else {
		return err
	}

}

func (ctx Mysql) Restore(filepath string) error {
	dstPath := "/tmp/" + tool.RandomString(20)
	keygen.AesDecryptCBCFile(filepath, dstPath)
	cmd_str := fmt.Sprintf("cat '%v' | mysql -h %v -P %v -u%v -p%v  ; rm -f %v", dstPath, ctx.Host, ctx.Port, ctx.Username, ctx.Password, dstPath)
	return cmd.Run(cmd_str, true)

}
