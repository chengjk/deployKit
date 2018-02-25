package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"ndp/common"
	"ndp/common/cmdhelper"
	"ndp/common/cmdparam"
	"ndp/common/filehelper"
	"ndp/common/model"
	"os"
	"strings"
)

func main() {
	cmdParam, err := cmdparam.Parse()
	if err != nil {
		log.Fatal(err)
		return
	}
	if cmdParam.ShowVersion {
		print("develop kit version: " + model.Version)
		return
	}
	cmdParam, config := mergeCfgFile(cmdParam)
	if cmdparam.Verify(cmdParam) {
		//deploy to servers
		for _, server := range config.Servers {
			log.Println("deploy " + config.Name + " to server " + server.Host + " start.")
			deploy(cmdParam, server)
			log.Println("deploy " + config.Name + " to server " + server.Host + " end.")
		}
	} else {
		cmdparam.ShowUsage()
	}
}
func mergeCfgFile(cmdParam model.CmdParam) (model.CmdParam, *model.Config) {
	if cmdParam.Name == "" {
		cmdParam.Name = "config"
	}
	config := parseConfig(cmdParam.Name)
	if cmdParam.Url == "" {
		cmdParam.Url = config.Url
	}
	if cmdParam.LocalUrl == "" {
		cmdParam.LocalUrl = config.LUrl
	}
	if cmdParam.Path == "" {
		cmdParam.Path = config.Path
	}
	if cmdParam.Tag == "" {
		cmdParam.Tag = config.Tag
	}
	if cmdParam.PrefixCmd == "" {
		cmdParam.PrefixCmd = config.PrefixCmd
	}
	if cmdParam.SuffixCmd == "" {
		cmdParam.SuffixCmd = config.SuffixCmd
	}
	return cmdParam, config
}

func parseConfig(cfgFileName string) *model.Config {
	workDir, _ := os.Getwd()
	cfgFilePath := workDir + "/" + cfgFileName + ".json"
	log.Println("using config file " + cfgFilePath)
	fd, error := ioutil.ReadFile(cfgFilePath)
	if error != nil {
		log.Fatal(error.Error())
	}
	var config = &model.Config{}
	json.Unmarshal(fd, config)
	return config
}

func deploy(cmdParam model.CmdParam, server model.ServerInfo) {
	sshClient, err := common.Connect(server)
	if err != nil {
		log.Fatal(err)
	}
	replaceVar(&cmdParam)

	//exe prefix cmd
	cmdhelper.ExecRemote(sshClient, server.WorkDir, []string{cmdParam.PrefixCmd})

	//upload and exe suffix cmd
	if cmdParam.Url != "" {
		log.Println("from internet repository :" + cmdParam.Url)
		cmdhelper.ExecRemote(sshClient, server.WorkDir, []string{
			"wget " + cmdParam.Url,
			cmdParam.SuffixCmd})
	} else if cmdParam.Path != "" {
		log.Println("from path :" + cmdParam.Path)
		filehelper.Upload(sshClient, cmdParam.Path, server.WorkDir)
		cmdhelper.ExecRemote(sshClient, server.WorkDir, []string{cmdParam.SuffixCmd})
	} else if cmdParam.LocalUrl != "" {
		log.Println("from local repository :" + cmdParam.LocalUrl)
		localPath := filehelper.Download(cmdParam.LocalUrl)
		log.Print("upload from " + localPath)
		filehelper.Upload(sshClient, localPath, server.WorkDir)
		cmdhelper.ExecRemote(sshClient, server.WorkDir, []string{cmdParam.SuffixCmd})
	}
}

func replaceVar(p *model.CmdParam) {
	p.Url = strings.Replace(p.Url, "{name}", p.Name, -1)
	p.Path = strings.Replace(p.Path, "{name}", p.Name, -1)
	p.LocalUrl = strings.Replace(p.LocalUrl, "{name}", p.Name, -1)
	p.PrefixCmd = strings.Replace(p.PrefixCmd, "{name}", p.Name, -1)
	p.SuffixCmd = strings.Replace(p.SuffixCmd, "{name}", p.Name, -1)

	p.Url = strings.Replace(p.Url, "{tag}", p.Tag, -1)
	p.Path = strings.Replace(p.Path, "{tag}", p.Tag, -1)
	p.LocalUrl = strings.Replace(p.LocalUrl, "{tag}", p.Tag, -1)
	p.PrefixCmd = strings.Replace(p.PrefixCmd, "{tag}", p.Tag, -1)
	p.SuffixCmd = strings.Replace(p.SuffixCmd, "{tag}", p.Tag, -1)
}
