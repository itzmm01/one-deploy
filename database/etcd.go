package database

import (
	"fmt"
	"one-backup/cmd"
	"one-backup/tool"
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

// Restore
func (ctx Etcd) Restore(filePath, dataDir string) error {
	message := `提示:
需要删除原有数据总目录，不然恢复的时候有有各种问题。命令中需要指定数据目录 -datadir /var/lib/etcd，否则会产生默认工作目录，一个default的名称
`
	if tool.PleaseConfirm(message) != "YES" {
		return nil
	}

	/*
		etcdctl snapshot restore snap1
		--name etcd-41
		--initial-cluster
		etcd-41=http://192.168.31.41:2380,etcd-42=http://192.168.31.42:2380,etcd-43=http://192.168.31.43:2380
		--initial-advertise-peer-urls http://192.168.31.41:2380
		--data-dir /var/lib/etcd/cluster.etcd
		快照恢复时，会重新生成客户端id和集群id，所有的节点统一使用一个快照恢复集群
	*/

	os.Setenv("ETCDAPI", "3")
	cmdStr := "etcdctl --command-timeout=300s "

	if ctx.Username != "" {
		cmdStr += fmt.Sprintf("--user=%v:%v ", ctx.Username, ctx.Password)
	}
	if ctx.Https == "yes" {
		cmdStr += fmt.Sprintf(
			"--cacert=%v --cert=%v --key=%v --endpoints=https://%v:%v snapshot restore %v --data-dir=\"%v\"",
			ctx.Cacert, ctx.Cert, ctx.Key, ctx.Host, ctx.Port, filePath, dataDir,
		)
	} else {
		cmdStr += fmt.Sprintf(
			"--endpoints=http://%v:%v snapshot restore %v --data-dir=\"%v\"", ctx.Host, ctx.Port, filePath, dataDir,
		)
	}

	return cmd.Run(cmdStr, false)
}
