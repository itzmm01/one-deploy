## 数据库一键备份工具



## 介绍

 基于`go 17.3`开发封装各种数据库各自的备份工具，简化备份操作，支持数据库备份和配置文件备份。

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
-file： 指定备份的配置文件
```

## 编译

```bash
# 需要go环境,拉取代码后执行,会输出one-backup.tar.gz文件
./build.sh

```

## 安装

```bash
# 将one-backup.tar.gz上传服务器解压即可
tar xf one-backup.tar.gz
```

## 备份

请根据`example.yml`修改自己需要备份的相关信息

```bash
cd one-backup
chmod +x ./one-backup
cp example.yml xx.yml
# 按照example.yml示例修改配置文件
./one-backup -file xx.yml
```

### 使用示例

#### mysql

通过配置文件备份

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

#### postgresql

通过配置文件备份

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

#### redis

通过配置文件备份

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
    mode: rdb
    # 数据库IP
    host: 192.168.146.134
    # 端口
    port: 6379
    # 密码
    password: Amt_2018
    # 数据库
    db: 0
```

#### mongo

通过配置文件备份

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

#### es

通过配置文件备份

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

#### etcd

通过配置文件备份

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
    port: 49514
    # 是否使用https, yes|no
    https: no
    # ca证书路径
    cacert: /etc/etcd/ssl/ca.pem
    # 客户端证书路径
    cert: /etc/etcd/ssl/etcd.pem
    # 客户端密钥路径
    key: /etc/etcd/ssl/etcd-key.pem
```

#### 配置文件

通过配置文件备份

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

## 恢复

暂时只支持redis-json格式的恢复

### redis

```bash
./one-backup -mode restore -type redis -host 192.168.146.134 -port 6380 -password xxx -db 0 -src "./dump.json"
```

### mysql

```bash
./one-backup -mode restore -type mysql -host 192.168.146.134 -port 3316 -username root -password xxx -src /tmp/backupdir/mysql/mysql/yc-Encrypt.sql
```



## 定期备份

目前脚本不支持周期备份，可与crontab配合进行定时备份
