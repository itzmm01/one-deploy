数据库一键备份工具

目前只支持x86（后续支持arm64架构）

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

