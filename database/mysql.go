package database

import (
	"fmt"
	"one-backup/cmd"
	"one-backup/keygen"
	"one-backup/tool"
)

// info
type Mysql struct {
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
func (ctx *Mysql) Backup() error {
	cmdStr := fmt.Sprintf(
		"mysqldump --lock-tables=0 -h %v -P %v -u%v -p'%v' ", ctx.Host, ctx.Port, ctx.Username, ctx.Password,
	)
	if ctx.Database == "alldatabase" {
		cmdStr = cmdStr + "--all-databases "
	} else {
		cmdStr = cmdStr + "--databases " + ctx.Database
	}

	cmdStr = cmdStr + fmt.Sprintf(" > %v/%v.sql", ctx.BackupDir, ctx.Database)

	err := cmd.Run(cmdStr, Debug)
	if err == nil {
		keygen.AesEncryptCBCFile(
			fmt.Sprintf("%v/%v.sql", ctx.BackupDir, ctx.Database),
			fmt.Sprintf("%v/%v-Encrypt.sql", ctx.BackupDir, ctx.Database),
		)
		return cmd.Run(fmt.Sprintf("rm -f %v/%v.sql", ctx.BackupDir, ctx.Database), Debug)
	} else {
		return err
	}

}

// restore
func (ctx Mysql) Restore(filepath string) error {
	dstPath := "/tmp/" + tool.RandomString(30)
	keygen.AesDecryptCBCFile(filepath, dstPath)
	cmdStr := fmt.Sprintf(
		"cat '%v' | mysql -h %v -P %v -u%v -p%v",
		dstPath, ctx.Host, ctx.Port, ctx.Username, ctx.Password,
	)

	if err := cmd.Run(cmdStr, Debug); err != nil {
		return err
	} else {
		return cmd.Run("rm -f "+dstPath, false)
	}

}
