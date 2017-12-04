package cmdparam

import (
	"flag"
	"fmt"
	"log"
	"ndp/common/model"
	"os"
)

func Parse() (model.CmdParam, error) {
	name := flag.String("name", "config", "名称,可作为变量{name}使用,对应配置文件名.不能使用变量. e.g. ec表示使用ec.json.")
	tag := flag.String("tag", "", "标签.可作为变量{tag}使用.不能使用变量. e.g. v1.0")
	url := flag.String("url", "", "外网仓库地址.直接在服务器上 wget. e.g. http://test.com/{tag}/{name}.zip.")
	localUrl := flag.String("lurl", "", "内网仓库地址,需要先下载到本地磁盘再上传服务器. e.g. http://127.0.0.1/{tag}/{name}.zip.")
	zipPath := flag.String("path", "", "目标文件在本地磁盘路径. e.g. /tmp/{tag}/{name}.zip.")
	pcmd := flag.String("pcmd", "", "prefix cmd,文件上传前在server的workDir中执行，分号隔开.e.g. mkdir p")
	scmd := flag.String("scmd", "", "suffix cmd,文件上传后在server的workDir中执行，分号隔开.e.g. rm -f *.zip")
	v:=flag.Bool("v", false,"show current version.")

	flag.Usage = ShowUsage
	flag.Parse()
	var param model.CmdParam
	param.Name = *name
	param.Tag = *tag
	param.Url = *url
	param.LocalUrl = *localUrl
	param.Path = *zipPath
	param.SuffixCmd = *scmd
	param.PrefixCmd = *pcmd
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
		log.Fatal("one and only one of url,lurl or path is required!")
		return false
	}
	if cmdParam.Tag == "" {
		log.Fatal("project tag is required!")
		return false
	}
	return true
}
