package database

import (
	"fmt"
	"one-backup/cmd"
)

// info
type Mongodb struct {
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
	// 数据库
	Database string
	// 认证数据库
	AuthDb string
}

// backup
func (ctx Mongodb) Backup() error {
	cmdStr := fmt.Sprintf("mongodump -h %v --port %v ", ctx.Host, ctx.Port)

	if ctx.Username != "" && ctx.Password != "" {
		cmdStr = cmdStr + fmt.Sprintf("-u %v -p '%v' --authenticationDatabase %v ", ctx.Username, ctx.Password, ctx.AuthDb)
	}

	if ctx.Database != "alldatabase" {
		cmdStr = cmdStr + fmt.Sprintf("-d %v ", ctx.Database)
	}

	cmdStr = cmdStr + fmt.Sprintf("-o  %v/ ", ctx.BackupDir)
	return cmd.Run(cmdStr, Debug)
}
