package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"ndp/common"
	"ndp/common/model"
	"net/http"
	"os"
	"path"
)

func main() {
	cmdParam, err := parseCmdParam()
	if err != nil {
		log.Fatal(err)
		return
	}
	config := parseConfig(cmdParam.CfgFileName)
	//todo 多线程
	for _, server := range config.Servers {
		log.Println("deploy " + config.Name + " to server " + server.Host + " start.")
		deploy(cmdParam, server)
		log.Println("deploy " + config.Name + " to server " + server.Host + " end.")
	}
}
func parseCmdParam() (model.CmdParam, error) {
	cfgFileName := flag.String("name", "", "项目名称，对应配置文件名. e.g. ec 代表使用配置文件ec.json.无效时报错。")
	version := flag.String("v", "", "版本，是zip文件名的一部分. e.g. v1.0.0.")
	url := flag.String("url", "", "外网仓库地址.可以直接在服务器上 wget. e.g. http://test.com/a.zip.")
	localUrl := flag.String("lurl", "", "内网仓库地址.服务器不能直接访问,需要先下载到本地磁盘再上传服务器. e.g. http://127.0.0.1/a.zip.")
	zipPath := flag.String("path", "", "本地磁盘路径，直接上传服务器. e.g. /tmp/a.zip.")
	flag.Usage= func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Println("注意: url,path,和lurl三个参数互斥,按照上述顺序检查到一个有效值时停止,否则报错.")
	}
	flag.Parse()

	var cmdParam model.CmdParam
	cmdParam.CfgFileName = *cfgFileName
	cmdParam.Version = *version
	cmdParam.Url = *url
	cmdParam.LocalUrl = *localUrl
	cmdParam.ZipPath = *zipPath

	if cmdParam.CfgFileName == "" {
		log.Fatal("name is required!")
	}
	if cmdParam.Version == "" {
		log.Fatal("version is required!")
	}

	if cmdParam.Url == "" && cmdParam.LocalUrl == "" && cmdParam.ZipPath == "" {
		log.Fatal("one of url,lurl or zf is required!")
	}

	log.Println("ConfigFile:" + cmdParam.CfgFileName + ".json; Version:" + cmdParam.Version)
	return cmdParam, nil
}

func parseConfig(cfgFileName string) (*model.Config) {
	workDir, _ := os.Getwd()
	cfgFilePath := workDir + "\\" + cfgFileName + ".json"
	log.Println("using config file :" + cfgFilePath)
	fd, error := ioutil.ReadFile(cfgFilePath)
	if error != nil {
		panic(error)
	}
	var config = &model.Config{}
	json.Unmarshal(fd, config)
	return config
}

func deploy(cmdParam model.CmdParam, server model.ServerInfo) {
	sshClient, err := common.Connect(server.Username, server.Password, server.Host, server.Port)
	if err != nil {
		log.Fatal(err)
	}
	var cmds []string
	if cmdParam.Url != "" {
		cmds = []string{
			"wget "+cmdParam.Url,
			"unzip -o " + path.Base(cmdParam.Url)}
		executeCmd(sshClient, server.WorkDir, cmds)
	}
	if cmdParam.ZipPath != "" {
		localFilePath:=cmdParam.ZipPath+"/"+cmdParam.CfgFileName+"/"+cmdParam.Version+".zip"
		upload(sshClient, localFilePath, server.WorkDir)
		cmds = []string{"unzip -o " + path.Base(localFilePath)}
		executeCmd(sshClient, server.WorkDir, cmds)

	}
	if cmdParam.LocalUrl != "" {
		localPath := downloadFromLocalRepo(cmdParam.LocalUrl)
		upload(sshClient, localPath, server.WorkDir)
		cmds = []string{"unzip -o " + path.Base(localPath)}
		executeCmd(sshClient, server.WorkDir, cmds)
	}
}
func downloadFromLocalRepo(url string) string{
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	filePath := path.Join("./upload/" + path.Base(url))
	os.MkdirAll("./upload", 0777)
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	io.Copy(f, resp.Body)
	return filePath
}
func executeCmd(sshClient *ssh.Client, basePath string, cmds []string) {
	// create session
	var session *ssh.Session
	var err error
	if session, err = sshClient.NewSession(); err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Run("cd " + basePath)


	for _, cmd := range cmds {
		session.Run(cmd)
	}
	log.Println("unzip finished！")
}

func upload(sshClient *ssh.Client, localFilePath, remoteDir string) {
	var err error
	// create sftp client
	var sftpClient *sftp.Client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()
	//upload files
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()
	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remoteDir, remoteFileName))
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}
	log.Println("copy file to remote server finished!")
}
