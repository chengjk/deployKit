# Deploy Kit 

tags： tool

---

DeployKit，部署工具,一个命令形式的可执行程序`dk`。在实际应用上线是再平常不过的事了，目前的一些平台工具可以覆盖大部分的上线工作，可是还是有一丢丢需要手动操作，这里就是来解决这部分问题的。

[Release Note][1]

## 起步
这些任务概括起来是把一个包部署到一组服务器，然后执行一些脚本，具体过程类似：1.上传，2 解压。文件包最可能在的几个地方：

1. 本地磁盘。需要上传服务器.并执行相应命令。
2. 公网地址。生产服务器可以访问，直接在服务器上 `wget`[wget for win][2]就可以获得。
3. 内网地址。生产服务器不能访问，需要先下载到本地磁盘再上传到服务器。

### 安装
运行根目录`install.sh`可将项目打包成可执行程序到`./bin`目录。不同操作系统需要在各自的GO环境下安装。或者直接从[github][3]下载可执行文件。

### 使用

工具使用命令风格，没有界面。编译结果为可执行文件`dk` 其参数用法如下:

```sh
$ ./dk.exe -h
Usage of E:\github\deployKit\bin\dk.exe:
  -lurl string
        内网仓库地址,需要先下载到本地磁盘再上传服务器. e.g. http://127.0.0.1/{tag}/{name}.zip.
  -name string
        名称,可作为变量{name}使用,对应配置文件名.不能使用变量. e.g. ec表示使用ec.json. (default "config")
  -path string
        目标文件路径. e.g. /tmp/{tag}/{name}.zip.
  -pcmd string
        prefix cmd,文件上传前在server的workDir中执行，分号隔开.e.g. mkdir p
  -scmd string
        suffix cmd,文件上传后在server的workDir中执行，分号隔开.e.g. rm -f *.zip
  -tag string
        标签.一般是版本信息,可作为变量{tag}使用.不能使用变量. e.g. v1.0
  -url string
        外网仓库地址.直接在服务器上 wget. e.g. http://test.com/{tag}/{name}.zip.
  -v    show current version.

Tips: url,path,和lurl三个参数互斥,按照上述顺序检查到一个有效值时停止,否则报错.

```
> `./dk -h` 查看Usage。
`./dk -v` 查看当前版本号。

配置文件`config.json`和dk.exe 同目录，自定义配置文件名时需要在`-name`参数中传入。结构如下：

```json
{
  "name": "web",
  "tag": "v0.8",
  "url": "",
  "path": "./upload/{name}.tar",
  "lurl": "",
  "prefixCmd": "mkdir {tag}",
  "suffixCmd": "tar -xvf {name}.tar -C ./{tag}",
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

- name对应参数 `-name` 配置文件名。可作为变量`{name}`使用，默认值"config"。不能使用变量；
- tag 对应参数 `-tag` 是目标项目的版本号，可作为变量`{tag}`使用，不能使用变量；
- url、path和lurl是资源文件路径对应参数 `-path -url -lurl`;
- prefixCmd 对应参数`-pcmd` 文件上传前在workDir中执行，分号隔开；
- suffixCmd 对应参数`-scmd` 文件上传后在workDir中执行，分号隔开.
- servers是目标服务器信息列表，没有参数对应。

**如果参数在配置文件和命令行都有设置，优先使用命令行。**

实例中展示了如何把web.tar的v0.8版本从本地磁盘"./upload/web.tar"部署到82和83两个环境中.
如果web.tar在内网服务器上，则可设置 `-lurl`替代`-path` 为 

```diff
- "path": "./upload/{name}.tar",
+ "lurl": "http://localserver/{project}/{tag}/{name}.tar",
```
设置好配置文件后双击dk.exe即可。这种方式的优点是便捷，缺点是每次都需要更新 `tag`。

### 换做用命令的方式
配置文件中只有servers是必须的，其他都用命令行替代,全参形式如下:

```sh
./dk -name=config -tag=v0.8.5 -path="./upload/web.tar" -pcmd='mkdir {tag};' -scmd='tar -xvf web.tar -C ./{tag}'
```
如果配置文件名是默认的 config.json 则-name 参数可省略。
```sh
./dk -tag=v0.8.5 -path="./upload/web.tar" -pcmd='mkdir {tag};' -scmd='tar -xvf web.tar -C ./{tag}'
```
如果配置文件中设置了 `tag`,`path`，`url`，`lurl`,`suffixCmd`如示例中那样，则对应的 `-tag,-path,-url,-lurl,-cmd` 都可以省略,也就是直接双击运行。
```sh
./dk
```


### 特殊用法
有一种情况，文件在local repository，希望做一些处理之后再上传服务器，步骤：下载，处理，部署。

```shell
#!/usr/bin/env bash
#download
wget -P ./upload http://172.30.10.171/FacebookPMD/EC/snapshots/v0.8.5_006/web.tar
#do sth
tar -zcvf ./upload/web.tar.gz ./upload/web.tar
#execute
./dk -name=ec -path=./upload/web.tar.gz -scmd='mkdir {tag};tar -zxvf web.tar.g -C  ./{tag};rm -f web.tar.gz;'
```

## 最佳实践
1. 傻瓜式双击。初次部署时在config.json配置所有配置项,以后每次部署编辑config.json中的`tag`，然后直接双击dk.exe。
2. 命令形式，不用每次都编辑配置文件。同样初次部署时在config.json配置所有配置项。以后每次使用命令传入参数`-tag`覆盖配置文件`tag`属性：

    ```shell
    ./dk -tag=v0.8.6
    ```
3. 多项目。为每个项目制作配置文件，如a.json,b.json。运行命令时指定 `-name`参数。

    ```shell
    ./dk -name=a -tag=v0.2
    ./dk -name=b -tag=v0.1
    ```

3. 特殊用法。从局域网中下载后在上传服务器之前希望做一些处理，把这些所有写成脚本。更好的办法是让这个脚本可以接受一个参数作为版本，使用起来就像是：

    ```sh
    deploy.sh v1.0.0
    ```
或者运行过程中输入参数。例：

    ```sh
    #!/usr/bin/env bash
    tag=""
    
    if [ ! $tag ]; then
    	read -p "please enter tag name :" tag
    	tag=$tag
    fi
    echo  "tag name is $tag."
    
    # dwonload
    if [  -f "./upload/web.tar" ] ; then
    	echo "target file already exsit, redownload?(y/n)"
    	read answer
    	if [ "$answer" == "y" ]; then
    	    rm -rf ./upload/*
            wget -P ./upload http://172.30.10.171/FacebookPMD/EC/snapshots/$tag/web.tar
    	fi
    else
      wget -P ./upload http://172.30.10.171/FacebookPMD/EC/snapshots/$tag/web.tar
    fi
    
    # do sth
    cd ./upload
    tar -zcvf web.tar.gz web.tar
    cd ..
    
    # execute
    ./dk -name=ec -tag=$tag -path=./upload/web.tar.gz
    
    read -s -n 1 -p "Press any key to exit..."
    echo
    echo bye...
    exit 0
    
    ```
  
以上。


  [1]: https://github.com/chengjk/deployKit/blob/master/RELEASE.md
  [2]: https://eternallybored.org/misc/wget/
  [3]: https://github.com/chengjk/deployKit/releases