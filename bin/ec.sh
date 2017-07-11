#!/usr/bin/env bash

wget -P ./upload http://172.30.10.171/FacebookPMD/EC/snapshots/v0.8.5_006/walle-web.tar
cd upload/
tar -zcvf walle-web.tar.gz walle-web.tar
cd ..
./dk.exe -name=ec -path=./upload/walle-web.tar.gz


## 后置命令
tar -zxvf walle-web.tar.gz
tar -xvf walle-web.tar
mkdir v0.8.5_006
tar -xvf ./walle-web/walle-web.tar -C ./v0.8.5_006