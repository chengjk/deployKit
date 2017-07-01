package model

type ServerInfo struct {
	Host     string  `json:"Ip"`
	Port     int  `json:"Port"`
	Username string  `json:"Username"`
	Password string  `json:"Password"`
	WorkDir  string `json:"WorkDir"`
}
type Config struct {
	Name    string  `json:"Name"`
	Url     string `json:"Url"`
	Servers []ServerInfo `json:"Servers"`
}

type CmdParam struct {
	CfgFileName string
	Version     string
	/**
	资源类型:
	1:公网ftp,default. 直接在服务器上 wget
	2:局域网ftp; 先下载到本地，再上传的服务器
	3:本地zip文件;从本地上传服务器
	*/
	SourceType string
}