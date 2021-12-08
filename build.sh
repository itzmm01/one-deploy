#!/bin/bash

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

go build main.go
mkdir -p one-backup
mv main one-backup/one-backup

cp -rf bin one-backup/bin
cp -rf example.yml one-backup/
chmod +x one-backup/bin/*
tar zcf one-backup.tar.gz one-backup/ --remove-files