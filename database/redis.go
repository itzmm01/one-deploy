package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"one-backup/cmd"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
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
	Databases   string
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
		cmd.Run(cmdStr + " save")
		cmdStr = cmdStr + fmt.Sprintf("--rdb %v/dump.rdb", ctx.BackupDir)
		return cmd.Run(cmdStr)
	} else {

		if err := BackupJson(ctx); err != nil {
			println(err.Error())
			return err
		}
		return nil
	}

}
func (ctx Redis) Restore(src string) error {
	jsonFile, err := os.Open(src)

	// 最好要处理以下错误
	if err != nil {
		return err
	}

	// 要记得关闭
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	allKeys := AllKey{}
	json.Unmarshal(byteValue, &allKeys)

	ctx1 := context.Background()
	r := Redis{
		Host:      ctx.Host,
		Port:      ctx.Port,
		Password:  ctx.Password,
		Databases: ctx.Databases,
	}

	if err := initRedis(r); err != nil {
		fmt.Println("init redis failed err :", err)
		return err
	}

	for _, key := range allKeys.StringKey {

		if err := rdb.Set(ctx1, key.Key, key.Val, 0).Err(); err != nil {
			return err
		}

	}
	for _, key := range allKeys.ListKey {

		if err := rdb.RPush(ctx1, key.Key, key.ListVal).Err(); err != nil {
			return err
		}

	}
	for _, key := range allKeys.HashKey {

		if err := rdb.HSet(ctx1, key.Key, key.HashVal).Err(); err != nil {
			return err
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

	return nil
}

func initRedis(r Redis) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	db, _ := strconv.Atoi(string(r.Databases))
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", r.Host, r.Port),
		Password: r.Password,
		DB:       db,
		PoolSize: 200,
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
	var keysList []string

	allKeys := AllKey{}
	for {
		var err error
		var keys []string
		//*扫描所有key，每次20条
		keys, cursor, err = rdb.Scan(ctx, cursor, "*", 10).Result()
		if err != nil {
			panic(err)
		}
		keysList = append(keysList, keys...)
		if cursor == 0 {
			break
		}
	}

	for _, key := range keysList {
		curType := Result{}
		sType, err := rdb.Type(ctx, key).Result()
		if err != nil {
			return err
		}

		expire, err := rdb.TTL(ctx, key).Result()
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
	if distFile, err := os.OpenFile(fmt.Sprintf("%v/dump.json", r.BackupDir), os.O_CREATE, 0666); err != nil {
		return err
	} else {
		enc := json.NewEncoder(distFile)
		if err := enc.Encode(allKeys); err != nil {
			return err
		}
	}
	return nil
}
