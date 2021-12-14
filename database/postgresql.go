package database

import (
	"fmt"
	"one-backup/cmd"
	"one-backup/keygen"
	"one-backup/tool"
	"os"
)

type Postgresql struct {
	/*
			name: postgresql
		    # 数据库类型
		    type: postgresql
		    # 需要备份的数据库，多个用英文逗号隔开, alldatabase 代表所有数据库
		    database: yc,test1
		    # 数据库IP
		    host: 127.0.0.1
		    # 端口
		    port: 5432
		    # 账号
		    username: root
		    # 密码
		    password: Amt_2018
	*/
	TarFilename string
	SaveDir     string
	BackupDir   string
	Name        string
	Host        string
	Port        string
	Username    string
	Password    string
	Db          string
}

func (ctx Postgresql) Backup() error {
	// pg_dump只支持单个数据库
	cmdStr := fmt.Sprintf("pg_dump -h %v -p %v -U %v ", ctx.Host, ctx.Port, ctx.Username)
	if ctx.Password != "" {
		os.Setenv("PGPASSWORD", ctx.Password)
	}

	cmdStr = cmdStr + fmt.Sprintf("-d %v -f %v/%v.sql", ctx.Db, ctx.BackupDir, ctx.Db)

	if err := cmd.Run(cmdStr, Debug); err == nil {
		keygen.AesEncryptCBCFile(fmt.Sprintf("%v/%v.sql", ctx.BackupDir, ctx.Db), fmt.Sprintf("%v/%v-Encrypt.sql", ctx.BackupDir, ctx.Db))
		return cmd.Run(fmt.Sprintf("rm -f %v/%v.sql", ctx.BackupDir, ctx.Db), false)
	} else {
		return err
	}
}

func (ctx Postgresql) Restore(filePath string) error {
	dstPath := "/tmp/" + tool.RandomString(30)
	keygen.AesDecryptCBCFile(filePath, dstPath)

	cmdStr := fmt.Sprintf("psql  -h %v -p %v -U %v ", ctx.Host, ctx.Port, ctx.Username)
	cmdStrCreate := fmt.Sprintf("num=`%v -c '\\l'|grep %v|wc -l`; if [ $num -eq 0 ]; then %v -c 'CREATE DATABASE %v'; fi", cmdStr, ctx.Db, cmdStr, ctx.Db)
	cmdStrRestore := cmdStr + fmt.Sprintf("-d %v -f %v", ctx.Db, dstPath)

	if ctx.Password != "" {
		os.Setenv("PGPASSWORD", ctx.Password)
	}

	if err := cmd.Run(cmdStrCreate, true); err != nil {
		return err
	}

	if err := cmd.Run(cmdStrRestore, true); err != nil {
		return err
	} else {
		return cmd.Run("rm -f "+dstPath, false)
	}
}
