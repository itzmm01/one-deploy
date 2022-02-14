package main

import (
	"flag"
	"fmt"
	"one-backup/cmd"
	"one-backup/config"
	"one-backup/database"
	"one-backup/keygen"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/wonderivan/logger"
)

func setPath(execPath string) {
	// 配置BIN环境变量
	oldPath := os.Getenv("PATH")
	os.Setenv("oneBackupPath", execPath+`/bin`)
	newPath := oldPath + ":" + execPath + "/bin"

	if runtime.GOOS == "linux" {
		pgStr := "cd " + os.Getenv("oneBackupPath") + " && tar xf pglib.tgz && tar xf node.tar.gz"
		cmd.Run(pgStr, config.Debug)
		os.Setenv("LD_LIBRARY_PATH", os.Getenv("oneBackupPath")+"/pglib")
		newPath1 := newPath + ":" + execPath + "/bin/node/bin"
		os.Setenv("PATH", newPath1)
		cmd.Run(fmt.Sprintf("cd %v/bin/ && chmod +x * ./node/bin/*", execPath), config.Debug)
	}

}

func promptInformation() {
	conf := 1  // 配置、终端默认设置
	bg := 32   // 背景色、终端默认设置
	text := 40 // 前景色、红色
	fmt.Printf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, conf, bg, text, "本工具适用于少量元数据的备份，如遇到超大数据请选择其它工具", 0x1B)
	fmt.Printf("%c[%d;%d;%dm%s%c[0m\n", 0x1B, conf, bg, text, "5s后继续， ctrl+c 取消", 0x1B)
	time.Sleep(time.Duration(5) * time.Second)
}

func main() {
	argsMap := map[string]*string{}

	argsMap["configfile"] = flag.String("file", "./backupdb.yml", "config file")
	argsMap["encrypt"] = flag.String("encrypt", "", "Input a string and output an encrypted string")
	argsMap["autoencrypt"] = flag.String("autoEncrypt", "yes", "yes|no")
	argsMap["mode"] = flag.String("mode", "backup", "run mode: backup|restore")
	argsMap["dbtype"] = flag.String("type", "", "database type: redis|mysql|mongodb|etcd|es|postgresql")
	argsMap["dbhost"] = flag.String("host", "", "database host: x.x.x.x")
	argsMap["dbport"] = flag.String("port", "", "database port: 6379")
	argsMap["db"] = flag.String("db", "0", "database: 0")
	argsMap["dbusername"] = flag.String("username", "root", "database username")
	argsMap["dbpassword"] = flag.String("password", "", "database password: xxx")
	argsMap["authdb"] = flag.String("authdb", "admin", "mongo authdb: admin")
	argsMap["https"] = flag.String("https", "no", "etcd https")
	argsMap["cacert"] = flag.String("etcdcacert", "/etc/kubernetes/pki/etcd/ca.crt", "etcd cacert")
	argsMap["cert"] = flag.String("etcdcert", "/etc/kubernetes/pki/etcd/server.crt", "etcd cert")
	argsMap["certkey"] = flag.String("etcdkey", "/etc/kubernetes/pki/etcd/server.key", "etcd key")
	argsMap["etcdservice"] = flag.String("etcdservice", "etcd", "etcdservice : etcd.service")
	argsMap["etcddatadir"] = flag.String("etcddatadir", "/var/lib/etcd", "etcd data-dir: /var/lib/etcd")
	argsMap["etcdname"] = flag.String("etcdname", "etcd1", "etcdname: etcd1")
	argsMap["etcdcluster"] = flag.String(
		"etcdcluster", "", "etcdcluster: etcd1=etcd1:2379,etcd2=etcd2:2379,etcd3=etcd3:2379",
	)
	argsMap["etcdclustertoken"] = flag.String("etcdclustertoken", "", "etcdclustertoken: ")
	argsMap["dockername"] = flag.String("dockername", "", "dockername: etcd1")
	argsMap["dockernetwork"] = flag.String("dockernetwork", "host", "dockernetwork: host|nat")
	argsMap["sshhost"] = flag.String("sshhost", "", "sshhost: 192.168.12.1")
	argsMap["sshport"] = flag.String("sshport", "22", "sshport: 22")
	argsMap["sshuser"] = flag.String("sshuser", "root", "sshuser: root")
	argsMap["sshpassword"] = flag.String("sshpassword", "", "sshpassword: root")
	argsMap["src"] = flag.String("src", "", "restore file/dir:  such './dump.json or ./mongodb-2021.12.27.01.35.24/'")

	flag.Parse()

	if *argsMap["encrypt"] != "" {
		logger.Info(keygen.AesEncryptCBC(*argsMap["encrypt"], "pass"))
		os.Exit(0)
	}

	promptInformation()

	execFile, _ := filepath.Abs(os.Args[0])
	execFileTmp1 := strings.Split(execFile, `/`)
	execPath := strings.Join(execFileTmp1[0:len(execFileTmp1)-1], `/`)
	setPath(execPath)
	database.Init(execPath, argsMap)

}
