package database

import (
	"fmt"
	"one-backup/cmd"
	"os"
)

type Etcd struct {
	/*
			name: etcd
		    # 数据库类型
		    type: etcd
		    # etcd链接信息
		    host: 127.0.0.1
		    port: 49514
		    # 是否使用https, yes|no
		    https: no
		    # ca证书路径
		    cacert: /etc/etcd/ssl/ca.pem
		    # 客户端证书路径
		    cert: /etc/etcd/ssl/etcd.pem
		    # 客户端密钥路径
		    key: /etc/etcd/ssl/etcd-key.pem
	*/
	TarFilename string
	SaveDir     string
	BackupDir   string
	Name        string
	Host        string
	Port        string
	Username    string
	Password    string
	Https       string
	Cacert      string
	Cert        string
	Key         string
}

func (ctx Etcd) Backup() error {
	os.Setenv("ETCDAPI", "3")
	cmdStr := "etcdctl --command-timeout=300s "
	if ctx.Cacert == "" {
		ctx.Cacert = "/etc/kubernetes/pki/etcd/ca.crt"
	}
	if ctx.Cert == "" {
		ctx.Cert = "/etc/kubernetes/pki/etcd/server.crt"
	}
	if ctx.Key == "" {
		ctx.Key = "/etc/kubernetes/pki/etcd/server.key"
	}
	if ctx.Username != "" {
		cmdStr += fmt.Sprintf("--user=%v:%v ", ctx.Username, ctx.Password)
	}
	if ctx.Https == "yes" {
		cmdStr += fmt.Sprintf(
			"--cacert=%v --cert=%v --key=%v --endpoints=https://%v:%v snapshot save %v/etcd.db",
			ctx.Cacert, ctx.Cert, ctx.Key, ctx.Host, ctx.Port, ctx.BackupDir,
		)
	} else {
		cmdStr += fmt.Sprintf(
			"--endpoints=http://%v:%v snapshot save %v/etcd.db",
			ctx.Host, ctx.Port, ctx.BackupDir,
		)
	}
	return cmd.Run(cmdStr, Debug)
}
