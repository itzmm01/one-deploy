package database

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"one-backup/cmd"
	"one-backup/keygen"
	"one-backup/tool"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wonderivan/logger"
)

type Redis struct {
	/*
		name: redis
		# 数据库类型
		type: redis
		# 备份方式目前支持sync
		mode: sync
		# 数据库IP
		host: 192.168.146.134
		# 端口
		port: 6379
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
	Database    string
	Model       string
}
type Result struct {
	Key     string
	Val     string
	ListVal []string
	HashVal map[string]string
	ZsetVal []string
	TTL     time.Duration
}
type AllKey struct {
	StringKey []Result
	ListKey   []Result
	HashKey   []Result
	ZsetKey   []Result
	SetKey    []Result
}

var rdb *redis.Client

func (ctx Redis) Backup() error {
	//
	if ctx.Model == "sync" {
		cmdStr := fmt.Sprintf("redis-cli -h %v -p %v ", ctx.Host, ctx.Port)
		if ctx.Password != "" {
			cmdStr = cmdStr + fmt.Sprintf("-a '%v' ", ctx.Password)
		}
		cmd.Run(cmdStr+" save", Debug)
		cmdStr = cmdStr + fmt.Sprintf("--rdb %v/dump.rdb", ctx.BackupDir)
		return cmd.Run(cmdStr, Debug)
	} else {

		if err := BackupJson(ctx); err != nil {
			return err
		}
		return nil
	}

}
func (ctx Redis) RestoreJson(filepath string) error {
	dstPath := "/tmp/" + tool.RandomString(30)
	keygen.AesDecryptCBCFile(filepath, dstPath)

	f, err := os.Open(dstPath)
	if err != nil {
		return err
	}
	// 要记得关闭
	defer f.Close()

	// 对于小于4096字节的数据，它会将所有文件内容读取到缓冲区，但对于大文件，只读取4096字节
	buf := bufio.NewReader(f)
	num := 0
	ctx1 := context.Background()
	r := Redis{
		Host:     ctx.Host,
		Port:     ctx.Port,
		Password: ctx.Password,
		Database: ctx.Database,
	}
	if err := initRedis(r); err != nil {
		fmt.Println("init redis failed err :", err)
		return err
	}
	for {
		num += 1
		// 按行读取
		line, err := buf.ReadString('\n')
		byteValue := []byte(line)

		allKeys := AllKey{}
		json.Unmarshal(byteValue, &allKeys)

		for _, key := range allKeys.StringKey {

			if err := rdb.Set(ctx1, key.Key, key.Val, -1).Err(); err != nil {
				return err
			}
			if key.TTL != -1 {
				rdb.Expire(ctx1, key.Key, key.TTL)
			}

		}
		for _, key := range allKeys.ListKey {

			if err := rdb.RPush(ctx1, key.Key, key.ListVal).Err(); err != nil {
				return err
			}
			if key.TTL != -1 {
				rdb.Expire(ctx1, key.Key, key.TTL)
			}

		}
		for _, key := range allKeys.HashKey {

			if err := rdb.HSet(ctx1, key.Key, key.HashVal).Err(); err != nil {
				return err
			}
			if key.TTL != -1 {
				rdb.Expire(ctx1, key.Key, key.TTL)
			}

		}
		for _, key := range allKeys.SetKey {
			rdb.SAdd(ctx1, key.Key, key.ListVal)

		}

		for _, key := range allKeys.ZsetKey {
			ranking := []*redis.Z{}
			for index, val := range key.ZsetVal {
				a := &redis.Z{}
				a.Score = float64(index)
				a.Member = val
				ranking = append(ranking, a)

			}
			if err := rdb.ZAdd(ctx1, key.Key, ranking...).Err(); err != nil {
				return err
			}

		}
		if err != nil {
			if err == io.EOF {
				f.Close()
				break
			}
			f.Close()
			return err
		}
	}

	return os.Remove(dstPath)
}

func initRedis(r Redis) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3600)
	defer cancel()
	db, _ := strconv.Atoi(string(r.Database))
	rdb = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%v:%v", r.Host, r.Port),
		Password:     r.Password,
		DB:           db,
		PoolSize:     20,
		MinIdleConns: 10,
	})
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func BackupJson(r Redis) error {
	if err := initRedis(r); err != nil {
		return err
	}
	ctx := context.Background()

	var cursor uint64

	for {
		var err error
		var keys []string

		allKeys := AllKey{}
		keys, cursor, err = rdb.Scan(ctx, cursor, "*", 10).Result()
		if err != nil {
			panic(err)
		}

		for _, key := range keys {
			curType := Result{}
			sType, err := rdb.Type(ctx, key).Result()
			if err != nil {
				return err
			}

			expire, err := rdb.TTL(ctx, key).Result()
			if err != nil {
				logger.Error(err)
			}
			curType.TTL = expire

			if sType == "string" {
				if val, err := rdb.Get(ctx, key).Result(); err != nil {
					return err
				} else {
					curType.Key = key
					curType.Val = val
					allKeys.StringKey = append(allKeys.StringKey, curType)
				}

			} else if sType == "list" {
				if val, err := rdb.LRange(ctx, key, 0, -1).Result(); err != nil {
					return err
				} else {
					curType.Key = key
					curType.ListVal = val
					allKeys.ListKey = append(allKeys.ListKey, curType)
				}
			} else if sType == "hash" {
				if val, err := rdb.HGetAll(ctx, key).Result(); err != nil {
					return err
				} else {
					curType.Key = key
					curType.HashVal = val
					allKeys.HashKey = append(allKeys.HashKey, curType)
				}
			} else if sType == "zset" {

				if val, err := rdb.ZRevRange(ctx, key, 0, -1).Result(); err != nil {
					return err
				} else {
					curType.Key = key
					curType.ZsetVal = val
					allKeys.ZsetKey = append(allKeys.ZsetKey, curType)
				}

			} else if sType == "set" {
				if val, err := rdb.SMembers(ctx, key).Result(); err != nil {
					return err
				} else {
					curType.Key = key
					curType.ListVal = val
					allKeys.SetKey = append(allKeys.SetKey, curType)
				}
			}

		}
		distFile, err := os.OpenFile(fmt.Sprintf("%v/%v.json", r.BackupDir, r.Database), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			return err
		} else {
			enc := json.NewEncoder(distFile)
			if err := enc.Encode(allKeys); err != nil {
				distFile.Close()
				return err
			}

		}
		if cursor == 0 {
			distFile.Close()
			break
		}
	}
	keygen.AesEncryptCBCFile(fmt.Sprintf("%v/%v.json", r.BackupDir, r.Database), fmt.Sprintf("%v/%v-Encrypt.json", r.BackupDir, r.Database))
	return os.Remove(fmt.Sprintf("%v/%v.json", r.BackupDir, r.Database))
}
