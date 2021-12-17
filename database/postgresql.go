package database

import (
	"fmt"
	"one-backup/cmd"
	"one-backup/keygen"
	"one-backup/tool"
	"os"
)

// info
type Postgresql struct {
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
}

// backup
func (ctx Postgresql) Backup() error {
	// pg_dump只支持单个数据库
	cmdStr := fmt.Sprintf("pg_dump -h %v -p %v -U %v ", ctx.Host, ctx.Port, ctx.Username)
	if ctx.Password != "" {
		os.Setenv("PGPASSWORD", ctx.Password)
	}

	cmdStr = cmdStr + fmt.Sprintf("-d %v -f %v/%v.sql", ctx.Database, ctx.BackupDir, ctx.Database)

	if err := cmd.Run(cmdStr, Debug); err == nil {
		keygen.AesEncryptCBCFile(
			fmt.Sprintf("%v/%v.sql", ctx.BackupDir, ctx.Database),
			fmt.Sprintf("%v/%v-Encrypt.sql", ctx.BackupDir, ctx.Database),
		)
		return cmd.Run(fmt.Sprintf("rm -f %v/%v.sql", ctx.BackupDir, ctx.Database), false)
	} else {
		return err
	}
}

// restore
func (ctx Postgresql) Restore(filePath string) error {
	dstPath := "/tmp/" + tool.RandomString(30)
	keygen.AesDecryptCBCFile(filePath, dstPath)

	cmdStr := fmt.Sprintf("psql  -h %v -p %v -U %v ", ctx.Host, ctx.Port, ctx.Username)
	cmdStrCreate := fmt.Sprintf(
		"num=`%v -c '\\l'|grep %v|wc -l`; if [ $num -eq 0 ]; then %v -c 'CREATE DATABASE %v'; fi",
		cmdStr, ctx.Database, cmdStr, ctx.Database,
	)
	cmdStrRestore := cmdStr + fmt.Sprintf("-d %v -f %v", ctx.Database, dstPath)

	if ctx.Password != "" {
		os.Setenv("PGPASSWORD", ctx.Password)
	}

	if err := cmd.Run(cmdStrCreate, true); err != nil {
		return err
	}

	if err := cmd.Run(cmdStrRestore, true); err != nil {
		return err
	} else {
		return cmd.Run("rm -f "+dstPath, false)
	}
}
