package database

import (
	"fmt"
	"one-backup/cmd"
	"one-backup/config"
	"one-backup/ssh"
	"os"
	"strconv"
	"strings"

	"github.com/wonderivan/logger"
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
	Key        string
	DockerName string
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
			"--cacert=%v --cert=%v --key=%v --endpoints=%v snapshot save %v/etcd.db",
			ctx.Cacert, ctx.Cert, ctx.Key, ctx.Host, ctx.BackupDir,
		)
	} else {
		cmdStr += fmt.Sprintf(
			"--endpoints=%v snapshot save %v/etcd.db",
			ctx.Host, ctx.BackupDir,
		)
	}
	return cmd.Run(cmdStr, config.Debug)
}

func getDockerCmdStr(dockerName, execPath, srcFilePath, cmdStr string, other map[string]string) string {
	if other["dockername"] != "" && other["dockernetwork"] == "nat" {
		cpCmdStr := fmt.Sprintf("docker cp %v %v:/root/ ; docker cp %v %v:/tmp/",
			execPath, dockerName, srcFilePath, dockerName,
		)
		dockerCleanStr := fmt.Sprintf("docker exec -i %s /bin/sh -c 'ls -d %v &>/dev/null && mv -f %v %v-%v || echo' ",
			dockerName, other["datadir"], other["datadir"], other["datadir"], other["nowtime"],
		)
		if other["sshhost"] == "" {
			cmd.Run(cpCmdStr, config.Debug)
			cmd.Run(dockerCleanStr, config.Debug)
		} else {
			sshClient := getSSH(other)
			if res, err := sshClient.RunShell(cpCmdStr); err != nil {
				logger.Error(res, err)
				os.Exit(1)
			}
			if res, err := sshClient.RunShell(dockerCleanStr); err != nil {
				logger.Error(res, err)
				os.Exit(1)
			}
		}
		return fmt.Sprintf("docker exec -i %s /bin/sh -c '%s snapshot restore /tmp/etcd.db'", dockerName, cmdStr)

	} else {
		CleanStr := fmt.Sprintf("ls -d %v &>/dev/null && mv -f %v %v-%v || echo",
			other["datadir"], other["datadir"], other["datadir"], other["nowtime"],
		)
		if other["sshhost"] == "" {
			cmd.Run(CleanStr, config.Debug)
		} else {
			sshClient := getSSH(other)
			if res, err := sshClient.RunShell(CleanStr); err != nil {
				logger.Error(res, err)
				os.Exit(1)
			}
			cmdStr = "/root/" + cmdStr
		}
		return cmdStr + "snapshot restore " + srcFilePath
	}
}

func getSSH(other map[string]string) *ssh.ClientConfig {
	sshClient := new(ssh.ClientConfig)
	sshPort, _ := strconv.ParseInt(other["sshport"], 10, 64)
	sshClient.CreateClient(other["sshhost"], sshPort, other["sshuser"], other["sshpassword"])
	return sshClient
}

func restartEtcd(dockerName, serviceName string, other map[string]string) error {
	restartCmdStr := ""
	if dockerName != "" {
		restartCmdStr = "docker restart " + dockerName
	} else {
		restartCmdStr = "systemctl restart " + serviceName
	}
	if other["cluster"] != "" {
		sshClient := getSSH(other)
		_, err := sshClient.RunShell(restartCmdStr)
		return err
	} else {
		return cmd.Run(restartCmdStr, false)
	}
}

// Restore
// func (ctx Etcd) Restore(filePath, dataDir, etcdName, cluster, clusertoken, docker string) error {
func (ctx Etcd) Restore(filePath string, other map[string]string) error {
	cmdStr := "etcdctl --command-timeout=300s "
	if ctx.Username != "" && ctx.Password != "" {
		cmdStr += fmt.Sprintf("--user=%v:%v ", ctx.Username, ctx.Password)
	}

	if ctx.Https == "yes" {
		cmdStr += fmt.Sprintf(
			"--cacert=%v --cert=%v --key=%v  ", ctx.Cacert, ctx.Cert, ctx.Key,
		)
	}

	if other["cluster"] != "" {
		nodeInfo := map[string]string{}
		clusterInfo := strings.Split(other["cluster"], `,`)
		for _, node := range clusterInfo {
			hostInfo := strings.Split(node, `=`)
			nodeInfo[hostInfo[0]] = hostInfo[1]
		}
		cmdStr += fmt.Sprintf(
			"--name %v --initial-cluster=\"%v\" --initial-cluster-token=%v --initial-advertise-peer-urls=%v --data-dir=%v ",
			other["name"], other["cluster"], other["clustertoken"], nodeInfo[other["name"]], other["datadir"],
		)
	} else {
		cmdStr += fmt.Sprintf("--data-dir=%v --endpoints=%v ", other["datadir"], other["name"])
	}
	execPath := other["execpath"]
	srcFilePath := filePath

	if other["sshhost"] != "" {
		execPath = "/root/one-backup/"
		srcFilePath = "/tmp/etcd.db"

		sshClient := getSSH(other)
		sshClient.Upload(other["execpath"]+"/bin/etcdctl", "/root/etcdctl")
		sshClient.RunShell("chmod +x /root/etcdctl")
		sshClient.Upload(filePath, srcFilePath)

		cmdStr = getDockerCmdStr(ctx.DockerName, execPath, srcFilePath, cmdStr, other)
		if res, err := sshClient.RunShell("ETCDCTL_API=3 " + cmdStr); err != nil {
			logger.Error(res)
			return err
		} else {
			return restartEtcd(ctx.DockerName, other["etcdservice"], other)
		}

	} else {
		cmdStr = getDockerCmdStr(ctx.DockerName, execPath, srcFilePath, cmdStr, other)
		if err := cmd.Run("ETCDCTL_API=3 "+cmdStr, config.Debug); err != nil {
			return err
		} else {
			return restartEtcd(ctx.DockerName, other["etcdservice"], other)
		}
	}
}
