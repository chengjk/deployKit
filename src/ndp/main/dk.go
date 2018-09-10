package main

import (
	"log"
	"ndp/common"
	"ndp/common/cmder"
	"ndp/common/filehelper"
	"ndp/common/model"
)

func main() {
	//解析命令行参数
	cmdParam, err := cmder.Parse()
	if err != nil {
		log.Fatal(err)
		return
	}
	if cmdParam.ShowVersion {
		println("develop kit version: " + model.Version)
		return
	}
	if cmdParam.Name == "" {
		cmdParam.Name = "config"
	}
	//解析配置文件
	config := model.ParseConfig(cmdParam.Name)
	cmdParam = model.MergeCfgFile(cmdParam, config)
	if cmder.Verify(cmdParam) {
		//deploy to servers
		for _, server := range config.Servers {
			log.Println("deploy " + config.Name + " to server " + server.Host + " start.")
			deploy(cmdParam, server)
			log.Println("deploy " + config.Name + " to server " + server.Host + " end.")
		}
	} else {
		cmder.ShowUsage()
	}
}

func deploy(cmdParam model.CmdParam, server model.ServerInfo) {
	sshClient, err := common.Connect(server)
	if err != nil {
		log.Fatal(err)
	}
	model.ReplaceVar(&cmdParam)
	//exe prefix cmd
	cmder.ExecRemote(sshClient, server.WorkDir, []string{cmdParam.PrefixCmd})
	//upload and exe suffix cmd
	if cmdParam.Url != "" {
		log.Println("from internet repository :" + cmdParam.Url)
		cmder.ExecRemote(sshClient, server.WorkDir, []string{
			"wget " + cmdParam.Url,
			cmdParam.SuffixCmd})
	} else if cmdParam.Path != "" {
		log.Println("from path :" + cmdParam.Path)
		filehelper.Upload(sshClient, cmdParam.Path, server.WorkDir)
		cmder.ExecRemote(sshClient, server.WorkDir, []string{cmdParam.SuffixCmd})
	} else if cmdParam.LocalUrl != "" {
		log.Println("from local repository :" + cmdParam.LocalUrl)
		localPath := filehelper.Download(cmdParam.LocalUrl)
		log.Print("upload from " + localPath)
		filehelper.Upload(sshClient, localPath, server.WorkDir)
		cmder.ExecRemote(sshClient, server.WorkDir, []string{cmdParam.SuffixCmd})
	} else {
		log.Println("nothing to upload,skip.")
		cmder.ExecRemote(sshClient, server.WorkDir, []string{cmdParam.SuffixCmd})
	}
}
