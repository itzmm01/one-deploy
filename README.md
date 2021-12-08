## 数据库一键备份工具

基于`go 17.3`开发封装各种数据库各自的备份工具，简化备份操作，

目前只支持x86（后续支持arm64架构）

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

以及配置文件备份

编译打包

```bash
./build.sh
```

输出 `one-backup.tar.gz`

安装

```bash
# 解压即可
tar xf one-backup.tar.gz
```

### 备份

请根据`example.yml`修改自己需要备份的相关信息

```bash
cd one-backup
chmod +x ./one-backup
cp example.yml db.yml
# 按照example.yml示例修改配置文件
./one-backup -f db.yml
```

### 恢复

暂时只支持redis-json格式的恢复

`redis`

```bash
./one-backup -restore yes -type redis -host 192.168.146.134 -port 6380 -password xxx -db 0 -src "./dump.json"
```

