package database

import (
	"fmt"
	"one-backup/cmd"
)

type Redis struct {
	/*
		name: redis
		# 数据库类型
		type: redis
		# 备份方式目前支持sync
		mode: sync
		# 数据库IP
		host: 192.168.146.134
		# 端口
		port: 6379
		# 密码
		password: Amt_2018
	*/
	TarFilename string
	SaveDir     string
	BackupDir   string
	Name        string
	Host        string
	Port        string
	Password    string
}

func (ctx Redis) Backup() error {
	//
	cmdStr := fmt.Sprintf("redis-cli -h %v -p %v ", ctx.Host, ctx.Port)
	if ctx.Password != "" {
		cmdStr = cmdStr + fmt.Sprintf("-a '%v' ", ctx.Password)
	}
	cmdStr = cmdStr + fmt.Sprintf("--rdb %v/dump.rdb", ctx.BackupDir)
	return cmd.Run(cmdStr)
}
