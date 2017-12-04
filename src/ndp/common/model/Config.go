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
	Tag       string       `json:"tag"`
	Url       string       `json:"url"`
	LUrl      string       `json:"lurl"`
	Path      string       `json:"path"`
	SuffixCmd string       `json:"suffixCmd"`
	PrefixCmd string       `json:"prefixCmd"`
	Servers   []ServerInfo `json:"servers"`
}
