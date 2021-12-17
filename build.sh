#!/bin/bash


build(){
	go build main.go
	mkdir -p one-backup
	mv main one-backup/one-backup
	if [ "$GOARCH" == "amd64" ]; then
	    cp -rf bin one-backup/bin
	else
		cp -rf arm_bin one-backup/bin
	fi
	cp -rf example.yml one-backup/
	cp -rf README.md one-backup/
	chmod +x one-backup/bin/*
	tar zcf one-backup-$GOOS-$GOARCH.tar.gz one-backup/ --remove-files
}

export CGO_ENABLED=0
export GOOS=linux

# amd64
export GOARCH=amd64
echo $GOOS-$GOARCH
build

# arm64
export GOARCH=arm64
echo $GOOS-$GOARCH
build
