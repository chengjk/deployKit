#!/usr/bin/env bash

wget -P ./upload http://172.30.10.171/FacebookPMD/EC/snapshots/v0.8.6_005/walle-web.tar
cd upload/
tar -zcvf walle-web.tar.gz walle-web.tar
cd ..
./dk.exe -name=ec -path=./upload/walle-web.tar.gz

