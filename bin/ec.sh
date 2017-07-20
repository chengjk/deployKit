#!/usr/bin/env bash

version=v0.8.6_005

if [! -f "./upload/walle-web.tar"]  then
	wget -P ./upload http://172.30.10.171/FacebookPMD/EC/snapshots/$version/walle-web.tar
else
	echo "文件已经有了，不再重复下载"
fi


cd upload/
tar -zcvf walle-web.tar.gz walle-web.tar
cd ..
./dk.exe -name=ec -path=./upload/walle-web.tar.gz

