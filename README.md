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
>
>zookeeper
>
>文件

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
  
#----------mongodb
  -authdb string
        mongo authdb: admin (default "admin")
#----------etcd
  -Key string
        etcd key (default "/etc/kubernetes/pki/etcd/server.key")
  -cacert string
        etcd cacert (default "/etc/kubernetes/pki/etcd/ca.crt")
  -cert string
        etcd cert (default "/etc/kubernetes/pki/etcd/server.crt")
  -datadir string
        etcd data-dir (default "/var/lib/etcd")
  -https string
        etcd https (default "no")


```



## 编译

```bash
# 需要go环境,拉取代码后执行,会输出one-backup-linux-架构.tar.gz文件(目前适配uos-arm,centos-x86)
./build.sh
# 或者直接从cos下载
# amd64
wget https://tencent-cloud-product-release-1258877907.cos.ap-guangzhou.myqcloud.com/commons/one-backup-linux-amd64.tar.gz
# arm64
wget https://tencent-cloud-product-release-1258877907.cos.ap-guangzhou.myqcloud.com/commons/one-backup-linux-arm64.tar.gz
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

## 备份选项

### 本地存储

```yaml
#存储选项
storewith:
  # 类型(暂时支持本地存储和sftp)
  type: sftp
  # 本地保存路径
  path: /tmp/backupdir
```

### sftp

```yaml
#存储选项
storewith:
  # 类型(暂时支持本地存储和sftp)
  type: sftp
  # 本地保存路径
  path: /tmp/backupdir
  # 远端主机
  host: 192.168.146.134
  # 端口
  port: "22"
  # 用户名
  username: root
  # 密码
  password: BhaBUTSg3lMXHLVUkHmOfw==
  # 远端存储路径
  dstpath: /root/
```

### ftp

```yaml
#存储选项
storewith:
  # 类型(暂时支持本地存储和sftp)
  type: ftp
  # 本地保存路径
  path: /tmp/backupdir
  # 远端主机
  host: 192.168.146.134
  # 端口
  port: "21"
  # 用户名
  username: ftp
  # 密码
  password: BhaBUTSg3lMXHLVUkHmOfw==
  # 远端存储路径
  dstpath: /root/
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

```bash
./one-backup -mode restore -type mongodb -host 192.168.146.134 -port 17017 -username root -password xxx -authdb admin -src /tmp/backupdir/mongodb/mongodb-2021.12.27.01.35.24/
```

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
  - name: etcd
    # 数据库类型
    type: etcd
    # etcd连接信息, http://xxx，集群使用"etcd1=http://xxx,etcd2=http://xxx,etcd3=xxx"
    endpoints: http://127.0.0.1:2379
    # 是否使用https, yes|no
    https: no
    # 用户
    username: xxx
    # 密码
    password: xxx
    # ca证书路径,https=yes时使用
    cacert: /etc/etcd/ssl/ca.pem
    # 客户端证书路径,https=yes时使用
    cert: /etc/etcd/ssl/etcd.pem
    # 客户端密钥路径,https=yes时使用
    key: /etc/etcd/ssl/etcd-key.pem
    
```

### 恢复

`提示:`

>执行时,如果datadir目录存在会自动将其重命名在原有目录下,
>

单机

```bash
# 示例
# -etcdname http://192.168.146.134:12380 需要恢复的etcd连接信息
# -dockername etcd1 docker部署时的容器名,非容器跳过此参数
# -dockernetwork host 网络模式nat/host(默认host),非容器跳过此参数
./one-backup -mode restore -type etcd -etcdname http://192.168.146.134:12380 -datadir /var/lib/etcd -src /tmp/backupdir/etcd/etcd-2021.12.27.22.43.14/etcd.db -dockername etcd
```

本机集群,集群内所有节点都需要还原

```bash
# 示例
# -etcdcluster "etcd1=http://192.168.146.134:12380,etcd2=http://192.168.146.134:22380,etcd3=192.168.146.134:32380" etcd集群信息
# -etcdname etcd1 需要恢复的etcd名称
# -datadir /etcd1.etcd etcd数据目录
# -etcdclustertoken etcd-cluster etcd集群token
# -src /tmp/backupdir/etcd/etcd1-2022.01.18.06.06.58/etcd.db 备份的文件路径
# -dockername etcd1 docker部署时的容器名,非容器跳过此参数
# -dockernetwork host 网络模式nat/host(默认host),非容器跳过此参数
./one-backup -mode restore -type etcd -datadir /var/lib/etcd -src /tmp/backupdir/etcd/etcd-cluter-2022.01.18.16.04.05/etcd.db -etcdcluster "etcd1=http://192.168.146.134:12380,etcd2=http://192.168.146.134:22380,etcd3=192.168.146.134:32380" -etcdname etcd1 -etcdclustertoken etcd-cluster -dockername etcd1 
```

远程集群，ssh连接,集群内所有节点都需要还原

```bash
# 示例
# -etcdcluster "etcd1=http://etcd1:2380,etcd2=http://etcd2:2380,etcd3=http://etcd3:2380" etcd集群信息
# -etcdname etcd1 需要恢复的etcd名称
# -etcdclustertoken etcd-cluster etcd集群token
# -datadir /etcd1.etcd etcd数据目录
# -src /tmp/backupdir/etcd/etcd1-2022.01.18.06.06.58/etcd.db 备份的文件路径
# -dockername etcd1 docker部署时的容器名,非容器跳过此参数
# -dockernetwork host 网络模式nat/host(默认host),非容器跳过此参数
# -sshhost 192.168.146.135 ssh连接信息
# -sshport 22 ssh连接信息
# -sshuser root ssh连接信息
# -sshpassword Amt_2018ssh连接信息
./one-backup -mode restore -type etcd -datadir /var/lib/etcd -src /tmp/backupdir/etcd/etcd-cluter-2022.01.18.16.04.05/etcd.db -etcdcluster "etcd1=http://192.168.146.134:2380,etcd2=http://192.168.146.135:2380,etcd3=http://192.168.146.136:2380" -etcdname etcd1 -etcdclustertoken etcd-cluster -dockername etcd -dockernetwork host -sshhost 192.168.146.135 -sshport 22 -sshuser root -sshpassword Amt_2018
```



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

```bash
# 手动解压备份目录 并且将zookeeper配置文件中的dataDir 指向备份目录然后重启zookeeper
```

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

```bash
# 手动解压出备份文件，并将其拷贝到对应目录
```

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

```bash
# 手动解压出备份文件，并将其拷贝到对应目录
```



## 定期备份

目前脚本不支持周期备份，可与crontab配合进行定时备份
