package database

import (
	"fmt"
	"one-backup/cmd"
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

	return cmd.Run(cmdStr)
}

func (es *Elasticsearch) Restore() {
	fmt.Println("start Restore")
}
