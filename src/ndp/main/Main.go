package main

import (
	"encoding/json"
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
	parseConfig()
	//connectSever()
}

func parseConfig() {
	fd, error := ioutil.ReadFile("E:/github/dk/src/config.json")
	if error != nil {
		panic(error)
	}
	log.Print(string(fd))
	var connInfo = &model.ConnectInfo{}
	json.Unmarshal(fd, connInfo)
	println(connInfo.Host)
}

func connectSever() {
	var (
		host          = "172.30.10.83"
		port          = 22
		username      = "root"
		pwd           = "12354"
		localFilePath = "e:/1.zip"
		remoteDir     = "/tmp/jackytest/"
	)
	sshClient, err := common.Connect(username, pwd, host, port)
	if err != nil {
		log.Fatal(err)
	}
	upload(sshClient, localFilePath, remoteDir)
	executeCmd(sshClient, remoteDir, localFilePath)
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
	session.Run("cd " + basePath + ";unzip " + path.Base(localFilePath))
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
