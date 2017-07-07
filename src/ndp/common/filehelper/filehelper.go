package filehelper

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func Download(url string) string {
	log.Println("download from local Repo start.")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	filePath := path.Join("./upload/" + path.Base(url))
	//磁盘有的话就不再下载
	if _, err := os.Stat(filePath); os.IsExist(err) {
		return filePath
	}
	os.MkdirAll("./upload", 0777)
	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	io.Copy(f, resp.Body)
	log.Println("download from local Repo succeed!")
	return filePath
}
func Upload(sshClient *ssh.Client, localFilePath, remoteDir string) {
	log.Println("upload file start!")
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
	fmt.Print("uploading...")
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		fmt.Print(".")
		dstFile.Write(buf)
	}
	log.Println("copy file to remote server finished!")
}
