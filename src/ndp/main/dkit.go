package main

import (
	"encoding/json"
	"flag"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"ndp/common"
	"ndp/common/model"
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
	cfgFileName := flag.String("name", "", "project config file name. e.g. ec means ec.json.")
	version := flag.String("v", "", "target version. e.g. v1.0.0.")
	url := flag.String("url", "", "ftp url at internet. we will wget at server. e.g. http://test.com/a.zip.")
	localUrl := flag.String("lurl", "", "ftp url at local lan. we will download zip to disk at first,then upload to server . e.g. http://127.0.0.1/a.zip.")
	zipPath := flag.String("zf", "", "zip file path. we will upload zip file to server. e.g. /tmp/a.zip.")
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
	localFilePath :=  "e:/" + cmdParam.Version + ".zip"
	//workDir, _ := os.Getwd()
	//localFilePath :=  workDir+"/upload/" +cmdParam.CfgFileName+ cmdParam.Version + ".zip"
	var cmds []string
	if cmdParam.Url != "" {
		cmds = []string{
			"wget "+cmdParam.Url,
			"unzip -o " + path.Base(localFilePath)}
		executeCmd(sshClient, server.WorkDir, cmds)
	}
	if cmdParam.ZipPath != "" {
		upload(sshClient, localFilePath, server.WorkDir)
		cmds = []string{"unzip -o " + path.Base(localFilePath)}
		executeCmd(sshClient, server.WorkDir, cmds)

	}
	if cmdParam.LocalUrl != "" {
		//todo download to disk ,upload to server
		downloadFromLocalRepo(cmdParam.LocalUrl)
		upload(sshClient, localFilePath, server.WorkDir)
		cmds = []string{"unzip -o " + path.Base(localFilePath)}
		executeCmd(sshClient, server.WorkDir, cmds)
	}
}
func downloadFromLocalRepo(url string) {
	//todo impl

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
