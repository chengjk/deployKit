package main

import (
	"encoding/json"
	"errors"
	"fmt"
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
		deploy(cmdParam.Version, server)
		log.Println("deploy " + config.Name + " to server " + server.Host + " end.")
	}
}
func parseCmdParam() (model.CmdParam, error) {
	fmt.Println(len(os.Args))
	if len(os.Args) < 3 {

		fmt.Println("invalid param!")
		fmt.Println("cmd {project} {version} [sourceType]")
		fmt.Println("sourceType :")
		fmt.Println("	1 :公网ftp，默认")
		fmt.Println("	2 :局域网ftp")
		fmt.Println("	3 :本地zip文件")
		return model.CmdParam{}, errors.New("invalid cmd  param!")
	}
	var cmdParam model.CmdParam
	cmdParam.CfgFileName = os.Args[1]
	cmdParam.Version = os.Args[2]
	if len(os.Args) < 4 {
		cmdParam.SourceType = "1"
	} else {
		cmdParam.SourceType = os.Args[3]
	}
	log.Println("ConfigFile:" + cmdParam.CfgFileName + ".json; Version:" + cmdParam.Version + ";SourceType:" + string(cmdParam.SourceType))
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

func deploy(version string, server model.ServerInfo) {
	sshClient, err := common.Connect(server.Username, server.Password, server.Host, server.Port)
	if err != nil {
		log.Fatal(err)
	}
	localFilePath := "e:/" + version + ".zip"
	upload(sshClient, localFilePath, server.WorkDir)
	executeCmd(sshClient, server.WorkDir, localFilePath)
}
func executeCmd(sshClient *ssh.Client, basePath string, localFilePath string) {
	// create session
	var session *ssh.Session
	var err error
	if session, err = sshClient.NewSession(); err != nil {
		log.Fatal(err)
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Run("cd " + basePath + ";unzip -o " + path.Base(localFilePath))
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
