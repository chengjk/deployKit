#!/usr/bin/env bash

version=v0.8.8_001

if [  -f "./upload/walle-web.tar" ] ; then
	echo "already exsit, redownload？(y/n)"
	read answer
	if [ "$answer" == "y" ]; then
	    rm -rf ./upload/*
        wget -P ./upload http://172.30.10.171/FacebookPMD/EC/snapshots/$version/walle-web.tar
	fi
else
  wget -P ./upload http://172.30.10.171/FacebookPMD/EC/snapshots/$version/walle-web.tar
fi

# do sth

cd ./upload
tar -zcvf walle-web.tar.gz walle-web.tar
cd ..

# execute
./dk.exe -name=zip_deploy -path=./upload/walle-web.tar.gz

