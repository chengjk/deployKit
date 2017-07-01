package main

import (
	"encoding/json"
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
	if len(os.Args)<3{
		fmt.Println("invalid param ")
		fmt.Println("cmd {project} {version}")
		return
	}
	var pName string=os.Args[1]
	var version string = os.Args[2]
	log.Println("Project:"+ pName+"; Version:"+ version)

	var workDir,_ =os.Getwd()
	var configPath= workDir +"\\"+pName+".json"
	log.Println("using config file :"+configPath)

	config := parseConfig(configPath)
	for _, server := range config.Servers {
		log.Println(server.Host)
		log.Println(version)
		deploy(version, server)
	}
}

func parseConfig(path string) (*model.Config) {
	fd, error := ioutil.ReadFile(path)
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
	localFilePath := "e:/"+version+".zip"
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
