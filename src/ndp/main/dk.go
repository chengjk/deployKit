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
		log.Println("ConfigFile:" + cmdParam.CfgFileName + ".json")
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
	if cmdParam.CfgFileName == "" {
		cmdParam.CfgFileName = "config"
	}
	config := parseConfig(cmdParam.CfgFileName)
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
	sshClient, err := common.Connect(server.Username, server.Password, server.Host, server.Port)
	if err != nil {
		log.Fatal(err)
	}
	replaceVar(&cmdParam)
	var cmdList []string
	if cmdParam.Url != "" {
		log.Println("from internet repository :" + cmdParam.Url)
		cmdList = []string{
			"wget " + cmdParam.Url,
			cmdParam.SuffixCmd}
		cmdhelper.ExecRemote(sshClient, server.WorkDir, cmdList)
		return
	}
	if cmdParam.Path != "" {
		log.Println("from path :" + cmdParam.Path)
		filehelper.Upload(sshClient, cmdParam.Path, server.WorkDir)
		cmdList = []string{cmdParam.SuffixCmd}
		cmdhelper.ExecRemote(sshClient, server.WorkDir, cmdList)
		return
	}
	if cmdParam.LocalUrl != "" {
		log.Println("from local repository :" + cmdParam.LocalUrl)
		localPath := filehelper.Download( cmdParam.LocalUrl)
		filehelper.Upload(sshClient, localPath, server.WorkDir)
		cmdList = []string{cmdParam.SuffixCmd}
		cmdhelper.ExecRemote(sshClient, server.WorkDir, cmdList)
	}
}

func replaceVar(p *model.CmdParam){
	p.Url= strings.Replace(p.Url, "{tag}", p.Tag, -1)
	p.Path= strings.Replace(p.Path, "{tag}", p.Tag, -1)
	p.LocalUrl= strings.Replace(p.LocalUrl, "{tag}", p.Tag, -1)
	p.SuffixCmd= strings.Replace(p.SuffixCmd, "{tag}", p.Tag, -1)
}