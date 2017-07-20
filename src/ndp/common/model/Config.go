package model

type ServerInfo struct {
	Host     string `json:"ip"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	WorkDir  string `json:"workDir"`
}
type Config struct {
	Name      string       `json:"name"`
	Version   string       `json:"version"`
	Url       string       `json:"url"`
	LUrl      string       `json:"lurl"`
	Path      string       `json:"path"`
	SuffixCmd string       `json:"suffixCmd"`
	Servers   []ServerInfo `json:"servers"`
}
