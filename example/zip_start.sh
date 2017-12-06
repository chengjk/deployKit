#!/usr/bin/env bash

tag=v0.8.8

if [ ! $tag ]; then
	read -p "please enter tag name :" tag
	tag=$tag
fi
echo  "tag name is $tag."

if [  -f "./upload/walle-web.tar" ] ; then
	echo "target file already exsit, redownload?(y/n)"
	read answer
	if [ "$answer" == "y" ]; then
	    rm -rf ./upload/*
        wget -P ./upload http://172.30.10.171/FacebookPMD/EC/snapshots/$tag/walle-web.tar
	fi
else
  wget -P ./upload http://172.30.10.171/FacebookPMD/EC/snapshots/$tag/walle-web.tar
fi

# do sth

cd ./upload
tar -zcvf walle-web.tar.gz walle-web.tar
cd ..

# execute
./dk -name=lan -v=$tag -path=./upload/walle-web.tar.gz

read -s -n 1 -p "Press any key to exit..."
echo
echo bye...
exit 0