package database

import (
	"fmt"
	"log"
	"one-backup/archive"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/wonderivan/logger"
)

const Debug = false

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
				Db:          db,
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
				Db:          db,
			}
			if err := postgresql.Backup(); err != nil {
				errList = append(errList, err)
			}
		}
	case "redis":
		redis := Redis{
			TarFilename: ctx.TarFilename,
			SaveDir:     ctx.SaveDir,
			BackupDir:   ctx.BackupDir,
			Host:        ctx.DbInfo["host"],
			Port:        ctx.DbInfo["port"],
			Password:    ctx.DbInfo["password"],
			Model:       ctx.DbInfo["model"],
			Database:    ctx.DbInfo["db"],
		}

		if err := redis.Backup(); err != nil {
			errList = append(errList, err)
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

	archive.ArchiveTar(true, ctx.SaveDir, ctx.DbInfo["name"], ctx.TarFilename)

	cleanHistoryFile(ctx.SaveDir, fmt.Sprintf("%v-*", ctx.DbInfo["name"]), ctx.BackupNum)
}
