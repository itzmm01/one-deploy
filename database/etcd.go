package database

import (
	"fmt"
	"one-backup/cmd"
	"os"
)

// info
type Etcd struct {
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
	// 是否使用https
	Https string
	// 证书路径
	Cacert string
	// 证书路径
	Cert string
	// 证书路径
	Key string
}

// backup
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
