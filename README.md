# Deploy Kit 

tags： tool

---
DeployKit，部署工具,一个命令形式的可执行程序`dk`。在实际应用上线是再平常不过的事了，目前的一些平台工具可以覆盖大部分的上线工作，可是还是有一丢丢需要手动操作，这里就是来解决这部分问题的。

## 起步
这些任务概括起来是把一个包部署到一组服务器，然后执行一些脚本，具体过程类似：1.上传，2 解压。文件包最可能在的几个地方：

1. 本地磁盘。需要上传服务器.并执行相应命令。
2. 公网地址。生产服务器可以访问，直接在服务器上 `wget`[wget for win][1]就可以获得。
3. 内网地址。生产服务器不能访问，需要先下载到本地磁盘再上传到服务器。

工具使用命令风格，没有界面。编译结果为可执行文件`dk` 其参数用法如下:

```sh
Usage of E:\github\deployKit\bin\dk.exe:
  -cmd string
        后置命令,文件上传成功后在server workDir 中执行的命令，多条以分号隔开。可以使用变量{version}。
  -lurl string
        内网仓库地址.服务器不能直接访问,需要先下载到本地磁盘再上传服务器. e.g.http://127.0.0.1/{version}/a.zip.
  -name string
        项目名称，对应配置文件名,默认是config. e.g. ec 代表使用配置文件ec.json.
  -path string
        本地磁盘路径，直接上传服务器. e.g. /tmp/{version}/a.zip.
  -url string
        外网仓库地址.可以直接在服务器上 wget. e.g. http://test.com/{version}/a.zip.
  -v string
        版本，是zip文件名的一部分. e.g. v1.0.0.

注意: url,path,和lurl三个参数互斥,按照上述顺序检查到一个有效值时停止,否则报错.
```
> `./dk -h` 可查看Usage。

配置文件`config.json`和dk.exe 同目录，自定义配置文件名时需要在`-name`参数中传入。结构如下：

```json
{
  "name": "ec-web",
  "version": "v0.8.5_006",
  "url": "",
  "path": "./upload/walle-web.tar",
  "lurl": "",
  "suffixCmd": "tar -xvf walle-web.tar -C ./{version}",
  "servers": [
    {
      "ip": "172.30.10.82",
      "port": 22,
      "username": "root",
      "password": "xxxxx",
      "workDir": "/tmp/jacky"
    },
    {
      "ip": "172.30.10.83",
      "port": 22,
      "username": "root",
      "password": "xxxxx",
      "workDir": "/tmp/jacky"
    }
  ]
}
```

其中各个属性的意义：

- name对应参数 `-name` 是项目名称不设置时使用默认值"config"；
- version 对应参数 `-v` 是目标项目的版本号；
- url、path和lurl是资源文件路径对应参数 `-path -url -lurl`;
- suffixCmd 对应参数`-cmd` 是在服务器上要执行的脚本；
- servers是目标服务器信息列表，没有参数对应。

**如果参数在配置文件和命令行都有设置，优先使用命令行。**

实例中展示了如何把walle-web.tar的v0.8.5_006版本从本地磁盘"./upload/walle-web.tar"部署到82和83两个环境中，设置好配置文件后双击dk.exe即可。

### 换做用命令的方式
配置文件中只有servers是必须的，其他都是可选的。

```sh
./dk -name=config -v=v0.8.5_006 -path="./upload/walle-web.tar" -cmd='mkdir {version};tar -xvf walle-web.tar -C ./{version}'
```
如果配置文件名是默认的 config.json 则-name 参数可省略。
```sh
./dk -v=v0.8.5_006 -path="./upload/walle-web.tar" -cmd='mkdir {version};tar -xvf walle-web.tar -C ./{version}'
```
如果配置文件中设置了 `version`,`path`，`url`，`lurl`,`suffixCmd`如示例中那样，则对应的 `-v,-path,-url,-lurl,-cmd` 都可以省略,也就是直接双击运行。
```sh
./dk
```


### 特殊用法
有一种情况，文件在local repository，希望做一些处理之后再上传服务器，步骤：下载，处理，部署。


```shell
#!/usr/bin/env bash
#download
wget -P ./upload http://172.30.10.171/FacebookPMD/EC/snapshots/v0.8.5_006/walle-web.tar
#do sth
tar -zcvf ./upload/walle-web.tar.gz ./upload/walle-web.tar
#execute
./dk -name=ec -path=./upload/walle-web.tar.gz -cmd='mkdir {version};tar -zxvf walle-web.tar.g -C ./{version};rm -f walle-web.tar.gz;'

```

## 最佳实践
1. 直接双击。初次部署时在config.json配置所有配置项,以后每次部署编辑config.json中的`version`，然后直接双击dk.exe。
2. 命令形式，不用编辑配置文件。同样初次部署时在config.json配置所有配置项。以后每次使用命令传入参数`-v`覆盖配置文件`version`属性：

    ```shell
    ./dk -v=v0.8.5_006 
    ```
3. 多项目。为每个项目制作配置文件，如a.json,b.json。运行命令时指定 `-name`参数。
    ```shell
    ./dk -name=a -v=v0.2
    ./dk -name=b -v=v0.1
    ```

3. 特殊用法。从局域网中下载后在上传服务器之前希望做一些处理，把这些所有写成脚本。每次部署新版本时编辑脚本后执行。还有更好的办法是让这个脚本可以接受一个参数作为版本，使用起来就像是：
    ```sh
    deploy.sh v1.0.0
    ```
一个例子：
    ```sh
    #!/usr/bin/env bash
    version=""
    
    if [ ! $version ]; then
      read -p "please enter tag name :" tag
      version=$tag
    fi
    echo  "tag name is $version."
    
    # dwonload
    if [  -f "./upload/walle-web.tar" ] ; then
      echo "target file already exsit, redownload?(y/n)"
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
    ./dk -name=ec -v=$version -path=./upload/walle-web.tar.gz
    
    read -s -n 1 -p "Press any key to exit..."
    echo
    echo bye...
    exit 0
    
    ```
        
    
    
以上。


  [1]: https://eternallybored.org/misc/wget/