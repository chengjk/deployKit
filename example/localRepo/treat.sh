#!/usr/bin/env bash

version=v0.8.6_005

if [  -f "./upload/walle-web.tar" ] ; then
	echo "文件已存在，是否要重新下载？(y/n)"
	read answer
	if [ "$answer" == "y" ]; then
	    rm -rf ./upload/*
        wget -P ./upload http://172.30.10.171/FacebookPMD/EC/snapshots/$version/walle-web.tar
	fi
else
  wget -P ./upload http://172.30.10.171/FacebookPMD/EC/snapshots/$version/walle-web.tar
fi

# do sth
tar -zcvf ./upload/walle-web.tar.gz ./upload/walle-web.tar

# execute
./dk.exe -name=treat -path=./upload/walle-web.tar.gz

