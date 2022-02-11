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

	configFile := flag.String("file", "./backupdb.yml", "config file")
	encrypt := flag.String("encrypt", "", "Input a string and output an encrypted string")
	autoEncrypt := flag.String("autoEncrypt", "yes", "yes|no")
	mode := flag.String("mode", "backup", "run mode: backup|restore")
	dbType := flag.String("type", "", "database type: redis|mysql|mongodb|etcd|es|postgresql")
	host := flag.String("host", "", "database host: x.x.x.x")
	port := flag.String("port", "", "database port: 6379")
	db := flag.String("db", "0", "database: 0")
	username := flag.String("username", "root", "database username")
	password := flag.String("password", "", "database password: xxx")
	authdb := flag.String("authdb", "admin", "mongo authdb: admin")
	https := flag.String("https", "no", "etcd https")
	cacert := flag.String("etcdcacert", "/etc/kubernetes/pki/etcd/ca.crt", "etcd cacert")
	cert := flag.String("etcdcert", "/etc/kubernetes/pki/etcd/server.crt", "etcd cert")
	certkey := flag.String("etcdkey", "/etc/kubernetes/pki/etcd/server.key", "etcd key")
	etcdservice := flag.String("etcdservice", "etcd", "etcdservice : etcd.service")
	etcddatadir := flag.String("etcddatadir", "/var/lib/etcd", "etcd data-dir: /var/lib/etcd")
	etcdname := flag.String("etcdname", "etcd1", "etcdname: etcd1")
	etcdcluster := flag.String("etcdcluster", "", "etcdcluster: etcd1=etcd1:2379,etcd2=etcd2:2379,etcd3=etcd3:2379")
	etcdclustertoken := flag.String("etcdclustertoken", "", "etcdclustertoken: ")
	dockername := flag.String("dockername", "", "dockername: etcd1")
	dockernetwork := flag.String("dockernetwork", "host", "dockernetwork: host|nat")
	sshhost := flag.String("sshhost", "", "sshhost: 192.168.12.1")
	sshport := flag.String("sshport", "22", "sshport: 22")
	sshuser := flag.String("sshuser", "root", "sshuser: root")
	sshpassword := flag.String("sshpassword", "", "sshpassword: root")
	src := flag.String("src", "", "restore file/dir:  such './dump.json or ./mongodb-2021.12.27.01.35.24/'")

	flag.Parse()

	if *encrypt != "" {
		logger.Info(keygen.AesEncryptCBC(*encrypt, "pass"))
		os.Exit(0)
	}

	promptInformation()

	execFile, _ := filepath.Abs(os.Args[0])
	execFileTmp1 := strings.Split(execFile, `/`)
	execPath := strings.Join(execFileTmp1[0:len(execFileTmp1)-1], `/`)
	setPath(execPath)

	if *mode == "backup" {
		if _, err := os.Lstat(*configFile); err != nil {
			logger.Error(*configFile, " not found!")
			flag.PrintDefaults()
			os.Exit(1)
		}

		absPath, _ := filepath.Abs(*configFile)
		absPath_format := strings.Replace(absPath, "\\", "/", -1)

		pathStrList := strings.Split(absPath_format, `/`)
		fileName := pathStrList[len(pathStrList)-1]

		filePath := strings.Replace(absPath_format, fileName, "", -1)
		configFlag := strings.Replace(fileName, ".yml", "", -1)

		configInfo := config.Init(*autoEncrypt, configFlag, filePath)

		if configInfo.StoreWith["password"] != "" && *autoEncrypt == "yes" {
			if configInfo.IsEncrypt {
				configInfo.StoreWith["password"] = keygen.AesDecryptCBC(configInfo.StoreWith["password"], "pass")
				if configInfo.StoreWith["password"] == "decrypted error" {
					logger.Error("password decrypted error: storewith")
				}
			}
		}
		for _, dbInfo := range configInfo.Databases {
			database.Run(configInfo, dbInfo, *autoEncrypt)
		}
	} else if *mode == "restore" {
		base := database.BaseModel{}
		base.DbInfo = map[string]string{
			"execpath": execPath,
			"dbType":   *dbType,
			"host":     *host,
			"port":     *port,
			"username": *username,
			"password": *password,
			"db":       *db,
			"src":      *src,
			"authdb":   *authdb,
			// etcd https
			"https": *https, "cacert": *cacert, "cert": *cert, "key": *certkey,
			// etcd-info
			"etcddatadir": *etcddatadir, "etcdName": *etcdname, "etcdservice": *etcdservice,
			// etcd-cluster
			"etcdCluster": *etcdcluster, "etcdCluserToken": *etcdclustertoken,
			// etcd-remote host
			"sshhost": *sshhost, "sshport": *sshport, "sshuser": *sshuser, "sshpassword": *sshpassword,
			// etcd-docker info
			"dockername": *dockername, "dockernetwork": *dockernetwork,
		}
		database.Restore(base)
	}

}
