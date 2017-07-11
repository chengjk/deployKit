package cmdparam

import (
	"flag"
	"fmt"
	"log"
	"ndp/common/model"
	"os"
)

func Parse() (model.CmdParam, error) {
	cfgFileName := flag.String("name", "", "项目名称，对应配置文件名,默认是config. e.g. ec 代表使用配置文件ec.json.无效时报错。")
	version := flag.String("v", "", "版本，是zip文件名的一部分. e.g. v1.0.0.")
	url := flag.String("url", "", "外网仓库地址.可以直接在服务器上 wget. e.g. http://test.com/{version}/a.zip.")
	localUrl := flag.String("lurl", "", "内网仓库地址.服务器不能直接访问,需要先下载到本地磁盘再上传服务器. e.g. http://127.0.0.1/{version}/a.zip.")
	zipPath := flag.String("path", "", "本地磁盘路径，直接上传服务器. e.g. /tmp/{version}/a.zip.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("注意: url,path,和lurl三个参数互斥,按照上述顺序检查到一个有效值时停止,否则报错.")
	}
	flag.Parse()

	var param model.CmdParam
	param.CfgFileName = *cfgFileName
	param.Version = *version
	param.Url = *url
	param.LocalUrl = *localUrl
	param.Path = *zipPath
	log.Println("ConfigFile:" + param.CfgFileName + ".json")
	return param, nil

}
func Verify(cmdParam model.CmdParam) {
	if cmdParam.Url == "" && cmdParam.LocalUrl == "" && cmdParam.Path == "" {
		log.Fatal("one of url,lurl or path is required!")
	}
	if cmdParam.Version == "" {
		log.Fatal("version is required!")
	}
}
