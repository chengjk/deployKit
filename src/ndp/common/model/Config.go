package model

type ServerInfo struct {
	Host     string  `json:"Ip"`
	Port     int  `json:"Port"`
	Username string  `json:"Username"`
	Password string  `json:"Password"`
	WorkDir  string `json:"WorkDir"`
}
type Config struct {
	Name string  `json:"Name"`
	Servers []ServerInfo `json:"Servers"`
}
