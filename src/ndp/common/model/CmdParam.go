package model

type CmdParam struct {
	CfgFileName string
	Tag         string
	Url         string //公网ftp,default. 直接在服务器上 wget
	LocalUrl    string //局域网ftp; 先下载到本地，再上传的服务器
	Path        string //本地zip文件;从本地上传服务器
	PrefixCmd        string //前置命令，文件上传前
	SuffixCmd        string //后置命令，文件上传后
	ShowVersion bool   //show current version
}
