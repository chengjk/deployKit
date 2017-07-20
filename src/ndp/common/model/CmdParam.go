package model

type CmdParam struct {
	CfgFileName string
	Version     string
	Url         string //公网ftp,default. 直接在服务器上 wget
	LocalUrl    string //局域网ftp; 先下载到本地，再上传的服务器
	Path        string //本地zip文件;从本地上传服务器
}
