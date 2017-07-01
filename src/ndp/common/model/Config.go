package model



type ConnectInfo struct {
	Host     string  `json:"Host"`
	Port     string  `json:"Port"`
	Username string  `json:"Username"`
	Password string  `json:"Password"`
}
type PathInfo struct {
	Local     string  `json:"host"`
	Remote     string  `json:"port"`
}

