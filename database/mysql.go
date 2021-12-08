package database

import (
	"fmt"
	"one-backup/cmd"
	"one-backup/keygen"
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
	Db          string
}

func (ctx *Mysql) Backup() error {
	cmdStr := fmt.Sprintf(
		"mysqldump -h %v -P %v -u%v -p'%v' ", ctx.Host, ctx.Port, ctx.Username, ctx.Password,
	)
	if ctx.Db == "alldatabase" {
		cmdStr = cmdStr + "--all-databases "
	} else {
		cmdStr = cmdStr + "--databases " + ctx.Db
	}

	cmdStr = cmdStr + fmt.Sprintf(" > %v/%v.sql", ctx.BackupDir, ctx.Db)

	err := cmd.Run(cmdStr)
	if err == nil {
		keygen.AesEncryptCBCFile(fmt.Sprintf("%v/%v.sql", ctx.BackupDir, ctx.Db), fmt.Sprintf("%v/%v-Encrypt.sql", ctx.BackupDir, ctx.Db))
		cmd.Run(fmt.Sprintf("rm -f %v/%v.sql", ctx.BackupDir, ctx.Db))
		return nil
	} else {
		return err
	}

}

func (ctx Mysql) Restore() {
	fmt.Println("start Restore")
}
