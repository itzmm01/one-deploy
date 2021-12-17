package database

import (
	"fmt"
	"log"
	"one-backup/archive"
	"one-backup/config"
	"one-backup/keygen"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/wonderivan/logger"
)

const Debug = false

// base info
type BaseModel struct {
	TarFilename string
	SaveDir     string
	BackupDir   string
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

	nameDir := fmt.Sprintf("%v-%v", dbinfo["type"], time.Now().Format("2006.01.02.15.04.05"))
	base := BaseModel{
		TarFilename: fmt.Sprintf("%v/%v/%v.tar.gz", configInfo.StoreWith["path"], dbinfo["type"], nameDir),
		SaveDir:     fmt.Sprintf("%v/%v/", configInfo.StoreWith["path"], dbinfo["type"]),
		BackupDir:   fmt.Sprintf("%v/%v/%v", configInfo.StoreWith["path"], dbinfo["type"], nameDir),
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
func Restore(dbType, host, port, username, password, db, src string) error {
	logger.Info("Restore starting")
	switch dbType {
	case "redis":
		rdb := Redis{
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
			Database: db,
		}
		if err := rdb.RestoreJson(src); err != nil {
			logger.Error(err)
			return err
		}
	case "mysql":
		mysql := Mysql{
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
			Database: db,
		}
		if err := mysql.Restore(src); err != nil {
			logger.Error(err)
			return err
		}
	case "es":
		es := Elasticsearch{
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
		}
		if err := es.Restore(src); err != nil {
			logger.Error(err)
			return err
		}
	case "postgresql":
		postgresql := Postgresql{
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
			Database: db,
		}
		if err := postgresql.Restore(src); err != nil {
			logger.Error(err)
			return err
		}
	default:
		logger.Info("no support db: ", dbType)
		return nil
	}

	logger.Info("Restore success")
	return nil

}

// Backup
func (ctx BaseModel) Backup() {
	var errList []error

	databaseList := strings.Split(ctx.DbInfo["database"], `,`)
	logger.Info("starting backup: ", ctx.DbInfo["name"])

	switch ctx.DbInfo["type"] {
	case "mysql":
		for _, db := range databaseList {
			mysql := Mysql{
				TarFilename: ctx.TarFilename,
				SaveDir:     ctx.SaveDir,
				BackupDir:   ctx.BackupDir,
				Host:        ctx.DbInfo["host"],
				Port:        ctx.DbInfo["port"],
				Username:    ctx.DbInfo["username"],
				Password:    ctx.DbInfo["password"],
				Database:    db,
			}
			if err := mysql.Backup(); err != nil {
				errList = append(errList, err)
			}
		}
	case "etcd":
		etcd := Etcd{
			TarFilename: ctx.TarFilename,
			SaveDir:     ctx.SaveDir,
			BackupDir:   ctx.BackupDir,
			Host:        ctx.DbInfo["host"],
			Port:        ctx.DbInfo["port"],
			Username:    ctx.DbInfo["username"],
			Password:    ctx.DbInfo["password"],
			Https:       ctx.DbInfo["https"],
			Cacert:      ctx.DbInfo["cacert"],
			Cert:        ctx.DbInfo["cert"],
			Key:         ctx.DbInfo["key"],
		}
		if err := etcd.Backup(); err != nil {
			errList = append(errList, err)
		}
	case "es":
		indexList := strings.Split(ctx.DbInfo["index"], `,`)
		for _, index := range indexList {
			es := Elasticsearch{
				TarFilename: ctx.TarFilename,
				SaveDir:     ctx.SaveDir,
				BackupDir:   ctx.BackupDir,
				Host:        ctx.DbInfo["host"],
				Port:        ctx.DbInfo["port"],
				Username:    ctx.DbInfo["username"],
				Password:    ctx.DbInfo["password"],
				Index:       index,
			}
			if err := es.Backup(); err != nil {
				errList = append(errList, err)
			}
		}

	case "mongodb":
		for _, db := range databaseList {
			mongodb := Mongodb{
				TarFilename: ctx.TarFilename,
				SaveDir:     ctx.SaveDir,
				BackupDir:   ctx.BackupDir,
				Host:        ctx.DbInfo["host"],
				Port:        ctx.DbInfo["port"],
				Username:    ctx.DbInfo["username"],
				Password:    ctx.DbInfo["password"],
				AuthDb:      ctx.DbInfo["authdb"],
				Database:    db,
			}
			if err := mongodb.Backup(); err != nil {
				errList = append(errList, err)
			}
		}
	case "postgresql":
		for _, db := range databaseList {
			postgresql := Postgresql{
				TarFilename: ctx.TarFilename,
				SaveDir:     ctx.SaveDir,
				BackupDir:   ctx.BackupDir,
				Host:        ctx.DbInfo["host"],
				Port:        ctx.DbInfo["port"],
				Username:    ctx.DbInfo["username"],
				Password:    ctx.DbInfo["password"],
				Database:    db,
			}
			if err := postgresql.Backup(); err != nil {
				errList = append(errList, err)
			}
		}
	case "redis":
		for _, db := range databaseList {
			redis := Redis{
				TarFilename: ctx.TarFilename,
				SaveDir:     ctx.SaveDir,
				BackupDir:   ctx.BackupDir,
				Host:        ctx.DbInfo["host"],
				Port:        ctx.DbInfo["port"],
				Password:    ctx.DbInfo["password"],
				Model:       ctx.DbInfo["model"],
				Database:    db,
			}

			if err := redis.Backup(); err != nil {
				errList = append(errList, err)
			}
		}
	case "file":
		file := File{
			TarFilename: ctx.TarFilename,
			SaveDir:     ctx.SaveDir,
			BackupDir:   ctx.BackupDir,
			Host:        ctx.DbInfo["host"],
			Port:        ctx.DbInfo["port"],
			Username:    ctx.DbInfo["username"],
			Password:    ctx.DbInfo["password"],
			Path:        ctx.DbInfo["path"],
		}
		if err := file.Backup(); err != nil {
			errList = append(errList, err)
		}
	default:
		logger.Info("no support ", ctx.DbInfo["name"])
	}

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
}
