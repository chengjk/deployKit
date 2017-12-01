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
	if cmdParam.ShowVersion{
		print("develop kit version: "+model.Version)
		return
	}
	log.Println("ConfigFile:" + cmdParam.CfgFileName + ".json")
	cmdParam, config := mergeWithCfgFile(cmdParam)
	//deploy to servers
	for _, server := range config.Servers {
		log.Println("deploy " + config.Name + " to server " + server.Host + " start.")
		deploy(cmdParam, server)
		log.Println("deploy " + config.Name + " to server " + server.Host + " end.")
	}
}
func mergeWithCfgFile(cmdParam model.CmdParam) (model.CmdParam, *model.Config) {
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
	cmdparam.Verify(cmdParam)
	return cmdParam, config
}

func parseConfig(cfgFileName string) *model.Config {
	workDir, _ := os.Getwd()
	cfgFilePath := workDir + "/" + cfgFileName + ".json"
	log.Println("using config file " + cfgFilePath)
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
		replace := strings.Replace(cmdParam.Url, "{tag}", cmdParam.Tag, -1)
		log.Println("from internet repository :" + replace)
		cmds = []string{
			"wget " + cmdParam.Url,
			strings.Replace(cmdParam.SuffixCmd, "{tag}", cmdParam.Tag, -1)}
		cmdhelper.ExecRemote(sshClient, server.WorkDir, cmds)
		return
	}
	if cmdParam.Path != "" {
		replace := strings.Replace(cmdParam.Path, "{tag}", cmdParam.Tag, -1)
		log.Println("from path :" + replace)
		filehelper.Upload(sshClient, replace, server.WorkDir)
		cmds = []string{strings.Replace(cmdParam.SuffixCmd, "{tag}", cmdParam.Tag, -1)}
		cmdhelper.ExecRemote(sshClient, server.WorkDir, cmds)
		return
	}
	if cmdParam.LocalUrl != "" {
		replace := strings.Replace(cmdParam.LocalUrl, "{tag}", cmdParam.Tag, -1)
		log.Println("from local repository :" + replace)
		localPath := filehelper.Download(replace)
		filehelper.Upload(sshClient, localPath, server.WorkDir)
		cmds = []string{strings.Replace(cmdParam.SuffixCmd, "{tag}", cmdParam.Tag, -1)}
		cmdhelper.ExecRemote(sshClient, server.WorkDir, cmds)
	}
}
