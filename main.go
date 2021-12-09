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

func backupRun(configInfo config.ModelConfig, dbinfo map[string]string, autoEncrypt string) {
	if dbinfo["password"] != "" && autoEncrypt == "yes" {
		dbinfo["password"] = keygen.AesDecryptCBC(dbinfo["password"])
	}

	base := database.BaseModel{
		TarFilename: fmt.Sprintf("%v/%v/%v-%v.tar.gz", configInfo.StoreWith["path"], dbinfo["type"], dbinfo["name"], time.Now().Format("2006.01.02.15.04.05")),
		SaveDir:     fmt.Sprintf("%v/%v/", configInfo.StoreWith["path"], dbinfo["type"]),
		BackupDir:   fmt.Sprintf("%v/%v/%v", configInfo.StoreWith["path"], dbinfo["type"], dbinfo["name"]),
		BackupNum:   configInfo.BackupNum,
		DbInfo:      dbinfo,
	}

	if _, err := os.Stat(base.BackupDir); err != nil {
		err := os.MkdirAll(base.BackupDir, 0755)
		if err != nil {
			return
		}
	}

	base.Backup()
}

func setPath() {
	// 配置BIN环境变量
	absPath1, _ := filepath.Abs("./")
	oldPath := os.Getenv("PATH")
	os.Setenv("oneBackupPath", absPath1+`/bin`)
	newPath := oldPath + ":" + absPath1 + "/bin"

	if runtime.GOOS == "linux" {
		pgStr := "cd " + os.Getenv("oneBackupPath") + " && tar xf pglib.tgz && tar xf node.tar.gz"
		cmd.Run(pgStr)
		os.Setenv("LD_LIBRARY_PATH", os.Getenv("oneBackupPath")+"/pglib")
		newPath1 := newPath + ":" + absPath1 + "/bin/node/bin"
		os.Setenv("PATH", newPath1)

		cmd.Run(fmt.Sprintf("cd %v/bin/ && chmod +x ./* ./node/bin/*", absPath1))
	}

}

func promptInformation() {
	conf := 1  // 配置、终端默认设置
	bg := 32   // 背景色、终端默认设置
	text := 40 // 前景色、红色
	fmt.Printf("\n %c[%d;%d;%dm%s%c[0m\n\n", 0x1B, conf, bg, text, "本工具适用于小量数据的备份，如遇到超大数据请选择其它工具", 0x1B)
	fmt.Printf("\n %c[%d;%d;%dm%s%c[0m\n\n", 0x1B, conf, bg, text, "5s后继续， ctrl+c 取消", 0x1B)
	time.Sleep(time.Duration(5) * time.Second)
}

func main() {

	configFile := flag.String("file", "./backupdb.yml", "config file")
	encrypt := flag.String("encrypt", "", "need encrypt string")
	autoEncrypt := flag.String("autoEncrypt", "no", "yes|no")

	mode := flag.String("mode", "backup", "run mode: backup|restore")
	dbType := flag.String("type", "", "database type: redis|mysql|mongodb|etcd|es|postgresql")
	host := flag.String("host", "", "database host: x.x.x.x")
	port := flag.String("port", "", "database port: 6379")
	db := flag.String("db", "0", "database: 0")
	username := flag.String("username", "", "database username: root")
	password := flag.String("password", "", "database password: xxx")
	src := flag.String("src", "", "restore dir/file:  such '/tmp/backupdir/redis/dump.json' ")

	flag.Parse()

	if *encrypt != "" {
		logger.Info(keygen.AesEncryptCBC(*encrypt))
		os.Exit(0)
	}

	if _, err := os.Lstat(*configFile); err != nil {
		logger.Error(*configFile, " not found!")
		flag.PrintDefaults()
		os.Exit(1)
	}

	promptInformation()

	setPath()

	absPath1, _ := filepath.Abs(*configFile)
	absDir := strings.Replace(absPath1, "\\", "/", -1)

	pathStrList := strings.Split(absDir, `/`)
	fileName := pathStrList[len(pathStrList)-1]

	filePath := strings.Replace(absDir, fileName, "", -1)
	configFlag := strings.Replace(fileName, ".yml", "", -1)

	configInfo := config.Init(*autoEncrypt, configFlag, filePath)

	if *mode == "backup" {
		for _, dbInfo := range configInfo.Databases {
			backupRun(configInfo, dbInfo, *autoEncrypt)
		}
	} else if *mode == "restore" {
		switch *dbType {
		case "redis":
			logger.Info("Restore starting")
			rdb := database.Redis{
				Host:      *host,
				Port:      *port,
				Username:  *username,
				Password:  *password,
				Databases: *db,
			}
			if err := rdb.Restore(*src); err != nil {
				logger.Error(err)
			} else {
				logger.Info("Restore success")
			}
		case "mysql":
			mysql := database.Mysql{}
			mysql.Restore(*src)
		}

	}

}
