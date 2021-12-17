package database

import (
	"fmt"
	"one-backup/cmd"
	"one-backup/keygen"
	"one-backup/tool"
)

// info
type Elasticsearch struct {
	// 压缩包文件名
	TarFilename string
	// 保存目录
	SaveDir string
	// 备份目录
	BackupDir string
	// 主机
	Host string
	// 端口
	Port string
	// 账号
	Username string
	// 密码
	Password string
	// 索引
	Index string
}

// backup
func (ctx *Elasticsearch) Backup() error {
	cmdStr := ""
	if ctx.Username != "" && ctx.Password != "" {
		cmdStr = fmt.Sprintf(
			"elasticdump --input=http://%v:%v@%v:%v/%v --output=/%v/%v.json --all=true ",
			ctx.Username, ctx.Password, ctx.Host, ctx.Port, ctx.Index, ctx.BackupDir, ctx.Index,
		)
	} else {
		cmdStr = fmt.Sprintf(
			"elasticdump --input=http://%v:%v/%v --output=/%v/%v.json --all=true ",
			ctx.Host, ctx.Port, ctx.Index, ctx.BackupDir, ctx.Index,
		)
	}
	if err := cmd.Run(cmdStr, Debug); err == nil {
		keygen.AesEncryptCBCFile(
			fmt.Sprintf("%v/%v.json", ctx.BackupDir, ctx.Index),
			fmt.Sprintf("%v/%v-Encrypt.json", ctx.BackupDir, ctx.Index),
		)
		return cmd.Run(fmt.Sprintf("rm -f %v/%v.json", ctx.BackupDir, ctx.Index), false)
	} else {
		return err
	}
}

// restore
func (ctx *Elasticsearch) Restore(filePath string) error {
	dstPath := "/tmp/" + tool.RandomString(30)
	keygen.AesDecryptCBCFile(filePath, dstPath)
	cmdStr := ""
	if ctx.Username != "" && ctx.Password != "" {
		cmdStr = fmt.Sprintf(
			"elasticdump --input=%v --output=http://%v:%v@%v:%v --all=true ",
			dstPath, ctx.Username, ctx.Password, ctx.Host, ctx.Port,
		)
	} else {
		cmdStr = fmt.Sprintf(
			"elasticdump --input=%v --output=http://%v:%v --all=true",
			dstPath, ctx.Host, ctx.Port,
		)
	}

	if err := cmd.Run(cmdStr, true); err != nil {
		cmd.Run("rm -f "+dstPath, false)
		return err
	} else {
		return cmd.Run("rm -f "+dstPath, false)
	}
}
