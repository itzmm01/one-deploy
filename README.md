## 数据库一键备份工具

基于`go 17.3`开发封装各种数据库各自得备份工具，简化备份操作，

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

使用

```bash
cd one-backup
chmod +x ./one-backup
cp example.yml db.yml
# 按照example.yml示例修改配置文件
./one-backup -f db.yml
```

