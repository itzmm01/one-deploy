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

const Debug = false

func setPath() {
	// 配置BIN环境变量
	absPath1, _ := filepath.Abs("./")
	oldPath := os.Getenv("PATH")
	os.Setenv("oneBackupPath", absPath1+`/bin`)
	newPath := oldPath + ":" + absPath1 + "/bin"

	if runtime.GOOS == "linux" {
		pgStr := "cd " + os.Getenv("oneBackupPath") + " && tar xf pglib.tgz && tar xf node.tar.gz"
		cmd.Run(pgStr, Debug)
		os.Setenv("LD_LIBRARY_PATH", os.Getenv("oneBackupPath")+"/pglib")
		newPath1 := newPath + ":" + absPath1 + "/bin/node/bin"
		os.Setenv("PATH", newPath1)

		cmd.Run(fmt.Sprintf("cd %v/bin/ && chmod +x ./* ./node/bin/*", absPath1), Debug)
	}

}

func promptInformation() {
	conf := 1  // 配置、终端默认设置
	bg := 32   // 背景色、终端默认设置
	text := 40 // 前景色、红色
	fmt.Printf("\n %c[%d;%d;%dm%s%c[0m\n\n", 0x1B, conf, bg, text, "本工具适用于少量元数据的备份，如遇到超大数据请选择其它工具", 0x1B)
	fmt.Printf("\n %c[%d;%d;%dm%s%c[0m\n\n", 0x1B, conf, bg, text, "5s后继续， ctrl+c 取消", 0x1B)
	time.Sleep(time.Duration(5) * time.Second)
}

func main() {

	configFile := flag.String("file", "./backupdb.yml", "config file")
	encrypt := flag.String("encrypt", "", "need encrypt string")
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
	cacert := flag.String("cacert", "/etc/kubernetes/pki/etcd/ca.crt", "etcd cacert")
	cert := flag.String("cert", "/etc/kubernetes/pki/etcd/server.crt", "etcd cert")
	certkey := flag.String("Key", "/etc/kubernetes/pki/etcd/server.key", "etcd key")
	dataDir := flag.String("datadir", "/var/lib/etcd", "etcd data-dir")
	src := flag.String("src", "", "restore file/dir:  such './dump.json or ./mongodb-2021.12.27.01.35.24/'")

	flag.Parse()

	if *encrypt != "" {
		logger.Info(keygen.AesEncryptCBC(*encrypt, "pass"))
		os.Exit(0)
	}

	promptInformation()

	setPath()

	if *mode == "backup" {
		if _, err := os.Lstat(*configFile); err != nil {
			logger.Error(*configFile, " not found!")
			flag.PrintDefaults()
			os.Exit(1)
		}

		absPath1, _ := filepath.Abs(*configFile)
		absDir := strings.Replace(absPath1, "\\", "/", -1)

		pathStrList := strings.Split(absDir, `/`)
		fileName := pathStrList[len(pathStrList)-1]

		filePath := strings.Replace(absDir, fileName, "", -1)
		configFlag := strings.Replace(fileName, ".yml", "", -1)

		configInfo := config.Init(*autoEncrypt, configFlag, filePath)

		for _, dbInfo := range configInfo.Databases {
			database.Run(configInfo, dbInfo, *autoEncrypt)
		}
	} else if *mode == "restore" {
		base := database.BaseModel{}
		base.DbInfo = map[string]string{
			"dbType":   *dbType,
			"host":     *host,
			"port":     *port,
			"username": *username,
			"password": *password,
			"db":       *db,
			"src":      *src,
			"authdb":   *authdb,
			"https":    *https,
			"cacert":   *cacert,
			"cert":     *cert,
			"key":      *certkey,
			"dataDir":  *dataDir,
		}
		database.Restore(base)
	}

}
