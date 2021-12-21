## 介绍

 基于`go 17.3`开发封装各种数据库各自的备份工具，简化备份操作，支持数据库备份和配置文件备份,主要解决少量元数据备份场景。

支持数据库有

>mysql
>
>redis
>
>mongodb
>
>pgsql
>
>etcd
>
>es

## 参数

```bash
# 是否自动加密密码,默认加密
  -autoEncrypt string
    	yes|no (default "yes")
# 指定需要加密的字符串,输出加密后的字符串
  -encrypt string
    	need encrypt string
# 工具运行模式,备份或者恢复，默认是恢复
  -mode string
    	run mode: backup|restore (default "backup")
#---------备份时使用参数
# 指定配置文件
  -file string
    	config file (default "./backupdb.yml")
#---------恢复时使用参数
# 数据库类型
  -type string
    	database type: redis|mysql|mongodb|etcd|es|postgresql    	
# 数据库IP
  -host string
    	database host: x.x.x.x
# 数据库端口
  -port string
    	database port: 6379
# 数据库用户
  -username string
    	database username: root
# 数据库密码
  -password string
    	database password: xxx
# 数据库名称
  -db string
    	database: 0 (default "0")
# 恢复来源文件
  -src string
    	restore dir/file:  such '/tmp/backupdir/redis/dump.json' 

```

## 编译

```bash
# 需要go环境,拉取代码后执行,会输出one-backup-linux-架构.tar.gz文件(目前适配uos-arm,centos-x86)
./build.sh
```

## 安装使用

```bash
# 将one-backup-linux-架构.tar.gz上传服务器解压即可
tar xf one-backup-linux-架构.tar.gz
```

请根据`example.yml`修改自己需要备份的相关信息

```bash
cd one-backup
chmod +x ./one-backup
cp example.yml xx.yml
# 按照example.yml示例修改配置文件
./one-backup -file xx.yml
```

## mysql

### 备份

> ./one-backup -file mysql.yml

```yaml
# 保留备份文件数量
backupnum: 5
# 压缩选项
compresstype: tgz
# 存储选项
storewith:
  # 类型(暂时支持本地存储)
  type: local
  # 路径
  path: /tmp/backupdir
databases:
    # 备份任务名
  - name: mysql
    # 数据库类型
    type: mysql
    # 数据库IP
    host: 192.168.146.134
    # 端口
    port: 3306
    # 账号
    username: root
    # 密码
    password: Amt_2018
    # 需要备份的数据库，多个用英文逗号隔开, alldatabase 代表所有数据库
    database: test1,yc
```

### 恢复

```bash
./one-backup -mode restore -type mysql -host 192.168.146.134 -port 3316 -username root -password xxx -src /tmp/backupdir/mysql/mysql/yc-Encrypt.sql
```

## postgresql

### 备份

> ./one-backup  -file  postgresql.yml

```yaml
# 保留备份文件数量
backupnum: 5
# 压缩选项
compresstype: tgz
# 存储选项
storewith:
  # 类型(暂时支持本地存储)
  type: local
  # 路径
  path: /tmp/backupdir
databases:
    # 备份任务名
  - name: postgresql
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
```

### 恢复

```bash
./one-backup -mode restore -type postgresql -host 192.168.146.134 -port 5432 -username root -password xxx -db test1 -src /tmp/backupdir/postgresql/postgresql-2021.12.14.12.39.32/test1-Encrypt.sql
```

## redis

### 备份

> ./one-backup  -file  redis.yml

```yaml
# 保留备份文件数量
backupnum: 5
# 压缩选项
compresstype: tgz
# 存储选项
storewith:
  # 类型(暂时支持本地存储)
  type: local
  # 路径
  path: /tmp/backupdir
databases:
    # 备份任务名
  - name: redis
    # 数据库类型
    type: redis
    # 备份方式目前支持rdb,json
    mode: json
    # 数据库IP
    host: 192.168.146.134
    # 端口
    port: 6379
    # 密码
    password: Amt_2018
    # 需要备份的数据库，多个用英文逗号隔开(json格式有效)
    database: 0
```

### 恢复

只支持redis-json模式备份数据的恢复

```bash
./one-backup -mode restore -type redis -host 192.168.146.134 -port 6380 -password xxx -db 0 -src "./dump.json"
```

### 性能数据

服务器配置: 16C,32G,150w个key

>    cpu 30%以下
>
>    内存 0.2%以下
>
>    redis链接数20
>
>    恢复 5分钟 
>
>    备份 14分钟

## mongo

### 备份

> ./one-backup  -file  mongo.yml

```yaml
# 保留备份文件数量
backupnum: 5
# 压缩选项
compresstype: tgz
# 存储选项
storewith:
  # 类型(暂时支持本地存储)
  type: local
  # 路径
  path: /tmp/backupdir
databases:
    # 备份任务名
  - name: mongo
    # 数据库类型
    type: mongodb
    # 数据库IP
    host: 192.168.146.134
    # 端口
    port: 27017
    # 账号
    username: root
    # 密码
    password: Amt_2018
    # 需要备份的数据库，多个用英文逗号隔开, alldatabase 代表所有数据库
    database: alldatabase
    # 验证用户数据库
    authdb: "admin"
```

### 恢复

暂时智能手动使用mongorestore 恢复，后续集成

## es

### 备份

> ./one-backup  -file  es.yml

```yaml
# 保留备份文件数量
backupnum: 5
# 压缩选项
compresstype: tgz
# 存储选项
storewith:
  # 类型(暂时支持本地存储)
  type: local
  # 路径
  path: /tmp/backupdir
databases:
    # 备份任务名
  - name: es
    # 数据库类型
    type: es
    # 数据库IP
    host: 192.168.146.134
    # 数据库端口
    port: 9200
    # 账号
    username: root
    # 密码
    password: Amt_2018
    # 备份的索引,多个用英文逗号隔开
    index: abc
```

### 恢复



```bash
./one-backup -mode restore -type es -host 192.168.146.134 -port 9200 -src /tmp/backupdir/es/es-2021.12.14.11.24.08/test111-Encrypt.json
```

## etcd

### 备份

> ./one-backup  -file  etcd.yml

```yaml
# 保留备份文件数量
backupnum: 5
# 压缩选项
compresstype: tgz
# 存储选项
storewith:
  # 类型(暂时支持本地存储)
  type: local
  # 路径
  path: /tmp/backupdir
databases:
    # 备份任务名
  - name: etcd
    # 数据库类型
    type: etcd
    # etcd链接信息
    host: 127.0.0.1
    # 端口
    port: 49514
    # 用户
    username: root
    # 密码
    password: Amt_2018
    # 是否使用https, yes|no
    https: no
    # ca证书路径
    cacert: /etc/etcd/ssl/ca.pem
    # 客户端证书路径
    cert: /etc/etcd/ssl/etcd.pem
    # 客户端密钥路径
    key: /etc/etcd/ssl/etcd-key.pem
```

### 恢复

手动使用etcdctl恢复，后续集成

## zookeeper

### 备份

> ./one-backup  -file  config.yml

```yaml
# 保留备份文件数量
backupnum: 5
# 压缩选项
compresstype: tgz
# 存储选项
storewith:
  # 类型(暂时支持本地存储)
  type: local
  # 路径
  path: /tmp/backupdir
# 配置文件备份
databases:
    # 备份名
  - name: "zookeeper"
    # 类型 file
    type: file
    # zookeeper的dataDir目录,注意需要/结尾
    path: "/data/hadoop/zookeeper/"
    # 需要备份的服务器,local代表本机且无需账号密码
    host: "x.x.x.x"
    # ssh端口
    port: "22"
    # ssh 账号
    username: "root"
    # ssh密码
    password: "xxxx"
```

### 恢复

手动cp恢复

## HDFS-fsimage

### 备份

> ./one-backup  -file  config.yml

```yaml
# 保留备份文件数量
backupnum: 5
# 压缩选项
compresstype: tgz
# 存储选项
storewith:
  # 类型(暂时支持本地存储)
  type: local
  # 路径
  path: /tmp/backupdir
# 配置文件备份
databases:
    # 备份名
  - name: "zookeeper"
    # 类型 file
    type: file
    # HDFS-fsimage的路径
    path: "/data/hadoop/fsimage-xxx"
    # 需要备份的服务器,local代表本机且无需账号密码
    host: "x.x.x.x"
    # ssh端口
    port: "22"
    # ssh 账号
    username: "root"
    # ssh密码
    password: "xxxx"
```

### 恢复

手动cp恢复

## 文件

### 备份

> ./one-backup  -file  config.yml

```yaml
# 保留备份文件数量
backupnum: 5
# 压缩选项
compresstype: tgz
# 存储选项
storewith:
  # 类型(暂时支持本地存储)
  type: local
  # 路径
  path: /tmp/backupdir
# 配置文件备份
databases:
    # 备份名
  - name: "mysql_cnf"
    # 类型 配置文件填file
    type: file
    # 备份的配置文件路径
    path: "/etc/my.cnf"
    # 需要备份的服务器,local代表本机且无需账号密码
    host: "x.x.x.x"
    # ssh端口
    port: "22"
    # ssh 账号
    username: "root"
    # ssh密码
    password: "xxxx"
```

### 恢复

手动cp恢复



## 定期备份

目前脚本不支持周期备份，可与crontab配合进行定时备份
