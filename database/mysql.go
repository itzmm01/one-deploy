package database

import (
	"fmt"
	"one-backup/cmd"
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
	return cmd.Run(cmdStr)
}

func (ctx Mysql) Restore() {
	fmt.Println("start Restore")
}
