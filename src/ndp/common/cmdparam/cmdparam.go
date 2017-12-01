package cmdparam

import (
	"flag"
	"fmt"
	"log"
	"ndp/common/model"
	"os"
)

func Parse() (model.CmdParam, error) {
	cfgFileName := flag.String("name", "config", "项目名称，对应配置文件名,默认是config. e.g. ec 代表使用配置文件ec.json.无效时报错。")
	tag := flag.String("tag", "", "标签.可作为变量{tag}使用。")
	url := flag.String("url", "", "外网仓库地址.可以直接在服务器上 wget. e.g. http://test.com/{tag}/a.zip.")
	localUrl := flag.String("lurl", "", "内网仓库地址.服务器不能直接访问,需要先下载到本地磁盘再上传服务器. e.g. http://127.0.0.1/{tag}/a.zip.")
	zipPath := flag.String("path", "", "本地磁盘路径，直接上传服务器. e.g. /tmp/{tag}/a.zip.")
	cmd := flag.String("cmd", "echo hello;", "后置命令,文件上传成功后在server的workDir中执行的命令，以分号隔开。可以使用变量{tag}")
	v:=flag.Bool("v", false,"show current version.")

	flag.Usage = ShowUsage
	flag.Parse()
	var param model.CmdParam
	param.CfgFileName = *cfgFileName
	param.Tag = *tag
	param.Url = *url
	param.LocalUrl = *localUrl
	param.Path = *zipPath
	param.SuffixCmd = *cmd
	param.ShowVersion=*v
	return param, nil
}

func ShowUsage(){
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println("")
	fmt.Println("Tips: url,path,和lurl三个参数互斥,按照上述顺序检查到一个有效值时停止,否则报错.")
}

func Verify(cmdParam model.CmdParam)(bool) {
	if cmdParam.Url == "" && cmdParam.LocalUrl == "" && cmdParam.Path == "" {
		log.Fatal("one of url,lurl or path is required!")
		return false
	}
	if cmdParam.Tag == "" {
		log.Fatal("project tag is required!")
		return false
	}
	if cmdParam.SuffixCmd == "" {
		log.Fatal("suffix cmd is required!")
		return false
	}
	return true
}
