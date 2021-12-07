package database

import (
	"fmt"
	"one-backup/cmd"
)

type Mongodb struct {
	/*
			name: mongo
		    # 数据库类型
		    type: mongodb
		    # 数据库IP
		    host: 192.168.146.134
		    # 端口
		    port: 27017
		    # 账号
		    username: root
		    # 密码
		    password: Amt_2018
		    # 需要备份的数据库，多个用英文逗号隔开, alldatabase 代表所有数据库
		    database: alldatabase
		    # 验证用户数据库
		    authdb: "admin"
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
	AuthDb      string
}

func (ctx Mongodb) Backup() error {
	cmdStr := fmt.Sprintf("mongodump -h %v --port %v ", ctx.Host, ctx.Port)

	if ctx.Username != "" && ctx.Password != "" {
		cmdStr = cmdStr + fmt.Sprintf("-u %v -p '%v' --authenticationDatabase %v ", ctx.Username, ctx.Password, ctx.AuthDb)
	}

	if ctx.Db != "alldatabase" {
		cmdStr = cmdStr + fmt.Sprintf("-d %v ", ctx.Db)
	}

	cmdStr = cmdStr + fmt.Sprintf("-o  %v/ ", ctx.BackupDir)
	return cmd.Run(cmdStr)
}
