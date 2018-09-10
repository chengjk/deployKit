package common

import (
	"testing"
	"log"
	"path/filepath"
	"os"
)

//password
func TestConnectPassword(t *testing.T) {
	client, e := ConnPassword("vagrant", "vagrant", "192.168.65.90", 22)
	if e != nil {
		println(e)
	}
	if client != nil {
		defer client.Close();
	} else {
		log.Fatal("failed")
	}
}

//public key
func TestConnPublicKey(t *testing.T) {
	client, err := ConnPublicKey("ubuntu", "~/.ssh/dev.pem", "172.30.10.216", 22)
	if err != nil {
		log.Fatalf("unable to connect: %v", err)
	}
	if client != nil {
		defer client.Close();
	} else {
		log.Fatal("failed")
	}
}


func  TestPath(t *testing.T) {
	p := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	println(p)
}