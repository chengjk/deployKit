package model

import "strings"

type CmdParam struct {
	Name string
	Tag         string
	Url         string //公网ftp,default. 直接在服务器上 wget
	LocalUrl    string //局域网ftp; 先下载到本地，再上传的服务器
	Path        string //本地zip文件;从本地上传服务器
	PrefixCmd        string //前置命令，文件上传前
	SuffixCmd        string //后置命令，文件上传后
	ShowVersion bool   //show current version
}

//合并命令行参数和配置文件
func MergeCfgFile(cmdParam CmdParam, config *Config) (CmdParam) {
	//配置文件名应该和 其中 name值一致。如果不一致，已配置文件里为准
	if config.Name != "" && cmdParam.Name != config.Name {
		cmdParam.Name = config.Name
	}
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
	return cmdParam
}

func ReplaceVar(p *CmdParam) {
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