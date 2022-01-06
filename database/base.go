package database

import (
	"fmt"
	"log"
	"one-backup/archive"
	"one-backup/config"
	"one-backup/ftpclient"
	"one-backup/keygen"
	"one-backup/ssh"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/wonderivan/logger"
)

const Debug = false

// base info
type BaseModel struct {
	TarFilename string
	TarName     string
	SaveDir     string
	BackupDir   string
	SaveInfo    map[string]string
	BackupNum   int
	DbInfo      map[string]string
}

// 清理文件
func cleanHistoryFile(path, match string, saveNum int) {
	var fileList []string
	filepathNames, err := filepath.Glob(filepath.Join(path, match))
	if err != nil {
		log.Fatal(err)

	}
	for i := range filepathNames {
		fileList = append(fileList, filepathNames[i])
	}
	sort.Strings(fileList)
	countNum := len(fileList) - saveNum
	if countNum > 0 {
		for _, v := range fileList[0:countNum] {
			err := os.Remove(v)
			if err != nil {
				fmt.Printf("remove %v fail\n", err)
			}
		}
	}
}

// run
func Run(configInfo config.ModelConfig, dbinfo map[string]string, autoEncrypt string) {
	if dbinfo["password"] != "" && autoEncrypt == "yes" {
		decryptRes := keygen.AesDecryptCBC(dbinfo["password"], "pass")
		if decryptRes != "base64 error" {
			dbinfo["password"] = decryptRes
		}
	}
	if configInfo.StoreWith["password"] != "" && autoEncrypt == "yes" {
		decryptRes := keygen.AesDecryptCBC(configInfo.StoreWith["password"], "pass")
		if decryptRes != "base64 error" {
			configInfo.StoreWith["password"] = decryptRes
		}
	}
	nameDir := fmt.Sprintf("%v-%v", dbinfo["name"], time.Now().Format("2006.01.02.15.04.05"))
	base := BaseModel{
		TarFilename: fmt.Sprintf("%v/%v/%v.tar.gz", configInfo.StoreWith["path"], dbinfo["type"], nameDir),
		TarName:     fmt.Sprintf("%v.tar.gz", nameDir),
		SaveDir:     fmt.Sprintf("%v/%v/", configInfo.StoreWith["path"], dbinfo["type"]),
		BackupDir:   fmt.Sprintf("%v/%v/%v", configInfo.StoreWith["path"], dbinfo["type"], nameDir),
		SaveInfo:    configInfo.StoreWith,
		BackupNum:   configInfo.BackupNum,
		DbInfo:      dbinfo,
	}

	if _, err := os.Stat(base.BackupDir); err != nil {
		err := os.MkdirAll(base.BackupDir, 0755)
		if err != nil {
			logger.Error(err)
		}
	}

	base.Backup()
}

// restore
func Restore(base BaseModel) error {
	logger.Info("Restore starting")
	errList := []error{}
	switch base.DbInfo["dbType"] {
	case "redis":
		rdb := redisObj(base, base.DbInfo["db"])
		if err := rdb.RestoreJson(base.DbInfo["src"]); err != nil {
			errList = append(errList, err)
		}
	case "mysql":
		mysql := mysqlObj(base, base.DbInfo["db"])
		if err := mysql.Restore(base.DbInfo["src"]); err != nil {
			errList = append(errList, err)
		}
	case "es":
		es := esObj(base, base.DbInfo["db"])
		if err := es.Restore(base.DbInfo["src"]); err != nil {
			errList = append(errList, err)
		}
	case "postgresql":
		postgresql := postgresqlObj(base, base.DbInfo["db"])
		if err := postgresql.Restore(base.DbInfo["src"]); err != nil {
			errList = append(errList, err)
		}
	case "mongodb":
		mongodb := mongodbObj(base, base.DbInfo["db"])
		if err := mongodb.Restore(base.DbInfo["src"]); err != nil {
			errList = append(errList, err)
		}
	case "etcd":
		etcd := etcdObj(base)
		if err := etcd.Restore(base.DbInfo["src"], base.DbInfo["dataDir"]); err != nil {
			errList = append(errList, err)
		}
	default:
		logger.Info("no support db: ", base.DbInfo["dbType"])
		return nil
	}

	if len(errList) != 0 {
		logger.Error(errList)
		logger.Error("Restore error")
		return nil
	} else {
		logger.Info("Restore done")
	}
	return nil

}

func mysqlObj(ctx BaseModel, db string) Mysql {
	return Mysql{
		TarFilename: ctx.TarFilename,
		SaveDir:     ctx.SaveDir,
		BackupDir:   ctx.BackupDir,
		Name:        ctx.DbInfo["name"],
		Host:        ctx.DbInfo["host"],
		Port:        ctx.DbInfo["port"],
		Username:    ctx.DbInfo["username"],
		Password:    ctx.DbInfo["password"],
		Database:    db,
	}
}
func etcdObj(ctx BaseModel) Etcd {
	return Etcd{
		TarFilename: ctx.TarFilename,
		SaveDir:     ctx.SaveDir,
		BackupDir:   ctx.BackupDir,
		Name:        ctx.DbInfo["name"],
		Host:        ctx.DbInfo["host"],
		Port:        ctx.DbInfo["port"],
		Username:    ctx.DbInfo["username"],
		Password:    ctx.DbInfo["password"],
		Https:       ctx.DbInfo["https"],
		Cacert:      ctx.DbInfo["cacert"],
		Cert:        ctx.DbInfo["cert"],
		Key:         ctx.DbInfo["key"],
	}
}

func esObj(ctx BaseModel, index string) Elasticsearch {
	return Elasticsearch{
		TarFilename: ctx.TarFilename,
		SaveDir:     ctx.SaveDir,
		BackupDir:   ctx.BackupDir,
		Name:        ctx.DbInfo["name"],
		Host:        ctx.DbInfo["host"],
		Port:        ctx.DbInfo["port"],
		Username:    ctx.DbInfo["username"],
		Password:    ctx.DbInfo["password"],
		Index:       index,
	}
}

func mongodbObj(ctx BaseModel, db string) Mongodb {
	return Mongodb{
		TarFilename: ctx.TarFilename,
		SaveDir:     ctx.SaveDir,
		BackupDir:   ctx.BackupDir,
		Name:        ctx.DbInfo["name"],
		Host:        ctx.DbInfo["host"],
		Port:        ctx.DbInfo["port"],
		Username:    ctx.DbInfo["username"],
		Password:    ctx.DbInfo["password"],
		AuthDb:      ctx.DbInfo["authdb"],
		Database:    db,
	}
}

func postgresqlObj(ctx BaseModel, db string) Postgresql {
	return Postgresql{
		TarFilename: ctx.TarFilename,
		SaveDir:     ctx.SaveDir,
		BackupDir:   ctx.BackupDir,
		Name:        ctx.DbInfo["name"],
		Host:        ctx.DbInfo["host"],
		Port:        ctx.DbInfo["port"],
		Username:    ctx.DbInfo["username"],
		Password:    ctx.DbInfo["password"],
		Database:    db,
	}
}

func redisObj(ctx BaseModel, db string) Redis {
	return Redis{
		TarFilename: ctx.TarFilename,
		SaveDir:     ctx.SaveDir,
		BackupDir:   ctx.BackupDir,
		Name:        ctx.DbInfo["name"],
		Host:        ctx.DbInfo["host"],
		Port:        ctx.DbInfo["port"],
		Password:    ctx.DbInfo["password"],
		Model:       ctx.DbInfo["model"],
		Database:    db,
	}
}
func fileObj(ctx BaseModel) File {
	return File{
		TarFilename: ctx.TarFilename,
		SaveDir:     ctx.SaveDir,
		BackupDir:   ctx.BackupDir,
		Name:        ctx.DbInfo["name"],
		Host:        ctx.DbInfo["host"],
		Port:        ctx.DbInfo["port"],
		Username:    ctx.DbInfo["username"],
		Password:    ctx.DbInfo["password"],
		Path:        ctx.DbInfo["path"],
	}
}

func choiceType(ctx BaseModel) []error {
	var errList []error
	databaseList := strings.Split(ctx.DbInfo["database"], `,`)
	switch ctx.DbInfo["type"] {
	case "mysql":
		for _, db := range databaseList {
			mysql := mysqlObj(ctx, db)
			if err := mysql.Backup(); err != nil {
				errList = append(errList, err)
			}
		}
	case "etcd":
		etcd := etcdObj(ctx)
		if err := etcd.Backup(); err != nil {
			errList = append(errList, err)
		}
	case "es":
		indexList := strings.Split(ctx.DbInfo["index"], `,`)
		for _, index := range indexList {
			es := esObj(ctx, index)
			if err := es.Backup(); err != nil {
				errList = append(errList, err)
			}
		}
	case "mongodb":
		for _, db := range databaseList {
			mongodb := mongodbObj(ctx, db)
			if err := mongodb.Backup(); err != nil {
				errList = append(errList, err)
			}
		}
	case "postgresql":
		for _, db := range databaseList {
			postgresql := postgresqlObj(ctx, db)
			if err := postgresql.Backup(); err != nil {
				errList = append(errList, err)
			}
		}
	case "redis":
		for _, db := range databaseList {
			redis := redisObj(ctx, db)
			if err := redis.Backup(); err != nil {
				errList = append(errList, err)
			}
		}
	case "file":
		file := fileObj(ctx)
		if err := file.Backup(); err != nil {
			errList = append(errList, err)
		}
	default:
		logger.Info("no support ", ctx.DbInfo["name"])
	}
	return errList
}

// Backup
func (ctx BaseModel) Backup() {
	logger.Info("starting backup: ", ctx.DbInfo["name"])
	errList := choiceType(ctx)
	if len(errList) != 0 {
		logger.Error(errList)
		logger.Error("backup error", ctx.DbInfo["name"])
		return
	} else {
		logger.Info("backup done", ctx.DbInfo["name"])
	}

	BackupDirName := strings.Split(ctx.BackupDir, `/`)
	archive.ArchiveTar(
		true, ctx.SaveDir,
		BackupDirName[len(BackupDirName)-1],
		ctx.TarFilename,
	)

	cleanHistoryFile(
		ctx.SaveDir,
		fmt.Sprintf("%v-*gz", ctx.DbInfo["name"]),
		ctx.BackupNum,
	)
	if ctx.SaveInfo["type"] == "sftp" {
		sftp := new(ssh.ClientConfig)
		sshPort, _ := strconv.ParseInt(ctx.SaveInfo["port"], 10, 64)
		sftp.CreateClient(ctx.SaveInfo["host"], sshPort, ctx.SaveInfo["username"], ctx.SaveInfo["password"])
		if err := sftp.Upload(ctx.TarFilename, fmt.Sprintf("%v/%v", ctx.SaveInfo["dstpath"], ctx.TarName)); err != nil {
			logger.Info("put sftp fail")
		} else {
			logger.Info("put sftp success")
		}

	} else if ctx.SaveInfo["type"] == "ftp" {
		ftp := ftpclient.FtpClient{
			Host:     ctx.SaveInfo["host"],
			Port:     ctx.SaveInfo["port"],
			Username: ctx.SaveInfo["username"],
			Password: ctx.SaveInfo["password"],
		}
		if err := ftp.Upload(ctx.TarFilename, ctx.SaveInfo["dstpath"], ctx.TarName); err != nil {
			logger.Info("put ftp fail")
		} else {
			logger.Info("put ftp success")
		}
	}
}
