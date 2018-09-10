package model

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type ServerInfo struct {
	Host      string `json:"ip"`
	Port      int    `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	PublicKey string `json:"publicKey"`
	WorkDir   string `json:"workDir"`
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



func ParseConfig(cfgFileName string) *Config {
	workDir, _ := os.Getwd()
	cfgFilePath := workDir + "/" + cfgFileName + ".json"
	log.Println("using config file " + cfgFilePath)
	fd, error := ioutil.ReadFile(cfgFilePath)
	if error != nil {
		log.Fatal(error.Error())
	}
	var config = &Config{}
	json.Unmarshal(fd, config)
	return config
}