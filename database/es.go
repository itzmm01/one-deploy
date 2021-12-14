package database

import (
	"fmt"
	"one-backup/cmd"
	"one-backup/keygen"
	"one-backup/tool"
)

type Elasticsearch struct {
	/*
		host: 192.168.146.134
		index: abc
		name: es
		password: BhaBUTSg3lMXHLVUkHmOfw==
		port: "9200"
		type: es
		username: root
	*/
	TarFilename string
	SaveDir     string
	BackupDir   string
	Host        string
	Port        string
	Username    string
	Password    string
	Index       string
}

func (ctx *Elasticsearch) Backup() error {
	cmdStr := ""
	if ctx.Username != "" && ctx.Password != "" {
		cmdStr = fmt.Sprintf("elasticdump --input=http://%v:%v@%v:%v/%v --output=/%v/%v.json --all=true ", ctx.Username, ctx.Password, ctx.Host, ctx.Port, ctx.Index, ctx.BackupDir, ctx.Index)
	} else {
		cmdStr = fmt.Sprintf("elasticdump --input=http://%v:%v/%v --output=/%v/%v.json --all=true ", ctx.Host, ctx.Port, ctx.Index, ctx.BackupDir, ctx.Index)
	}
	if err := cmd.Run(cmdStr, Debug); err == nil {
		keygen.AesEncryptCBCFile(fmt.Sprintf("%v/%v.json", ctx.BackupDir, ctx.Index), fmt.Sprintf("%v/%v-Encrypt.json", ctx.BackupDir, ctx.Index))
		return cmd.Run(fmt.Sprintf("rm -f %v/%v.json", ctx.BackupDir, ctx.Index), false)
	} else {
		return err
	}
}

func (ctx *Elasticsearch) Restore(filePath string) error {
	dstPath := "/tmp/" + tool.RandomString(30)
	keygen.AesDecryptCBCFile(filePath, dstPath)
	cmdStr := ""
	if ctx.Username != "" && ctx.Password != "" {
		cmdStr = fmt.Sprintf("elasticdump --input=%v --output=http://%v:%v@%v:%v --all=true ", dstPath, ctx.Username, ctx.Password, ctx.Host, ctx.Port)
	} else {
		cmdStr = fmt.Sprintf("elasticdump --input=%v --output=http://%v:%v --all=true", dstPath, ctx.Host, ctx.Port)
	}

	if err := cmd.Run(cmdStr, true); err != nil {
		return err
	} else {
		return cmd.Run("rm -f "+dstPath, false)
	}
}
